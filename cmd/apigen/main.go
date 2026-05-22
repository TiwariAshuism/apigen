package main

import (
	"flag"
	"log"
	"strings"

	"github.com/TiwariAshuism/apigen/pkg/generator"
	"github.com/TiwariAshuism/apigen/pkg/model"
	"github.com/TiwariAshuism/apigen/pkg/parser"
)

func main() {
	input  := flag.String("input",   "api/routes.go", "Path to a Go file or directory")
	output := flag.String("output",  ".",             "Root output directory")
	module := flag.String("module",  "",              "Go module path (auto-detected from go.mod if empty)")
	layers := flag.String("layers",  "",              "Comma-separated layers: handler,service,repository,openapi (default: all)")
	dryRun := flag.Bool("dry-run", false, "Print output without writing files")
	force  := flag.Bool("force",   false, "Write files even when content is unchanged")
	flag.Parse()

	mod := *module
	if mod == "" {
		var err error
		mod, err = parser.ReadModulePath(*output)
		if err != nil {
			log.Fatalf("could not detect module path from go.mod: %v\n(use -module to set it explicitly)", err)
		}
		log.Printf("detected module: %s", mod)
	}

	var (
		defs     []model.APIDefinition
		warnings []string
		err      error
	)
	if parser.IsDir(*input) {
		defs, warnings, err = parser.ParseDir(*input)
	} else {
		defs, warnings, err = parser.ParseFile(*input)
	}
	if err != nil {
		log.Fatalf("parse %s: %v", *input, err)
	}
	for _, w := range warnings {
		log.Printf("WARNING: %s", w)
	}
	if len(defs) == 0 {
		log.Fatalf("no interface definitions found in %s", *input)
	}

	gen, err := generator.New(*output, mod)
	if err != nil {
		log.Fatalf("init generator: %v", err)
	}
	gen.DryRun = *dryRun
	gen.Force = *force
	if *layers != "" {
		gen.Layers = strings.Split(*layers, ",")
	}

	if err := gen.Generate(defs); err != nil {
		log.Fatalf("generate: %v", err)
	}
}
