package remix

import (
	"go/types"
)

type Types struct {
	raw *types.Package
}

type Type struct {
	raw types.Type
}

func String() Type {
	return Type{types.Typ[types.String]}
}

func NewType(ty types.Type) Type {
	return Type{ty}
}

func (t Types) Get(name string) Type {
	ty := t.raw.Scope().Lookup(name)
	if ty == nil {
		panic("type not found")
	}

	return NewType(ty.Type())
}

func (t Type) Instantiate(generics ...Type) Type {
	if named, ok := t.raw.(*types.Named); ok {
		if len(generics) != named.TypeParams().Len() {
			panic("wrong number of generics")
		}

		targs := make([]types.Type, 0, len(generics))
		for _, g := range generics {
			targs = append(targs, g.raw)
		}
		result, err := types.Instantiate(nil, named, targs, true)
		if err != nil {
			panic(err)
		}
		return NewType(result)
	}
	panic("not a named type")
}

func (t Type) Slice() Type {
	return NewType(types.NewSlice(t.raw))
}

func (t Type) WithKeys(key Type) Type {
	return NewType(types.NewMap(key.raw, t.raw))
}

func (t Type) WithValues(value Type) Type {
	return NewType(types.NewMap(t.raw, value.raw))
}

func (t Type) String() string {
	return types.TypeString(t.raw, nil)
}

func (t Type) RelativeTo(pkg Package) string {
	return types.TypeString(t.raw, types.RelativeTo(pkg.pkg))
}
