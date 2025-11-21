package tools

type set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](items ...T) *set[T] {
	s := &set[T]{
		m: make(map[T]struct{}),
	}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *set[T]) Add(items ...T) *set[T] {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *set[T]) Remove(item T) *set[T] {
	delete(s.m, item)
	return s
}

func (s *set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *set[T]) Clear() *set[T] {
	s.m = make(map[T]struct{})
	return s
}

func (s *set[T]) Size() int {
	return len(s.m)
}

func (s *set[T]) IsEmpty() bool {
	return len(s.m) == 0
}

func (s *set[T]) Slice() []T {
	items := make([]T, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}
