package main

import (
	"log"
	"strconv"
	"fmt"
	"flag"
	"io/ioutil"

	"gopkg.in/qml.v1"
)

var inline = `import QtQuick 2.0
import QtQuick.Layouts 1.1
import QtQuick.Controls 1.2

TabView { {{ range . }}
	{{$function := .Name}}
	Tab {
		title: "{{.Name}}"

		ScrollView {
			GridLayout {
				anchors.centerIn: parent
				columns: 2

				Rectangle {
					Layout.columnSpan: 2
					height: 10
				}

				{{ range .Params }}Text {
					text: "{{.Name}}"
				}
				TextField {
					text: "{{zeroString .Type}}"
					onEditingFinished: ctrl.set("{{$function}}", "{{.Name}}", text)
				}
				{{ end }}

				Button {
					Layout.columnSpan: 2
					text: "Run {{.Name}}"
					onClicked: ctrl.run("{{.Name}}")
				}

				{{ range $i, $result := .Results }}Text {
					text: "{{$result.Name}}"
				}
				TextField {
					text: ctrl.{{id $function}}result{{$i}}
					readOnly: true
				}
				{{ end }}

				Rectangle{
					Layout.columnSpan: 2
					height: 10
				}
			}
		}
	}
{{ end }}}`

func main() {
	log.SetFlags(log.Lshortfile)

	export := flag.Bool("export", false, "export embedded qml to stdout")
	flag.Parse()

	if *export {
		fmt.Println(inline)
		return
	}

	// use a custom qml
	if flag.NArg() == 1 {
		log.Println(flag.Arg(0))

		b, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatalln(err)
		}
		inline = string(b)
	}

	err := qml.Run(run)
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	engine := qml.NewEngine()
	component, err := engine.LoadString("inline", inline)
	if err != nil {
		log.Fatalln(err)
	}

	ctrl := &Control{
		Params: make(map[string]map[string]string),
	}

	{{range .}}ctrl.Params["{{.Name}}"] = make(map[string]string)
	{{end}}

	engine.Context().SetVar("ctrl", ctrl)

	win := component.CreateWindow(nil)
	win.Show()
	win.Wait()

	return nil
}

type Control struct {
	Params map[string]map[string]string

	{{range .}}{{$function := .Name}}{{range $i, $result := .Results}}
	{{$function}}result{{$i}} string{{end}}{{end}}
}

func (ctrl *Control) Run(function string) {
	switch function { {{range .}}{{$function := .Name}}
	case "{{.Name}}": {{range .Params}}
		{{.Name}}Str := ctrl.Params["{{$function}}"]["{{.Name}}"]
		if {{.Name}}Str == "" {
			{{.Name}}Str = "{{zeroString .Type}}"
		}

		{{convert .Name .Type}}{{end}}

		var (
			{{range $i, $result := .Results}}result{{$i}} {{$result.Type}}
		{{end}})

		{{resultList (len .Results)}} {{$function}}({{range .Params}}{{.Name}},{{end}})

		{{range $i, $result := .Results}}
		ctrl.{{$function}}result{{$i}} = fmt.Sprint(result{{$i}})
		qml.Changed(ctrl, &ctrl.{{$function}}result{{$i}})
		{{end}}
	{{end}}
	default:
		log.Println("unhandled function", function)
	}
}

func (ctrl *Control) Set(function, param, value string) {
	ctrl.Params[function][param] = value
}
