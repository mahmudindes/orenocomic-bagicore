package model

var (
	ErrGeneric  GenericError
	ErrNotFound notFoundError
	ErrDatabase DatabaseError
	ErrCache    cacheError
)

type GenericError string

func (e GenericError) Error() string { return string(e) }

type wrappedError struct {
	msg string
	err error
}

func (e wrappedError) Error() string   { return e.msg }
func (e wrappedError) Unwrap() []error { return []error{e.err, GenericError(e.msg)} }

func WrappedError(err error, msg string) error {
	return wrappedError{msg, err}
}

type notFoundError struct {
	err error
}

func (e notFoundError) Error() string { return "not found: " + e.err.Error() }
func (e notFoundError) Unwrap() error { return e.err }

func NotFoundError(err error) error {
	return notFoundError{err}
}

type DatabaseError struct {
	Name string
	Code string
	Err  error
}

func (e DatabaseError) Error() string { return "database error: " + e.Err.Error() }
func (e DatabaseError) Unwrap() error { return e.Err }

type cacheError struct {
	err error
}

func (e cacheError) Error() string { return "cache failed: " + e.err.Error() }
func (e cacheError) Unwrap() error { return e.err }

func CacheError(err error) error {
	return cacheError{err}
}
