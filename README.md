# remix those go types

a research project to explore the possibilities of go type manipulation

this will become a usable library at some point, promise :)

## fluid type api

```go
package main

import "github.com/nickdbush/remix"

func main() {
	pkg := remix.Import("github.com/your/package/here")
	types := pkg.Types()
	types.Get("Option"). // Option[T any]
		Instantiate(types.Get("User")). // Option[User]
		WithKeys(remix.String()) // map[string]Option[User]
}
```

## struct conversion
 
given these types...

```go
package user

type Params struct {
	Name  string
	Email fp.Option[string]
	Years int
}

type Model struct {
	ID    string
	Name  string
	Email fp.Option[string]
	Age   int
}
```

and this config...

```go
package main

import "github.com/nickdbush/remix"

func main() {
	remix.SetBasePath(".")
	remix.SetBasePackage("github.com/your/rootpkg")

	pkg := remix.Import("user")

	out := remix.ToSource(pkg, "user/convert.go")

	model := pkg.Get("Model")
	params := pkg.Get("Params").Rename("Years", "Age")
	
	// If you want to output the types in the generated source.
	// out.Add("User", model)
	// out.Add("UserParams", params)
	
	out.Convert("FromParams", params, model, nil)

	out.Finish()
}
```

the following code will be generated

```go
package user

type FromParamsFields struct {
	ID string
}

func FromParams(from Params, extra FromParamsFields) Model {
	return Model{
		Name:  from.Name,
		Email: from.Email,
		Age:   from.Age,
		ID:    extra.ID,
	}
}

```