package helper

func Apply[T1 any, T2 any](slice []T1, f func(T1) T2) []T2 {
	result := make([]T2, 0, len(slice))

	for _, item := range slice {
		result = append(result, f(item))
	}
	return result
}

// apply with indeces
func Apply2[T1 any, T2 any](slice []T1, f func(int, T1) T2) []T2 {
	result := make([]T2, 0, len(slice))

	for i, item := range slice {
		result = append(result, f(i, item))
	}
	return result
}

func Fold[T1 any, T2 any](slice []T1, init T2, f func(T2, T1) T2) T2 {
	acc := init
	for _, item := range slice {
		acc = f(acc, item)
	}
	return acc
}

func ApplyInPlace[T any](slice []T, f func(T) T) {
	for i, item := range slice {
		slice[i] = f(item)
	}
}

func Filter[T any](slice []T, f func(T) bool) []T {
	result := []T{}

	for _, item := range slice {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}
