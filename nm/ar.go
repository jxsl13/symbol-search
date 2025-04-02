package nm

import (
	"errors"
	"io"
	"io/fs"
	"time"

	"github.com/erikgeiser/ar"
	"github.com/jxsl13/archivewalker"
)

// https://www.abhirag.com/blog/ar/
func WalkAR(file archivewalker.File, walkFunc archivewalker.WalkFunc) (err error) {

	arr, err := ar.NewReader(file)
	if err != nil {
		return err
	}

	var header *ar.Header
	for {
		// defines a sub error in the loop scope
		header, err = arr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			err = walkFunc("", nil, nil, err)
			if err != nil {
				return err
			}
		}

		fi := &arInfo{header}

		err = walkFunc(header.Name, fi, arr, err)
		if err != nil {
			return err
		}

	}
}

var (
	_ fs.FileInfo = (*arInfo)(nil)
)

type arInfo struct {
	h *ar.Header
}

func (fi *arInfo) Name() string {
	return fi.h.Name
}

func (fi *arInfo) Size() int64 {
	return fi.h.Size
}

func (fi *arInfo) Mode() fs.FileMode {
	return fs.FileMode(fi.h.Mode)
}

func (fi *arInfo) ModTime() time.Time {
	return fi.h.ModTime
}

func (fi *arInfo) IsDir() bool {
	return fi.Mode().IsDir()
}

func (fi *arInfo) Sys() any {
	return fi.h
}
