package driver

import (
	"io"

	"github.com/danielkrainas/tinkersnest/blobs"
)

type Driver interface {
	Inspect(name string) (*blobs.Blob, error)
	Writer(name string) (io.WriteCloser, error)
	Reader(name string) (io.ReadCloser, error)
	WriteMeta(name string, b *blobs.Blob) error
	Drop(name string) (bool, error)
}
