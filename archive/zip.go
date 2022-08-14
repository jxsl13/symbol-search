package archive

import (
	"archive/zip"
	"io/fs"
	"os"
)

func WalkZip(file *os.File, fileSize int64, walkFunc WalkFunc) error {
	f, err := zip.NewReader(file, fileSize)
	if err != nil {
		return err
	}
	return fs.WalkDir(f, "/", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return walkFunc(path, nil, nil, err)
		}

		currentFile, err := f.Open(path)
		if err != nil {
			return walkFunc(path, info, nil, err)
		}

		ra, err := newReaderAt(currentFile, info.Size())
		return walkFunc(path, info, ra, err)
	})
}
