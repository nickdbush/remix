package remix

import (
	"fmt"
	"go/types"
	"os"

	orderedmap "github.com/wk8/go-ordered-map/v2"
	"golang.org/x/tools/go/packages"
)

var basePath = "."
var basePackage string

func SetBasePath(path string) {
	basePath = path
}

func SetBasePackage(pkg string) {
	basePackage = pkg
}

func packagePath(pkg string) string {
	if basePackage != "" {
		return basePackage + "/" + pkg
	}
	return pkg
}

type Package struct {
	pkg *types.Package
}

func (p Package) Name() string {
	return p.pkg.Name()
}

func (p Package) Get(name string) *Node {
	ty := p.pkg.Scope().Lookup(name)
	if ty == nil {
		panic(fmt.Sprintf("type %q not found in %q", name, p.pkg.Path()))
	}

	return importType(ty.Type())
}

func (p Package) Types() Types {
	return Types{raw: p.pkg}
}

func Import(pkg string) Package {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedTypes,
		Dir:  basePath,
	}, packagePath(pkg))
	if err != nil {
		panic(err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return Package{pkg: pkgs[0].Types}
}

func importType(ty types.Type) *Node {
	switch t := ty.(type) {
	case *types.Named:
		if s, ok := t.Underlying().(*types.Struct); ok {
			obj := t.Obj()
			return importStruct(obj.Type(), s)
		}
	}
	panic(fmt.Sprintf("unhandled type %T: %s", ty, ty))
}

func importStruct(ty types.Type, obj *types.Struct) *Node {
	fields := orderedmap.New[string, Field]()
	for i := 0; i < obj.NumFields(); i++ {
		field := obj.Field(i)
		if !field.Exported() {
			continue
		}

		fields.Set(field.Name(), Field{
			name:   field.Name(),
			actual: field.Name(),
			ty:     field.Type(),
		})
	}

	return load(ty, Fields{storage: fields})
}

func id(pkg string, name string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}
