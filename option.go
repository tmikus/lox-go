package main

type Option[T interface{}] struct {
	value interface{}
}

func (o *Option[T]) IsEmpty() bool {
	return o.value == nil
}

func (o *Option[T]) IsSome() bool {
	return o.value != nil
}

func (o *Option[T]) Value() T {
	return o.value.(T)
}

func (o *Option[T]) ValueOr(defaultValue T) T {
	if o.IsEmpty() {
		return defaultValue
	}
	return o.Value()
}

func (o *Option[T]) ValueOrPanic() T {
	if o.IsEmpty() {
		panic("Value is empty")
	}
	return o.Value()
}

func NewOption[T interface{}](value interface{}) Option[T] {
	return Option[T]{value}
}

func NewEmptyOption[T interface{}]() Option[T] {
	return Option[T]{nil}
}

func NewSomeOption[T interface{}](value T) Option[T] {
	return Option[T]{value}
}

func NewNoneOption[T interface{}]() Option[T] {
	return Option[T]{nil}
}
