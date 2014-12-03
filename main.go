package main

// TODO: hook result from execution into ui
// TOOD: add other data types

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

type field struct {
	Name string
	Type string
}

type function struct {
	Name    string
	Params  []field
	Results []field
	Doc     string
}

func filter(fi os.FileInfo) bool {
	if fi.Name() == "autoui.go" {
		return false
	}
	return true
}

func main() {
	log.SetFlags(log.Lshortfile)

	// find all the exported functions in the current package main

	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", filter, parser.ParseComments)
	if err != nil {
		log.Fatalln(err)
	}

	pkg := doc.New(pkgs["main"], "main", 0)

	functions := []function{}
	for _, f := range pkg.Funcs {
		// only exported functions
		// (allows helper functions without ui)
		if !ast.IsExported(f.Name) {
			continue
		}

		// only deal with fucntions, not methods
		if f.Decl.Recv != nil {
			continue
		}

		functions = append(functions, function{
			Name:    f.Name,
			Params:  fieldsFor(f.Decl.Type.Params),
			Results: fieldsFor(f.Decl.Type.Results),
			Doc:     f.Doc,
		})

		log.Println(f.Name)
	}

	// setup the template

	tmpl := template.New("qml.go")
	tmpl.Funcs(template.FuncMap{
		"id": func(id string) string {
			return strings.ToLower(id[:1]) + id[1:]
		},
		"zeroString": func(t string) string {
			switch t {
			case "int":
				return "0"
			case "error":
				return "<nil>"
			default:
				log.Println("unsupported type:", t)
				return ""
			}
		},
		"convert": func(name, t string) string {
			switch t {
			case "int":
				return fmt.Sprintf(`
					%s64, err := strconv.ParseInt(%[1]sStr, 10, 64)
					if err != nil {
						log.Println(err)
						return
					}
					%[1]s := int(%[1]s64)
				`, name)
			default:
				log.Println("unsupported type:", t)
				return ""
			}
		},
		"resultList": func(n int) string {
			if n == 0 {
				return ""
			}

			results := make([]string, n)
			for i := range results {
				results[i] = fmt.Sprintf("result%d", i)
			}

			return fmt.Sprintf("%s =", strings.Join(results, ", "))
		},
	})

	// find and parse the template

	b, err := rice.FindBox("templates")
	if err != nil {
		log.Fatalln(err)
	}

	str, err := b.String("qml.go")
	if err != nil {
		log.Fatalln(err)
	}

	tmpl, err = tmpl.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}

	// run the template

	w, err := os.Create("autoui.go")
	if err != nil {
		log.Fatalln(err)
	}
	defer w.Close()

	w.WriteString("// THIS FILE GENERATED BY github.com/nathankerr/autoui\n// DO NOT EDIT\n\n")

	err = tmpl.Execute(w, functions)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt the result
	err = exec.Command("go", "fmt", "autoui.go").Run()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("done")
}

func fieldsFor(list *ast.FieldList) []field {
	fields := make([]field, list.NumFields())
	fieldNum := 0

	for _, f := range list.List {
		t := f.Type.(*ast.Ident).Name

		// fields with names
		for _, ident := range f.Names {
			fields[fieldNum] = field{
				Name: ident.String(),
				Type: t,
			}
			fieldNum++
		}

		// handle anonymous fields
		if len(f.Names) == 0 {
			fields[fieldNum] = field{
				Type: t,
			}
			fieldNum++
		}
	}
	return fields
}
