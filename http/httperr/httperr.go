package httperr

type APIError interface {
	Status() int
	Error() string
}
