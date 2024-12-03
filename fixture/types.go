package fixture

type Generic[T any] struct {
	Value T
}

type Concrete struct {
	Value bool
}
