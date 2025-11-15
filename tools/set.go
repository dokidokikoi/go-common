package tools

type set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable](items ...T) *set[T] {
	s := &set[T]{}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

func (s *set[T]) Add(item T) {
	s.m[item] = struct{}{}
}

func (s *set[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}

func (s *set[T]) Slice() []T {
	items := make([]T, 0, len(s.m))
	for item := range s.m {
		items = append(items, item)
	}
	return items
}
