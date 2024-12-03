package remix_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nickdbush/remix"
)

func TestType_Instantiate(t *testing.T) {
	pkg := remix.Import("github.com/nickdbush/remix/fixture")
	types := pkg.Types()
	generic := types.Get("Generic")
	assert.Equal(t, "Generic[T any]", generic.RelativeTo(pkg))
	wrapped := generic.Instantiate(types.Get("Concrete"))
	assert.Equal(t, "Generic[Concrete]", wrapped.RelativeTo(pkg))
	assert.Equal(t, "Generic[T any]", generic.RelativeTo(pkg))
}

func TestType_Slice(t *testing.T) {
	pkg := remix.Import("github.com/nickdbush/remix/fixture")
	ty := pkg.Types().Get("Concrete")
	assert.Equal(t, "[]Concrete", ty.Slice().RelativeTo(pkg))
	assert.Equal(t, "[][]Concrete", ty.Slice().Slice().RelativeTo(pkg))
	assert.Equal(t, "Concrete", ty.RelativeTo(pkg))
}

func TestType_WithKeys(t *testing.T) {
	pkg := remix.Import("github.com/nickdbush/remix/fixture")
	ty := pkg.Types().Get("Concrete")
	assert.Equal(t, "map[string]Concrete", ty.WithKeys(remix.String()).RelativeTo(pkg))
	assert.Equal(t, "Concrete", ty.RelativeTo(pkg))
}

func TestType_WithValues(t *testing.T) {
	pkg := remix.Import("github.com/nickdbush/remix/fixture")
	ty := pkg.Types().Get("Concrete")
	assert.Equal(t, "map[Concrete]string", ty.WithValues(remix.String()).RelativeTo(pkg))
	assert.Equal(t, "Concrete", ty.RelativeTo(pkg))
}
