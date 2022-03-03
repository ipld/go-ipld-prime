package selector

import "io"

type readerat struct {
	io.ReadSeeker
}

func (r readerat) ReadAt(p []byte, off int64) (n int, err error) {
	_, err = r.Seek(off, 0)
	if err != nil {
		return 0, err
	}
	return r.Read(p)
}
