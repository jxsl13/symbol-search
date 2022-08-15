package archive

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var supportedExtensions = map[string]bool{
	".gz":  true,
	".tgz": true,
	".xz":  true,
	".tar": true,
	".zip": true,
	".7z":  true,
}

// WalkFunc defines the function in order to efficiently walk over the archive
type WalkFunc func(path string, info fs.FileInfo, file io.ReaderAt, err error) error

func IsSupported(path string) bool {
	return supportedExtensions[filepath.Ext(path)]
}

func Walk(path string, walkcFunc WalkFunc) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".gz", ".tgz":
		return WalkTarGzip(f, walkcFunc)
	case ".xz":
		return WalkTarXz(f, walkcFunc)
	case ".tar":
		return WalkTar(f, walkcFunc)
	case ".zip":
		return WalkZip(f, stat.Size(), walkcFunc)
	case ".7z":
		return Walk7Zip(f, stat.Size(), walkcFunc)
	}
	return fmt.Errorf("unknown file extension: %s", ext)
}

// newReaderAt closes the passed file handle
func newReaderAt(fi io.Reader, size int64) (io.ReaderAt, error) {
	if c, ok := fi.(io.Closer); ok {
		defer c.Close()
	}

	buf := bytes.NewBuffer(make([]byte, 0, size))
	written, err := io.Copy(buf, fi)
	if err != nil {
		return nil, err
	}
	if written != size {
		return nil, fmt.Errorf("could buffer file in archive: size mismatch: expected %d, got %d", size, written)
	}

	return bytes.NewReader(buf.Bytes()), nil
}
