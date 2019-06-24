package kahla

import (
	"io"
)

type RequestFile interface {
	io.Reader
	Name() string
}
