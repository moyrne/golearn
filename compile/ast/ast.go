package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

var (
	filename, structName string
)

func init() {
	flag.StringVar(&filename, "file", "", "generate file name")
	flag.StringVar(&structName, "struct", "", "generate struct name")
	flag.Parse()
	fmt.Println("[flag]", filename, structName)
}

func main() {
	// TODO 数据导出
	result, err := getStruct(filename, structName)
	if err != nil {
		log.Panicln(err)
	}
	if err := generateAntdTitle(filename, result); err != nil {
		log.Panicln(err)
	}
}

type StructResult struct {
	Name   string        `json:"name"`
	Fields []StructField `json:"fields"`
}

type StructField struct {
	Name       string `json:"name"`
	ArrayCount int    `json:"array_count"`
	Type       string `json:"type"`
	Tags       string `json:"tags"`
	Omitempty  bool   `json:"omitempty"`
	Comment    string `json:"comment"`
}

func getStruct(filename, structName string) (*StructResult, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// 初始化 scanner
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO 如何获取注释
	var result StructResult
	ast.Inspect(f, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.GenDecl:
			if node.Tok != token.TYPE {
				return false
			}
			for _, spec := range node.Specs {
				typeSpec := spec.(*ast.TypeSpec)
				name := typeSpec.Name.Name
				if name != structName {
					return true
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					err = errors.WithStack(errors.New("typeSpec.Type is not *ast.StructType"))
					return false
				}
				result.Name = name
				for _, field := range structType.Fields.List {
					// TODO 支持多 Names
					sf := StructField{
						Name: field.Names[0].Name,
					}

					fty := field.Type
					for {
						switch ty := fty.(type) {
						case *ast.Ident:
							sf.Type = ty.Name
							goto listout
						case *ast.ArrayType:
							fty = ty.Elt
							sf.ArrayCount++
						default:
							err = errors.WithStack(errors.New("field.Type is not in *ast.Ident and *ast.ArrayType"))
							return false
						}
					}
				listout:
					if field.Tag != nil {
						tag := field.Tag.Value
						if len(tag) > 2 {
							tag = tag[1 : len(tag)-1]
						}
						sf.Tags = tag
					}

					result.Fields = append(result.Fields, sf)
				}
			}
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func generateAntdTitle(filename string, result *StructResult) error {
	tpl, err := template.New("antd").Parse(`let {{.Name}} = [
	{{range $i,$field := .Fields}}{id:"{{$field.Name}}",title:"{{$field.Name}}",key:"{{$field.Name}}"},
	{{end}}
]`)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(filename[:len(filename)-3] + "_" + result.Name + ".js")
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	return errors.WithStack(tpl.Execute(f, result))
}
