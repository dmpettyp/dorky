package mapper

import "fmt"

type Mapper[From, To comparable] struct {
	to   map[From]To
	from map[To]From
}

func New[From, To comparable](
	values ...any,
) (
	*Mapper[From, To],
	error,
) {
	m := &Mapper[From, To]{
		to:   make(map[From]To),
		from: make(map[To]From),
	}

	if len(values)%2 == 1 {
		return nil, fmt.Errorf("odd number of key/values")
	}

	for {
		if len(values) == 0 {
			break
		}

		k, ok := values[0].(From)

		if !ok {
			return nil, fmt.Errorf(
				"expected key of type %T, got %T", *new(From), values[0],
			)
		}

		if _, ok := m.to[k]; ok {
			return nil, fmt.Errorf("key already exists")
		}

		v, ok := values[1].(To)

		if !ok {
			return nil, fmt.Errorf(
				"expected value of type %T, got %T", *new(To), values[1],
			)
		}

		if _, ok := m.from[v]; ok {
			return nil, fmt.Errorf("value already exists")
		}

		m.to[k] = v
		m.from[v] = k

		values = values[2:]
	}

	return m, nil
}

func MustNew[From, To comparable](
	values ...any,
) *Mapper[From, To] {
	m, err := New[From, To](values...)

	if err != nil {
		panic(err)
	}

	return m
}

func (m *Mapper[From, To]) To(from From) (To, error) {
	if to, ok := m.to[from]; ok {
		return to, nil
	}
	var zero To
	return zero, fmt.Errorf("no mapping found")
}

func (m *Mapper[From, To]) ToWithDefault(from From, def To) To {
	if to, ok := m.to[from]; ok {
		return to
	}
	return def
}

func (m *Mapper[From, To]) From(to To) (From, error) {
	if from, ok := m.from[to]; ok {
		return from, nil
	}
	var zero From
	return zero, fmt.Errorf("no mapping found")
}

func (m *Mapper[From, To]) FromWithDefault(to To, def From) From {
	if from, ok := m.from[to]; ok {
		return from
	}
	return def
}
