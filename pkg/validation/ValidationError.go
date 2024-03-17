// Package validationutil provides helpers to validate
package validationutil

// ValidationError represents one validation error
type ValidationError[T any] struct {
	Problem string
	Source  T
}

// ValidationErrorList is a container to hold multiple ValidationErrors
type ValidationErrorList[T any] struct {
	errors []ValidationError[T]
}

func (l *ValidationErrorList[T]) Add(err ...ValidationError[T]) {
	l.errors = append(l.errors, err...)
}

func (l *ValidationErrorList[T]) Join(other ValidationErrorList[T]) {
	l.errors = append(l.errors, other.errors...)
}

func (l ValidationErrorList[T]) GetAll() []ValidationError[T] {
	return l.errors
}

func (l ValidationErrorList[T]) HasErrors() bool {
	return len(l.errors) > 0
}

func (l ValidationErrorList[T]) HasNoErrors() bool {
	return len(l.errors) == 0
}
