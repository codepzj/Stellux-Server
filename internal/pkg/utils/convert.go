package utils

import "github.com/samber/lo"

func PtrToValList[T any](ptrSlice []*T) []T {
	return lo.Map(ptrSlice, func(item *T, _ int) T {
		if item == nil {
			return *new(T)
		}
		return *item
	})
}

func ValToPtrList[T any](valueSlice []T) []*T {
	return lo.Map(valueSlice, func(item T, _ int) *T {
		return &item
	})
}
