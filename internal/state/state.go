package state

type State[T any] struct {
	value   T
	Changed bool
}

func UseState[T any](value T) *State[T] {
	return &State[T]{
		value: value,
	}
}

func (s *State[T]) SetState(fn func(prev T) T) {
	s.value = fn(s.value)
	s.Changed = true
}

func (s *State[T]) Value() T {
	return s.value
}

type StateAny interface {
	changed() bool
	toggleChanged()
}

func (s *State[T]) changed() bool {
	return s.Changed
}

func (s *State[T]) toggleChanged() {
	s.Changed = !s.Changed
}
func UseEffect(fn func(), dep []StateAny) {
	for {
		for _, s := range dep {
			if s.changed() {
				fn()
				s.toggleChanged()
			}
		}
	}
}
