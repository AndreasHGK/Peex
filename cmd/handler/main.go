package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/andreashgk/peex/eventid"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const template = `// Package %s
// This file was generated using the handler generator. Do not edit.
package %s

import (
	"github.com/andreashgk/peex/eventid"
	"github.com/zyedidia/generic/mapset"
)

// Events returns a list of all events that are implemented by the handler. Any event that is not part of this set will
// never be called. This method is usually generated using the handler tool.
func (%s) Events() mapset.Set[eventid.EventId] {
	set := mapset.New[eventid.EventId]()
	%s
	return set
}
`

func main() {
	outFile := flag.String("o", "", "output file for handler info")
	inputType := flag.String("handler", "", "the handler's type name")
	isPointer := flag.Bool("ptr", false, "whether the method receiver should be a pointer")
	flag.Parse()

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, "./", nil, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing package: %s\n", err.Error())
		os.Exit(1)
	}
	if len(pkgs) != 1 {
		fmt.Fprintf(os.Stderr, "expected exactly one package, got %v", len(pkgs))
		os.Exit(1)
	}

	methods := []*ast.FuncDecl{}
	pkgName := ""
	// First, find all methods on the type specified
	for _, pkg := range pkgs {
		pkgName = pkg.Name

		// Remove all non-exported nodes
		ast.PackageExports(pkg)

		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				f, _ := decl.(*ast.FuncDecl)
				// Filter all types that are not methods.
				if f == nil || f.Recv == nil {
					continue
				}

				recv := f.Recv.List[0]
				// Find the name of the receiver's type.
				var typeName string
				switch recvType := recv.Type.(type) {
				case *ast.StarExpr:
					if ptrTypr, ok := recvType.X.(*ast.Ident); ok {
						typeName = ptrTypr.Name
					}
				case *ast.Ident:
					typeName = recvType.Name
				default:
					panic("unexpected receiver type")
				}

				// We only want a receiver to be of the type we are generating a method for.
				if typeName != *inputType {
					continue
				}

				methods = append(methods, f)
			}
		}
	}

	var eventPuts []string
	// Go over each method to find which are part of the Handler interface, and which event they handle.
	for _, f := range methods {
		if f.Name == nil {
			panic("unexpected method name")
		}
		name := f.Name.Name
		_, ok := eventid.GetEventId(name)
		if !ok {
			// The method is not an event handler, skip this method.
			continue
		}

		eventPuts = append(eventPuts, "set.Put(eventid.Event"+strings.TrimPrefix(name, "Handle")+")")
	}

	if *isPointer {
		*inputType = "*" + *inputType
	}
	fileContent := fmt.Sprintf(template, pkgName, pkgName, *inputType, strings.Join(eventPuts, "\n\t"))
	err = os.WriteFile(*outFile, bytes.NewBufferString(fileContent).Bytes(), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %s\n", err.Error())
		os.Exit(1)
	}
}
