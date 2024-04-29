package modelchecker

import (
	"fmt"
	"github.com/fizzbee-io/fizzbee/lib"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var (
	roleMethods = map[string]*starlark.Builtin{
		//"clear": starlark.NewBuiltin("clear", dict_clear),
	}
)

var (
	roleRefs = map[string]int{}
)
type Role struct {
	ref int
	Name string
	Params *starlarkstruct.Struct
	Fields *lib.Struct
}

func (r *Role) SetField(name string, val starlark.Value) error {
	// If name is found in Params or BuiltinAttrNames, then fail
	// Otherwise call Fields.SetField
	if _, err := r.Params.Attr(name); err == nil {
		return fmt.Errorf("cannot set immutable field %s on role %s", name, r.Name)
	} else if _, ok := err.(starlark.NoSuchAttrError); !ok {
		return err
	} else if v, _ := lib.BuiltinAttr(r, name, roleMethods); v != nil {
		return fmt.Errorf("cannot override builtins %s on role %s", name, r.Name)
	}
	return r.Fields.SetField(name, val)
}

func (r *Role) Attr(name string) (starlark.Value, error) {
	if v, err := r.Fields.Attr(name); err == nil {
		return v, nil
	} else if _, ok := err.(starlark.NoSuchAttrError); !ok {
		return v, err
	} else if v, err := r.Params.Attr(name); err == nil{
		return v, nil
	} else if _, ok := err.(starlark.NoSuchAttrError); !ok {
		return v, err
	}
	return lib.BuiltinAttr(r, name, roleMethods)
}

func (r *Role) AttrNames() []string {
	return lib.BuiltinAttrNames(roleMethods)
}

func (r *Role) String() string {
	return fmt.Sprintf("role %s#%d (%s,%s)", r.Name, r.ref, r.Params.String(), r.Fields.String())
}

func (r *Role) Type() string {
	return "role"
}

func (r *Role) Freeze() {

}

func (r *Role) Truth() starlark.Bool {
	return true
}

func (r *Role) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: role")
}

var _ starlark.HasAttrs = (*Role)(nil)
var _ starlark.HasSetField = (*Role)(nil)
var _ starlark.Value = (*Role)(nil)

func CreateRoleBuiltin(name string) *starlark.Builtin {
	return starlark.NewBuiltin(name, func(t *starlark.Thread, b *starlark.Builtin,
		args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		params := starlarkstruct.FromKeywords(starlark.String("params"), kwargs)
		nextRef := roleRefs[name]
		if roleRefs[name] > 0 {
			roleRefs[name]++
		} else {
			roleRefs[name] = 1
		}
		fields := lib.FromStringDict(starlark.String("fields"), starlark.StringDict{})
		return &Role{ref: nextRef+1, Name: name, Params: params, Fields: fields}, nil
	})
}