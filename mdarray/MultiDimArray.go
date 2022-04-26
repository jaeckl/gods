package mdarray

type MultiDimArray[T any] struct {
	shape []int
	axes  map[int] map[int]T
}

type axe struct {
	size int
	data map[int]T
	axis *axe
}

func (md *MultiDimArray[T]) Shape() []int {
	return md.shape
}

func (md *MultiDimArray[T]) Slice(indices []int) view {
	if len(indices) > len(md.shape)   {
		panic( "dimensions out of bounds")
	}
	if len(indices) == len(md.shape)   {
		// validate indices in bound

	}
}

func New[T any](shape []int) *MultiDimArray[T] {
	return &MultiDimArray[T]{
		shape: shape,
		axes:  make(map[int][]T),
	}
}

func (md *MultiDimArray[T]) SetField(indices []int, data T) {

}
func (md *MultiDimArray[T]) GetField(indices []int, ) T {

}

struct view interface {
}

