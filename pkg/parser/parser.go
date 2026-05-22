package parser

import (
	"fmt"
	"go/ast"
	gparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/TiwariAshuism/apigen/pkg/model"
)

// ParseFile reads a single Go source file and returns all interface definitions found,
// along with any validation warnings (e.g. methods missing HTTP annotations).
func ParseFile(filePath string) ([]model.APIDefinition, []string, error) {
	fset := token.NewFileSet()
	f, err := gparser.ParseFile(fset, filePath, nil, gparser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return parseASTFile(f)
}

// ParseDir parses all non-generated Go files in a directory and merges their
// interface definitions. Generated files (containing "_gen.go") are skipped.
func ParseDir(dirPath string) ([]model.APIDefinition, []string, error) {
	fset := token.NewFileSet()
	pkgs, err := gparser.ParseDir(fset, dirPath, func(fi os.FileInfo) bool {
		return !strings.Contains(fi.Name(), "_gen.go")
	}, gparser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("parse dir %s: %w", dirPath, err)
	}

	var (
		allDefs     []model.APIDefinition
		allWarnings []string
	)
	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			defs, warns, err := parseASTFile(f)
			if err != nil {
				return nil, nil, err
			}
			allDefs = append(allDefs, defs...)
			allWarnings = append(allWarnings, warns...)
		}
	}
	return allDefs, allWarnings, nil
}

// parseASTFile extracts APIDefinitions from a single parsed *ast.File.
func parseASTFile(f *ast.File) ([]model.APIDefinition, []string, error) {
	var (
		defs     []model.APIDefinition
		warnings []string
	)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			iface, ok := typeSpec.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			modelImport, modelAlias := detectModelImport(f)

			def := model.APIDefinition{
				PackageName:     f.Name.Name,
				Name:            typeSpec.Name.Name,
				Resource:        extractResource(typeSpec.Name.Name),
				ModelImportPath: modelImport,
				ModelAlias:      modelAlias,
			}

			for _, field := range iface.Methods.List {
				if len(field.Names) == 0 {
					continue // embedded interface — skip
				}
				funcType, ok := field.Type.(*ast.FuncType)
				if !ok {
					continue
				}

				m := model.Method{
					Name: field.Names[0].Name,
				}

				// Parse HTTP method + path from the doc comment preceding the method.
				// e.g. "// GET /users/:id"
				if field.Doc != nil {
					for _, c := range field.Doc.List {
						parseHTTPComment(c.Text, &m)
					}
				}

				// Validate: skip methods with no HTTP annotation and emit a warning.
				if m.HTTPMethod == "" || m.Path == "" {
					warnings = append(warnings, fmt.Sprintf(
						"%s.%s: no HTTP comment found, skipping (expected e.g. \"// GET /path\")",
						def.Name, m.Name,
					))
					continue
				}

				// Extract path params first so parseParams can cross-reference them.
				m.PathParams = extractPathParams(m.Path)
				m.Request, m.Response, m.Args = parseParams(
					funcType, m.PathParams, modelAlias,
				)

				def.Methods = append(def.Methods, m)
			}

			// Only include interfaces that have at least one valid method.
			if len(def.Methods) > 0 {
				defs = append(defs, def)
			}
		}
	}

	return defs, warnings, nil
}

// parseHTTPComment extracts verb and path from "// GET /users/:id".
func parseHTTPComment(text string, m *model.Method) {
	text = strings.TrimPrefix(text, "//")
	text = strings.TrimSpace(text)
	parts := strings.SplitN(text, " ", 2)
	if len(parts) == 2 {
		verb := strings.ToUpper(parts[0])
		switch verb {
		case "GET", "POST", "PUT", "DELETE", "PATCH":
			m.HTTPMethod = verb
			m.Path = strings.TrimSpace(parts[1])
		}
	}
}

// detectModelImport finds the import used for domain types (internal/types or internal/model).
func detectModelImport(f *ast.File) (importPath, alias string) {
	alias = "model"
	for _, decl := range f.Imports {
		if decl == nil {
			continue
		}
		path := strings.Trim(decl.Path.Value, `"`)
		if !strings.HasSuffix(path, "/internal/types") &&
			!strings.HasSuffix(path, "/internal/model") {
			continue
		}
		importPath = path
		if decl.Name != nil {
			alias = decl.Name.Name
		} else {
			alias = path[strings.LastIndex(path, "/")+1:]
		}
		return importPath, alias
	}
	return "", alias
}

// typeExpr returns a valid Go type reference for generated code.
func typeExpr(typeStr, modelAlias string) string {
	if strings.Contains(typeStr, ".") {
		return typeStr
	}
	if modelAlias == "" {
		return typeStr
	}
	return modelAlias + "." + typeStr
}

// parseParams extracts the request body param and response type from a method signature.
// It skips context.Context (first param) and any params whose names match path params.
func parseParams(
	fn *ast.FuncType,
	pathParams []string,
	modelAlias string,
) (req *model.Param, resp *model.Param, args []model.Argument) {
	pathParamSet := make(map[string]bool, len(pathParams))
	for _, p := range pathParams {
		pathParamSet[p] = true
	}

	if fn.Params != nil {
		for i, p := range fn.Params.List {
			if i == 0 {
				continue // skip context.Context
			}

			name := ""
			if len(p.Names) > 0 {
				name = p.Names[0].Name
			}

			// Check if this is a path param
			if pathParamSet[name] {
				args = append(args, model.Argument{Kind: model.ArgPathParam, Name: name})
				continue
			}

			// If not a path param, assume it's the request object
			typeStr := exprToString(p.Type)
			req = &model.Param{
				Name: name,
				Type: typeExpr(typeStr, modelAlias),
			}
			args = append(args, model.Argument{
				Kind: model.ArgRequest, Name: name, Type: req.Type,
			})
		}
	}

	if fn.Results != nil {
		for _, r := range fn.Results.List {
			t := exprToString(r.Type)
			if t == "error" {
				continue
			}
			isSlice := strings.HasPrefix(t, "[]")
			base := strings.TrimPrefix(t, "[]")
			resp = &model.Param{
				Type:    typeExpr(base, modelAlias),
				IsSlice: isSlice,
			}
			break
		}
	}
	return
}

// exprToString converts an AST expression to its Go source representation.
func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	}
	return ""
}

// extractPathParams returns the list of named segments like ":id" → "id".
func extractPathParams(path string) []string {
	var params []string
	for _, seg := range strings.Split(path, "/") {
		if strings.HasPrefix(seg, ":") {
			params = append(params, strings.TrimPrefix(seg, ":"))
		}
	}
	return params
}

// extractResource strips the "API" suffix: "UserAPI" → "User".
func extractResource(name string) string {
	return strings.TrimSuffix(name, "API")
}

// IsDir reports whether path is a directory.
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// AbsPath returns the absolute form of path, panicking on error (for CLI use).
func AbsPath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
