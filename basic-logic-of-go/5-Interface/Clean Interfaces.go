package __Interface

import (
	"io"
	"os"
)

// 1. Keep Interfaces Small
type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

// 2. Interfaces Should Have No Knowledge of Satisfying Types
type car interface {
	Color() string
	Speed() int
	IsFiretruck() bool
}
type firetruck interface {
	car
	HoseLength() int
}
