package locales

type Reader interface {
	Text(s string) string
}
