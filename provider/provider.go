package provider

import (
	"io"
)

type Provider interface {
	Copy(in io.Reader, name string) string
}
