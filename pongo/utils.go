package pongo

func ListDiff[T comparable](a, b []T) []T {
	mb := make(map[T]any, len(b))
	for _, x := range b {
		mb[x] = nil
	}

	return ListMapDiff(a, mb)
}

func ListMapDiff[T comparable](a []T, mb map[T]interface{}) []T {
	var diff = []T{}
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func Values[K comparable, V any](m map[K]V) []V {
	var arr []V
	for _, v := range m {
		arr = append(arr, v)
	}

	return arr
}
