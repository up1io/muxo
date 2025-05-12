package local

import "github.com/leonelquinteros/gotext"

func Text(s string) string {
	return gotext.Get(s)
}
