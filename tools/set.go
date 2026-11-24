package tools

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](items ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]struct{}),
	}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *Set[T]) Has(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[T]) Add(items ...T) *Set[T] {
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *Set[T]) Remove(item T) *Set[T] {
	delete(s.m, item)
	return s
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *Set[T]) Clear() *Set[T] {
	s.m = make(map[T]struct{})
	return s
}

func (s *Set[T]) Size() int {
	return len(s.m)
}

func (s *Set[T]) IsEmpty() bool {
	return len(s.m) == 0
}

func (s *Set[T]) Slice() []T {
	items := make([]T, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}
