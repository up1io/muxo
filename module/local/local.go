package local

import "github.com/leonelquinteros/gotext"

type Getter interface {
	Text(s string) string
}

func Text(s string) string {
	return gotext.Get(s)
}
