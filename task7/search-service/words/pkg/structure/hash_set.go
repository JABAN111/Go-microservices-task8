package hashset

type HashSet[T comparable] map[T]struct{}

func New[T comparable]() HashSet[T] {
	return make(HashSet[T])
}

func (s HashSet[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s HashSet[T]) Remove(value T) {
	delete(s, value)
}

func (s HashSet[T]) Contains(value T) bool {
	_, exists := s[value]
	return exists
}

func (s HashSet[T]) Size() int {
	return len(s)
}

func (s HashSet[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for key := range s {
		result = append(result, key)
	}
	return result
}
