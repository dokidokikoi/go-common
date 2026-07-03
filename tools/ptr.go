package tools

func Ptr[T any](t T) *T {
	return &t
}

func PeelPtr[T any](t *T) T {
	return *t
}
