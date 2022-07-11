package preload

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"sync"
	"sync/atomic"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/linking"
)

type bytesLoader interface {
	Bytes() []byte
}

type preloadingLink struct {
	refCnt     uint64
	loadSyncer sync.Once
	loaded     chan struct{}
	data       []byte
	err        error
}

type request struct {
	lnkCtx linking.LinkContext
	lnk    datamodel.Link
}

type BufferedLoader struct {
	allocationLimit uint64
	concurrency     uint64
	avgBlockSize    uint64
	allocated       uint64
	dealloc         chan struct{}
	originalLoader  linking.BlockReadOpener
	preloadsLk      sync.RWMutex
	preloads        map[ipld.Link]*preloadingLink
	requests        chan request
}

func NewBufferedLoader(loader linking.BlockReadOpener, allocationLimit uint64, avgBlockSize uint64, concurrency uint64) *BufferedLoader {
	return &BufferedLoader{
		allocationLimit: allocationLimit,
		avgBlockSize: avgBlockSize,
		concurrency: concurrency,
		originalLoader: loader,
		dealloc: make(chan struct{}, 1),
		preloads: make(map[datamodel.Link]*preloadingLink),
		requests: make(chan request),
	}
}

func (bl *BufferedLoader) Preloader(ctx Context, links []Link) {
	bl.preloadsLk.Lock()
	defer bl.preloadsLk.Unlock()
	for _, l := range links {
		if pl, existing := bl.preloads[l.Link]; existing {
			pl.refCnt++
			continue
		}
		bl.preloads[l.Link] = &preloadingLink{
			loaded: make(chan struct{}),
			refCnt: 1,
		}
		select {
		case <-ctx.Ctx.Done():
		case bl.requests <- request{
			lnkCtx: linking.LinkContext{
				Ctx:        ctx.Ctx,
				LinkPath:   ctx.BasePath.AppendSegment(l.Segment),
				LinkNode:   l.LinkNode,
				ParentNode: ctx.ParentNode,
			},
			lnk: l.Link,
		}:
		}

	}
}

type byteReader interface {
	Bytes() []byte
}

func (bl *BufferedLoader) Start(ctx context.Context) {
	go bl.run(ctx)
}

func (bl *BufferedLoader) run(ctx context.Context) {
	feed := make(chan request)
	loadComplete := make(chan struct{})

	var wg sync.WaitGroup

	loaderCtx, cancel := context.WithCancel(ctx)
	defer wg.Wait()
	defer cancel()
	for i := uint64(0); i < bl.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for request := range feed {
				bl.preloadsLk.RLock()
				pl, ok := bl.preloads[request.lnk]
				bl.preloadsLk.RUnlock()
				if !ok {
					continue
				}
				bl.preloadLink(pl, request.lnkCtx, request.lnk)
				select {
				case <-loaderCtx.Done():
				case loadComplete <- struct{}{}:
				}
			}
		}()
	}
	defer close(feed)

	var send chan<- request
	var requestBuffer []request
	var next *request
	inProgress := uint64(0)
	for {
		select {
		case request := <-bl.requests:
			if next == nil {
				next = &request
				if bl.roomToAllocate(inProgress) {
					send = feed
				}
			} else {
				requestBuffer = append(requestBuffer, request)
			}
		case <-bl.dealloc:
			if next != nil && bl.roomToAllocate(inProgress) {
				send = feed
			}
		case send <- *next:
			inProgress++
			if len(requestBuffer) > 0 {
				next = &requestBuffer[0]
				requestBuffer = requestBuffer[1:]
			} else {
				next = nil
			}
			if next == nil || bl.roomToAllocate(inProgress) {
				send = nil
			}
		case <-loadComplete:
			inProgress--
			if next != nil && bl.roomToAllocate(inProgress) {
				send = feed
			}
		case <-ctx.Done():
			return
		}
	}
}

func (bl *BufferedLoader) roomToAllocate(inProgress uint64) bool {
	return atomic.LoadUint64(&bl.allocated)+(inProgress*bl.avgBlockSize) < bl.allocationLimit
}

func (bl *BufferedLoader) preloadLink(pl *preloadingLink, lnkCtx linking.LinkContext, lnk datamodel.Link) {
	pl.loadSyncer.Do(func() {
		defer close(pl.loaded)
		reader, err := bl.originalLoader(lnkCtx, lnk)
		if err != nil {
			pl.err = err
		} else {
			if br, ok := reader.(byteReader); ok {
				pl.data = br.Bytes()
			} else {
				pl.data, pl.err = ioutil.ReadAll(reader)
			}
		}
		atomic.AddUint64(&bl.allocated, uint64(len(pl.data)))
	})
}

func (bl *BufferedLoader) Load(lnkCtx linking.LinkContext, lnk datamodel.Link) (io.Reader, error) {
	bl.preloadsLk.Lock()
	pl, ok := bl.preloads[lnk]
	if ok {
		pl.refCnt--
		if pl.refCnt <= 0 {
			delete(bl.preloads, lnk)
		}
	}
	bl.preloadsLk.Unlock()
	if !ok {
		return bl.originalLoader(lnkCtx, lnk)
	}
	bl.preloadLink(pl, lnkCtx, lnk)
	select {
	case <-lnkCtx.Ctx.Done():
		return nil, lnkCtx.Ctx.Err()
	case <-pl.loaded:
		if pl.err != nil {
			return nil, pl.err
		}
		if pl.refCnt <= 0 {
			atomic.AddUint64(&bl.allocated, ^uint64(len(pl.data)))
			select {
			case bl.dealloc <- struct{}{}:
			default:
			}
		}
		return bytes.NewBuffer(pl.data), nil
	}
}
