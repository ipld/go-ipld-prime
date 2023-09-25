package storage

import (
	"context"
	"io"
)

type teeStorage struct {
	ReadableStorage
	out WritableStorage
}

type teeReadCloser struct {
	io.Reader
	readCloser io.Closer
	key        string
	writeClose func(string) error
}

func (trc teeReadCloser) Close() error {
	err := trc.readCloser.Close()
	if err != nil {
		return err
	}
	err = trc.writeClose(trc.key)
	if err != nil {
		return err
	}
	return nil
}

func (ts teeStorage) GetStream(ctx context.Context, key string) (io.ReadCloser, error) {
	rdr, err := GetStream(ctx, ts.ReadableStorage, key)
	if err != nil {
		return nil, err
	}
	writer, committer, err := PutStream(ctx, ts.out)
	if err != nil {
		return nil, err
	}
	return teeReadCloser{
		Reader:     io.TeeReader(rdr, writer),
		readCloser: rdr,
		writeClose: committer,
		key:        key,
	}, nil
}

func (ts teeStorage) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := ts.ReadableStorage.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	err = ts.out.Put(ctx, key, data)
	return data, err
}

func TeeStorage(in ReadableStorage, out WritableStorage) ReadableStorage {
	return teeStorage{
		ReadableStorage: in,
		out:             out,
	}
}
