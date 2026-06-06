package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

type storageUnit struct {
	key  string
	data []byte
}

// onceError is an object that will only store an error once.
type onceError struct {
	sync.Mutex // guards following
	err        error
}

func (a *onceError) Store(err error) {
	a.Lock()
	defer a.Unlock()
	if a.err != nil {
		return
	}
	a.err = err
}
func (a *onceError) Load() error {
	a.Lock()
	defer a.Unlock()
	return a.err
}

// ErrMismatchedKey is the error used when read and write keys don't match
var ErrMismatchedKey = errors.New("put/get keys do not match")

// A pipeStorage is the shared pipeStorage structure underlying ReadablePipeStorage and WritablePipeStorage.
type pipeStorage struct {
	wrMu sync.Mutex // Serializes Write operations
	wrCh chan storageUnit
	rdCh chan struct{}

	once   sync.Once // Protects closing done
	done   chan struct{}
	rerr   onceError
	werr   onceError
	keysLk sync.RWMutex
	keys   map[string]struct{}
}

func (p *pipeStorage) has(ctx context.Context, key string) (bool, error) {
	select {
	case <-p.done:
		return false, p.readableFinalizeError()
	default:
	}
	p.keysLk.RLock()
	defer p.keysLk.RUnlock()
	_, ok := p.keys[key]
	return ok, nil
}

func (p *pipeStorage) get(ctx context.Context, key string) ([]byte, error) {
	select {
	case <-p.done:
		return nil, p.readableFinalizeError()
	default:
	}

	select {
	case su := <-p.wrCh:
		p.rdCh <- struct{}{}
		if su.key != key {
			return nil, fmt.Errorf("%w: put %s, got %s", ErrMismatchedKey, su.key, key)
		}
		return su.data, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-p.done:
		return nil, p.readableFinalizeError()
	}
}

func (p *pipeStorage) finalizeReadable(err error) error {
	if err == nil {
		err = io.ErrClosedPipe
	}
	p.rerr.Store(err)
	p.once.Do(func() { close(p.done) })
	return nil
}

func (p *pipeStorage) put(ctx context.Context, key string, b []byte) error {
	select {
	case <-p.done:
		return p.writableFinalizeError()
	default:
		p.wrMu.Lock()
		defer p.wrMu.Unlock()
	}

	select {
	case p.wrCh <- storageUnit{key, b}:
		p.keysLk.Lock()
		p.keys[key] = struct{}{}
		p.keysLk.Unlock()
		<-p.rdCh
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.done:
		return p.writableFinalizeError()
	}
}

func (p *pipeStorage) finalizeWritable(err error) error {
	if err == nil {
		err = io.EOF
	}
	p.werr.Store(err)
	p.once.Do(func() { close(p.done) })
	return nil
}

// readableFinalizeError is considered internal to the pipe type.
func (p *pipeStorage) readableFinalizeError() error {
	rerr := p.rerr.Load()
	if werr := p.werr.Load(); rerr == nil && werr != nil {
		return werr
	}
	return io.ErrClosedPipe
}

// writableFinalizeError is considered internal to the pipe type.
func (p *pipeStorage) writableFinalizeError() error {
	werr := p.werr.Load()
	if rerr := p.rerr.Load(); werr == nil && rerr != nil {
		return rerr
	}
	return io.ErrClosedPipe
}

// A ReadablePipeStorage is the read half of a pipe.
type ReadablePipeStorage struct {
	p *pipeStorage
}

// Has implements the Storage interface
func (r *ReadablePipeStorage) Has(ctx context.Context, key string) (bool, error) {
	return r.p.has(ctx, key)
}

// Get implements the ReadableStorage interface:
// it reads data from the pipe, blocking until a writer
// arrives or the write end is closed.
// If the write end is closed with an error, that error is
// returned as err; otherwise err is EOF.
func (r *ReadablePipeStorage) Get(ctx context.Context, key string) ([]byte, error) {
	return r.p.get(ctx, key)
}

// Finalize closes the reader; subsequent writes to the
// write half of the pipe will return the error ErrClosedPipe.
func (r *ReadablePipeStorage) Finalize() error {
	return r.FinalizeWithError(nil)
}

// FinalizeWithError closes the reader; subsequent writes
// to the write half of the pipe will return the error err.
//
// FinalizeWithError never overwrites the previous error if it exists
// and always returns nil.
func (r *ReadablePipeStorage) FinalizeWithError(err error) error {
	return r.p.finalizeReadable(err)
}

// A WritablePipeStorage is the write half of a pipe.
type WritablePipeStorage struct {
	p *pipeStorage
}

// Has implements the Storage interface
func (w *WritablePipeStorage) Has(ctx context.Context, key string) (bool, error) {
	return w.p.has(ctx, key)
}

// Put implements the standard Write interface:
// it writes data to the pipe, blocking until one or more readers
// have consumed all the data or the read end is closed.
// If the read end is closed with an error, that err is
// returned as err; otherwise err is ErrClosedPipe.
func (w *WritablePipeStorage) Put(ctx context.Context, key string, data []byte) error {
	return w.p.put(ctx, key, data)
}

// Finalize closes the writer; subsequent reads from the
// read half of the pipe will return no bytes and EOF.
func (w *WritablePipeStorage) Finalize() error {
	return w.FinalizeWithError(nil)
}

// FinalizeWithError closes the writer; subsequent reads from the
// read half of the pipe will return no bytes and the error err,
// or EOF if err is nil.
//
// FinalizeWithError never overwrites the previous error if it exists
// and always returns nil.
func (w *WritablePipeStorage) FinalizeWithError(err error) error {
	return w.p.finalizeWritable(err)
}

// PipeStorage creates a synchronous in-memory pipe.
// It can be used to connect code expecting an io.Reader
// with code expecting an io.Writer.
//
// Reads and Writes on the pipe are matched one to one
// except when multiple Reads are needed to consume a single Write.
// That is, each Write to the PipeWriter blocks until it has satisfied
// one or more Reads from the PipeReader that fully consume
// the written data.
// The data is copied directly from the Write to the corresponding
// Read (or Reads); there is no internal buffering.
//
// It is safe to call Read and Write in parallel with each other or with Close.
// Parallel calls to Read and parallel calls to Write are also safe:
// the individual calls will be gated sequentially.
func PipeStorage() (*ReadablePipeStorage, *WritablePipeStorage) {
	p := &pipeStorage{
		wrCh: make(chan storageUnit),
		rdCh: make(chan struct{}),
		done: make(chan struct{}),
	}
	return &ReadablePipeStorage{p}, &WritablePipeStorage{p}
}
