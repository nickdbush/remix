package remix

import (
	"go/types"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

type Fields struct {
	storage *orderedmap.OrderedMap[string, Field]
}

func (f Fields) Count() int {
	return f.storage.Len()
}

func (f Fields) Get(name string) (Field, bool) {
	return f.storage.Get(name)
}

func (f Fields) clone() Fields {
	clone := Fields{
		storage: orderedmap.New[string, Field](),
	}
	f.iter(func(name string, field Field) bool {
		clone.storage.Set(name, field)
		return true
	})
	return clone
}

func (f Fields) cloneMap(fn func(Field) *Field) Fields {
	clone := Fields{
		storage: orderedmap.New[string, Field](),
	}
	f.iter(func(name string, field Field) bool {
		if newField := fn(field); newField != nil {
			clone.storage.Set(newField.name, *newField)
		}
		return true
	})
	return clone
}

func (f Fields) iter(fn func(string, Field) bool) {
	for pair := f.storage.Oldest(); pair != nil; pair = pair.Next() {
		if !fn(pair.Key, pair.Value) {
			break
		}
	}
}

type Field struct {
	name   string
	ty     types.Type
	actual string
}

func (f Field) Name() string {
	return f.name
}

func (f Field) Type() types.Type {
	return f.ty
}
