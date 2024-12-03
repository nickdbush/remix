package remix

import (
	"fmt"
	"go/types"
)

type Conversion struct {
	From      *Node
	To        *Node
	Qualifier types.Qualifier
	Casts     *Casts
}

func (c Conversion) source(name string) string {
	fromFields := c.From.resolve()
	fromLoad := c.From.findNearestLoad()
	if fromLoad == nil {
		panic("fromLoad is nil")
	}

	toLoad, isToLoad := c.To.Operation.(loadOp)
	if !isToLoad {
		panic("toLoad is not a loadOp")
	}

	fromPath := fromLoad.Relative(c.Qualifier)
	toPath := toLoad.Relative(c.Qualifier)

	var unknownFields []unknownField
	var knownFields []knownField
	toLoad.fields.iter(func(_ string, field Field) bool {
		fromField, isKnown := fromFields.Get(field.name)
		if !isKnown {
			unknown := newUnknownField(field)
			unknownFields = append(unknownFields, unknown)
			return true
		}

		variable := fmt.Sprintf("from.%s", fromField.actual)
		if types.AssignableTo(fromField.ty, field.ty) {
			// no-op
		} else if types.ConvertibleTo(fromField.ty, field.ty) {
			variable = fmt.Sprintf("%s(%s)", types.TypeString(field.ty, c.Qualifier), variable)
		} else {
			panic(fmt.Sprintf("cannot convert %s to %s", fromField.ty, field.ty))
		}

		knownFields = append(knownFields, newKnownField(field, variable))
		return true
	})

	out := ""
	extra := ""
	if len(unknownFields) > 0 {
		extra = fmt.Sprintf("%sFields", name)
		out += fmt.Sprintf("type %s struct {\n", extra)
		for _, unknown := range unknownFields {
			out += fmt.Sprintf("%s %s\n", unknown.alias, types.TypeString(unknown.field.ty, c.Qualifier))
		}
		out += "}\n"
	}

	if extra != "" {
		out += fmt.Sprintf("func %s(from %s, extra %s) %s {\n", name, fromPath, extra, toPath)
	} else {
		out += fmt.Sprintf("func %s(from %s) %s {\n", name, fromPath, toPath)
	}
	out += fmt.Sprintf("return %s{\n", toPath)

	for _, known := range knownFields {
		out += fmt.Sprintf("%s: %s,\n", known.field.name, known.variable)
	}
	for _, unknown := range unknownFields {
		out += fmt.Sprintf("%s: extra.%s,\n", unknown.field.name, unknown.alias)
	}

	out += "}\n}\n"

	return out
}

type knownField struct {
	field    Field
	variable string
}

func newKnownField(field Field, source string) knownField {
	return knownField{
		field:    field,
		variable: source,
	}
}

type unknownField struct {
	field Field
	alias string
}

func newUnknownField(field Field) unknownField {
	return unknownField{
		field: field,
		alias: field.name,
	}
}
