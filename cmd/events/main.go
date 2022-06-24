package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
)

const (
	Package   = "/pkg/mod/github.com/df-mc/dragonfly@" // note: path does not include DF version!
	File      = "/server/player/handler.go"
	Interface = "Handler"
)

// very messy code ahead, read at your own risk
func main() {
	optOutFile := flag.String("o", "", "output file")
	optPkgName := flag.String("p", "", "the package name")
	goModFile := flag.String("m", "", "go.mod file")
	flag.Parse()

	// Read the file containing the handler interface and store it as a string.
	var sourceContent string
	{
		// Find the GOPATH
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}

		if *goModFile == "" {
			fmt.Fprintf(os.Stderr, "must provide a go.mod file")
			os.Exit(1)
		}
		// Read the dragonfly version in the go.mod file
		goMod, err := os.Open(*goModFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening go modules file: %s\n", err.Error())
			os.Exit(1)
		}
		goModBuf := bytes.Buffer{}
		goModBuf.ReadFrom(goMod)

		matches := regexp.MustCompile("github\\.com/df-mc/dragonfly (v.+)\\n").FindStringSubmatch(goModBuf.String())
		if len(matches) < 2 {
			fmt.Fprintf(os.Stderr, "could not find dragonfly version in specified go modules file. is it installed?\n")
			os.Exit(1)
		}
		// Open the actual file in dragonfly
		sourceFile, err := os.Open(gopath + Package + matches[1] + File)
		defer sourceFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening handler file: %s\n", err.Error())
			os.Exit(1)
		}
		buff := bytes.Buffer{}
		buff.ReadFrom(sourceFile)

		sourceContent = buff.String()
	}

	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "", sourceContent, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing file: %s\n", err.Error())
		os.Exit(1)
	}

	handlerInterface := func() *ast.InterfaceType {
		for _, decl := range file.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, spec := range gen.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if typ.Name.Name != Interface {
					continue
				}
				i, ok := typ.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				return i
			}
		}
		return nil
	}()
	if handlerInterface == nil {
		fmt.Fprintf(os.Stderr, "could not find %s interface", Interface)
		os.Exit(1)
	}

	imports := ""
	for _, decl := range file.Decls {
		imp, ok := decl.(*ast.GenDecl)
		if !ok || imp.Tok != token.IMPORT {
			continue
		}
		imports += "\n" + sourceContent[imp.Pos()-1:imp.End()-1]
		break
	}
	imports = strings.TrimPrefix(imports, "\n")

	eventIds := ""
	allEvents := ""
	handlers := ""
	interfaces := ""
	getHandlerEvents := ""
	for _, event := range handlerInterface.Methods.List {
		name := event.Names[0].Name
		eventName := "event" + strings.TrimPrefix(name, "Handle")
		interfaceName := eventName + "Handler"

		method, ok := event.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		if eventIds == "" {
			eventIds = eventName + " eventId = iota"
		} else {
			eventIds += "\n\t" + eventName
		}

		allEvents += "\n\t" + eventName + ","
		funcParams := sourceContent[method.Pos()-1 : method.End()-1]
		funcDecl := name + funcParams
		interfaces += "\n" + fmt.Sprintf(interfaceTemplate, interfaceName, funcDecl)
		getHandlerEvents += fmt.Sprintf(getHandlerEventsTemplate, interfaceName, eventName)

		eventArgs := ""
		for _, p := range method.Params.List {
			for _, n := range p.Names {
				eventArgs += ", " + n.Name
			}
		}
		eventArgs = strings.TrimPrefix(eventArgs, ", ")

		methodBody := fmt.Sprintf(methodBodyTemplate, eventName, interfaceName, name, eventArgs)
		if eventName == "eventQuit" {
			methodBody += "\n\ts.doQuit()" // to run the session specific logic when the player quits
		}
		handlers += "\n\n" + fmt.Sprintf(methodTemplate, funcDecl, methodBody)
	}

	pkgName := *optPkgName
	filledTemplate := fmt.Sprintf(template,
		pkgName,
		pkgName,
		imports,
		eventIds,
		getHandlerEvents,
		strings.TrimPrefix(allEvents, "\n\t"),
		strings.TrimPrefix(interfaces, "\n\n"),
		strings.TrimPrefix(handlers, "\n\n"),
	)

	err = os.WriteFile(*optOutFile, bytes.NewBufferString(filledTemplate).Bytes(), os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing output file: %s\n", err.Error())
		os.Exit(1)
	}
}

const template = `// Package %s
// This file was generated using the event generator. Do not edit.
package %s

%s

type eventId uint

const (
	%s
)

// getHandlerEvents returns which events a handler implements. Since it is impossible to distinguish actually imlemented
// methods from ones embedded using player.NopHandler, it is recommended to not embed it at all. Most Peex handlers 
// won't implement player.Handler!
func getHandlerEvents(h Handler) map[eventId]struct{} {
	m := make(map[eventId]struct{})
	%s
	return m
}

var allEvents = []eventId {
	%s
}

%s

%s
`

const methodTemplate = `func (s *Session) %s {
	%s
}`

const methodBodyTemplate = `s.handleEvent(%s, func(h Handler) {
		h.(%s).%s(%s)
	})`

// interfaceTemplate is a template used for creating interfaces for event handlers.
const interfaceTemplate = `
type %s interface {
	%s
}`

// getHandlerEventsTemplate is a template used for the getHandlerEvents function, to detect
// what events a handler implements.
const getHandlerEventsTemplate = `
	if _, ok := h.(%s); ok {
		m[%s] = struct{}{}
	}`
