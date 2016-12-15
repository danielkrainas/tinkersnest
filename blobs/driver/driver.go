package driver

import (
	"io"

	"github.com/danielkrainas/tinkersnest/blobs"
)

type Driver interface {
	Store(b *blobs.Blob, r io.Reader) error
	Exists(name string) (bool, error)
	Get(name string) (*blobs.Blob, io.Reader, error)
	Drop(name string) (bool, error)
}
