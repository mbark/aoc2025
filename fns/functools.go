package fns

func Every[T any](s []T, fn func(t T) bool) bool {
	for i := 0; i < len(s); i++ {
		if !fn(s[i]) {
			return false
		}
	}

	return true
}

func EveryMap[K comparable, V any](s map[K]V, fn func(key K, val V) bool) bool {
	for k, v := range s {
		if !fn(k, v) {
			return false
		}
	}

	return true
}

func EveryIndex[T any](s []T, fn func(t T, idx int) bool) bool {
	for i := 0; i < len(s); i++ {
		if !fn(s[i], i) {
			return false
		}
	}

	return true
}

func Some[T any](s []T, fn func(t T) bool) bool {
	for i := 0; i < len(s); i++ {
		if fn(s[i]) {
			return true
		}
	}

	return false
}

func FilterMap[T, R any](s []T, fn func(t T) (R, bool)) []R {
	var arr []R
	for _, t := range s {
		r, ok := fn(t)
		if ok {
			arr = append(arr, r)
		}
	}
	return arr
}

func Map[T, R any](s []T, fn func(t T) R) []R {
	arr := make([]R, len(s))
	for i, t := range s {
		arr[i] = fn(t)
	}
	return arr
}

func Associate[V any, K comparable](s []V, fn func(t V) K) map[K]V {
	m := make(map[K]V, len(s))
	for _, t := range s {
		m[fn(t)] = t
	}
	return m
}

func AsMap[T any, V any, K comparable](s []T, fn func(t T) (K, V)) map[K]V {
	m := make(map[K]V, len(s))
	for _, t := range s {
		k, v := fn(t)
		m[k] = v
	}
	return m
}

func Find[T any](s []T, fn func(t T) bool) (T, bool) {
	for _, t := range s {
		if fn(t) {
			return t, true
		}
	}

	var d T
	return d, false
}

func Filter[T any](s []T, fn func(t T) bool) []T {
	var arr []T
	for _, t := range s {
		if fn(t) {
			arr = append(arr, t)
		}
	}
	return arr
}

func Flatten[T any](ts [][]T) []T {
	length := 0
	for i := range ts {
		length += len(ts[i])
	}

	arr := make([]T, 0, length)
	for i := range ts {
		arr = append(arr, ts[i]...)
	}

	return arr
}

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i += 1
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, len(m))
	i := 0
	for _, v := range m {
		values[i] = v
		i += 1
	}
	return values
}

func Repeat[T any](t T, n int) []T {
	arr := make([]T, n)
	for i := 0; i < n; i++ {
		arr[i] = t
	}
	return arr
}

func Sum[T ~int | ~float64](s []T) T {
	var sum T
	for _, t := range s {
		sum += t
	}
	return sum
}

func Intersection[K comparable, V any](a, b map[K]V) map[K]V {
	out := map[K]V{}
	for k := range a {
		if _, ok := b[k]; ok {
			out[k] = b[k]
		}
	}
	return out
}

func Union[K comparable, V any](a, b map[K]V) map[K]V {
	out := map[K]V{}
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
