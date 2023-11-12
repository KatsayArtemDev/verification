package result

import "fmt"

type Result[T any] struct {
	Error string `json:"error"`
	Value T      `json:"value"`
}

func Fail(err error) Result[any] {
	return Result[any]{
		Error: err.Error(),
		Value: nil,
	}
}

func HttpFail(path string, code int, err error) Result[any] {
	return Fail(fmt.Errorf("error along the way: %s; with code: %d; message: %w", path, code, err))
}

func Success[T any](value T) Result[T] {
	return Result[T]{
		Value: value,
	}
}
