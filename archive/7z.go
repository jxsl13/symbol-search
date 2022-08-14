package archive

import (
	"io/fs"
	"os"

	"github.com/bodgit/sevenzip"
)

func Walk7Zip(file *os.File, fileSize int64, walkFunc WalkFunc) error {
	f, err := sevenzip.NewReader(file, fileSize)
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
		if err != nil {
			return walkFunc(path, info, nil, err)
		}

		return walkFunc(path, info, ra, nil)
	})
}
