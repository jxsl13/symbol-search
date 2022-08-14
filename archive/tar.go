package archive

import (
	"archive/tar"
	"errors"
	"io"
)

// WalkTar may be passed a compressed reader instead of an explicit file
func WalkTar(file io.Reader, walkFunc WalkFunc) error {

	tr := tar.NewReader(file)

	for {
		// defines a sub error in the loop scope
		header, err := tr.Next()

		switch {
		// if no more files are found return
		case errors.Is(err, io.EOF):
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		fi := header.FileInfo()
		ra, err := newReaderAt(tr, fi.Size())

		// the target location where the dir/file should be created
		err = walkFunc(header.Name, fi, ra, err)
		if err != nil {
			return err
		}
	}
}
