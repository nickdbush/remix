package main

import (
	"go/types"

	"github.com/nickdbush/remix"
)

func main() {
	remix.SetBasePath(".")
	remix.SetBasePackage("github.com/nickdbush/remix")
	pkg := remix.Import("bin")

	out := remix.ToSource(pkg, "bin/convert.txt")

	model := pkg.Get("User")
	params := pkg.Get("UserParams").Rename("Id", "ID")

	out.Add("User", model)
	out.Add("UserParams", params)
	out.Convert("ProtoToModel", params, model, casts())

	out.Finish()
}

func casts() *remix.Casts {
	casts := remix.NewCasts()
	// Converts a generic main.Option[T] to a T, falling back to the default value
	casts.Add(func(t types.Type) *remix.Cast {
		if named, ok := t.(*types.Named); ok {
			if named.Obj().Pkg().Path() == "main" && named.Obj().Name() == "Option" {
				return &remix.Cast{
					Result: named.Underlying().(*types.Struct).Field(0).Type(),
					Source: func(v string) string {
						return v + ".Value"
					},
				}
			}
		}
		return nil
	})
	return casts
}
