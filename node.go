package remix

import (
	"fmt"
	"go/types"
)

type Operation interface {
	resolve() Fields
	parents() []*Node
	debug() string
}

type loadOp struct {
	ty     types.Type
	fields Fields
}

func load(ty types.Type, fields Fields) *Node {
	return wrap(loadOp{ty: ty, fields: fields})
}

func (op loadOp) resolve() Fields {
	return op.fields
}

func (op loadOp) parents() []*Node {
	return nil
}

func (op loadOp) debug() string {
	return types.TypeString(op.ty, nil)
}

func (op loadOp) Relative(qualifier types.Qualifier) string {
	return types.TypeString(op.ty, qualifier)
}

type withoutOp struct {
	parent *Node
	field  string
}

func (n *Node) Without(field string) *Node {
	return wrap(withoutOp{parent: n, field: field})
}

func (op withoutOp) resolve() Fields {
	return op.parent.resolve().cloneMap(func(field Field) *Field {
		if field.Name() == op.field {
			return nil
		}
		return &field
	})
}

func (op withoutOp) parents() []*Node {
	return []*Node{op.parent}
}

func (op withoutOp) debug() string {
	return fmt.Sprintf("without(%s)", op.field)
}

type renameOp struct {
	parent *Node
	from   string
	to     string
}

func (n *Node) Rename(from, to string) *Node {
	return wrap(renameOp{parent: n, from: from, to: to})
}

func (op renameOp) resolve() Fields {
	return op.parent.resolve().cloneMap(func(field Field) *Field {
		if field.Name() == op.from {
			return &Field{
				name:   op.to,
				ty:     field.Type(),
				actual: field.actual,
			}
		}
		return &field
	})
}

func (op renameOp) parents() []*Node {
	return []*Node{op.parent}
}

func (op renameOp) debug() string {
	return fmt.Sprintf("rename(%s, %s)", op.from, op.to)
}

var nextNodeID int

type Node struct {
	Operation
	id int
}

func wrap(op Operation) *Node {
	node := &Node{
		Operation: op,
		id:        nextNodeID,
	}
	nextNodeID++
	return node
}
