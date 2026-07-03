package tools

import "slices"

func MergeSlice[T comparable](s1, s2 []T) []T {
	s1 = append(s1, s2...)
	return slices.Compact(s1)
}
