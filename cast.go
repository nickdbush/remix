package remix

import "go/types"

type Cast struct {
	Result types.Type
	Source func(v string) string
}

type CastFn func(types.Type) *Cast

type Casts struct {
	fns []CastFn
}

func NewCasts() *Casts {
	return &Casts{}
}

func (c *Casts) Add(cast CastFn) *Casts {
	c.fns = append(c.fns, cast)
	return c
}

func (c *Casts) apply(ty types.Type) *Cast {
	var convertible *Cast
	for _, cast := range c.fns {
		if result := cast(ty); result != nil {
			if types.AssignableTo(ty, result.Result) {
				return result
			}
			if types.ConvertibleTo(ty, result.Result) {
				convertible = result
			}
		}
	}
	return convertible
}
