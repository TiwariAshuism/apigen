package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/TiwariAshuism/apigen/pkg/model"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// allLayers is the default set of layers produced when -layers is not specified.
var allLayers = []layer{
	{"handler.tmpl", "handler", "_handler_gen.go", false},
	{"service.tmpl", "service", "_service_gen.go", false},
	{"repository.tmpl", "repository", "_repository_gen.go", false},
	{"openapi.tmpl", ".", "_openapi.yaml", true}, // opt-in only
}

type layer struct {
	tmpl       string
	dir        string // relative to outputDir/internal/, or "." for root
	fileSuffix string
	optIn      bool // only rendered when explicitly requested in Layers
}

// Generator renders MVC layers for a set of APIDefinitions.
type Generator struct {
	OutputDir string
	Module    string
	// Layers is the set of layer names to render, e.g. ["handler","service","repository","openapi"].
	// Defaults to all non-opt-in layers when empty.
	Layers []string
	// DryRun prints what would be written without touching the filesystem.
	DryRun    bool
	templates *template.Template
}

// New creates a Generator that writes output under outputDir using the given module path.
func New(outputDir, module string) (*Generator, error) {
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
		"title": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},
		"groupByPath": groupByPath,
	}

	tmpl, err := template.New("").Funcs(funcMap).ParseFS(templateFS, "templates/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("loading templates: %w", err)
	}

	return &Generator{OutputDir: outputDir, Module: module, templates: tmpl}, nil
}

// Generate writes the selected layers for each definition.
func (g *Generator) Generate(defs []model.APIDefinition) error {
	active := g.activeLayers()

	var errs []error
	written, unchanged := 0, 0

	for _, def := range defs {
		def.Module = g.Module
		if def.ModelImportPath == "" {
			def.ModelImportPath = g.Module + "/internal/model"
		}
		if def.ModelAlias == "" {
			def.ModelAlias = "model"
		}

		for _, l := range active {
			var outDir, outPath string
			if l.dir == "." {
				outDir = g.OutputDir
			} else {
				outDir = filepath.Join(g.OutputDir, "internal", l.dir)
			}
			outPath = filepath.Join(outDir, strings.ToLower(def.Resource)+l.fileSuffix)

			// Render template into an in-memory buffer.
			var buf bytes.Buffer
			if err := g.templates.ExecuteTemplate(&buf, l.tmpl, def); err != nil {
				errs = append(errs, fmt.Errorf("render %s for %s: %w", l.tmpl, def.Name, err))
				continue
			}

			// Apply gofmt for .go files.
			content := buf.Bytes()
			if strings.HasSuffix(outPath, ".go") {
				formatted, err := format.Source(content)
				if err != nil {
					errs = append(errs, fmt.Errorf("gofmt %s: %w (writing unformatted)", outPath, err))
				} else {
					content = formatted
				}
			}

			if g.DryRun {
				fmt.Printf("=== [dry-run] %s ===\n%s\n", outPath, content)
				continue
			}

			// Overwrite guard: skip if content is identical to existing file.
			if existing, err := os.ReadFile(outPath); err == nil && bytes.Equal(existing, content) {
				unchanged++
				continue
			}

			if err := os.MkdirAll(outDir, 0755); err != nil {
				errs = append(errs, fmt.Errorf("mkdir %s: %w", outDir, err))
				continue
			}
			if err := os.WriteFile(outPath, content, 0644); err != nil {
				errs = append(errs, fmt.Errorf("write %s: %w", outPath, err))
				continue
			}
			written++
		}
	}

	if !g.DryRun {
		fmt.Printf("apigen: %d file(s) written, %d unchanged\n", written, unchanged)
	}

	if len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		return fmt.Errorf("generation errors:\n%s", strings.Join(msgs, "\n"))
	}
	return nil
}

// activeLayers returns the layers to render based on g.Layers.
// If g.Layers is empty, all non-opt-in layers are returned.
func (g *Generator) activeLayers() []layer {
	if len(g.Layers) == 0 {
		var out []layer
		for _, l := range allLayers {
			if !l.optIn {
				out = append(out, l)
			}
		}
		return out
	}

	set := make(map[string]bool, len(g.Layers))
	for _, name := range g.Layers {
		set[strings.TrimSpace(name)] = true
	}

	var out []layer
	for _, l := range allLayers {
		if set[l.dir] || set[strings.TrimSuffix(l.tmpl, ".tmpl")] {
			out = append(out, l)
		}
	}
	return out
}

// groupByPath groups methods by their URL path, preserving insertion order of paths.
// Used by the OpenAPI template so each path key appears only once.
func groupByPath(methods []model.Method) []pathGroup {
	seen := map[string]int{}
	var groups []pathGroup
	for _, m := range methods {
		if i, ok := seen[m.Path]; ok {
			groups[i].Methods = append(groups[i].Methods, m)
		} else {
			seen[m.Path] = len(groups)
			groups = append(groups, pathGroup{Path: m.Path, Methods: []model.Method{m}})
		}
	}
	return groups
}

// pathGroup is a template-visible grouping of methods sharing the same URL path.
type pathGroup struct {
	Path    string
	Methods []model.Method
}
