package model

// APIDefinition is the parsed result of one interface.
type APIDefinition struct {
	PackageName     string
	Name            string // e.g. "UserAPI"
	Resource        string // e.g. "User"
	Module          string // injected by generator (Go module path)
	ModelImportPath string // e.g. "github.com/foo/internal/types"
	ModelAlias      string // e.g. "types" — import alias for domain types
	Methods         []Method
}

// Method represents a single interface method mapped to an HTTP endpoint.
type Method struct {
	Name       string   // e.g. "ListUsers"
	HTTPMethod string   // GET / POST / PUT / DELETE / PATCH
	Path       string   // e.g. /users or /users/:id
	Request    *Param   // input struct param (nil if none)
	Response   *Param   // output type (nil if none)
	PathParams []string // list of path param names (e.g. ["id"])
	Args       []Argument
}

type Argument struct {
	Kind ArgumentKind
	Name string
	Type string // only for request kind
}

type ArgumentKind string

const (
	ArgRequest   ArgumentKind = "request"
	ArgPathParam ArgumentKind = "pathParam"
)

// Param describes a single Go parameter or return value.
type Param struct {
	Name    string // Go variable name
	Type    string // Go type string (without leading [])
	IsSlice bool
}
