package handler

// contains checks if the given element `e` is presented in the slice `s`.
// it makes use of the GO generics supported from 1.18
func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
