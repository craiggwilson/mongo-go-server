package mrpc

import (
	"fmt"
	"html/template"
	"io"
	"strings"

	"github.com/craiggwilson/mongo-go-server/mrpc/tree"
)

type generateConfig struct {
}

type generateOpt func(*generateConfig)

func Generate(w io.Writer, t *tree.Tree, opts ...generateOpt) error {
	cfg := &generateConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	g := generator{
		t: t,
		w: w,
	}

	return g.generate()
}

type generator struct {
	t *tree.Tree
	w io.Writer
}

func (g *generator) generate() error {
	tmpl := template.Must(template.New("mux").Funcs(template.FuncMap{
		"commandName": func(cmd *tree.Command) string {
			return strings.Title(cmd.Name)
		},
		"fieldName": func(f *tree.Field) string {
			if len(f.Name) < 3 {
				return strings.ToUpper(f.Name)
			}
			return strings.Title(f.Name)
		},
		"handlerName": func(cmd *tree.Command) string {
			return fmt.Sprintf("%sCommandHandler", strings.Title(cmd.Name))
		},
		"handlerImplName": func(cmd *tree.Command) string {
			return fmt.Sprintf("%sCommandHandlerImpl", cmd.Name)
		},
		"requestName": func(cmd *tree.Command) string {
			if cmd.Request.Name != "" {
				return cmd.Request.Name
			}

			return strings.Title(cmd.Name + "Request")
		},
		"responseName": func(cmd *tree.Command) string {
			if cmd.Response.Name != "" {
				return cmd.Response.Name
			}

			return strings.Title(cmd.Name + "Response")
		},
		"serviceName": func(svc *tree.Service) string {
			return fmt.Sprintf("%sService", strings.Title(svc.Name))
		},
		"serviceCommandHandlerName": func(cmdName string) string {
			return fmt.Sprintf("%sCommandHandler", strings.Title(cmdName))
		},
	}).Parse(tmpl))
	return tmpl.Execute(g.w, g.t)
}

var tmpl = `
	package {{.AttributeOrDefault "go" "package" "def" }}

	import (
		"context"
		"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
		{{range (.Attributes "go" "import") -}}
			"{{.}}"
		{{end -}}
	)

	{{range .Commands}}
		type {{handlerName .}} interface {
			{{commandName .}}(context.Context, *mongo.CommandRequest, *{{requestName .}}) (*{{responseName .}}, error)
		}

		type {{handlerName .}}Func func(context.Context, *mongo.CommandRequest, *{{requestName .}}) (*{{responseName .}}, error)

		func (f {{handlerName .}}Func) {{commandName .}}(ctx context.Context, orig *mongo.CommandRequest, req *{{requestName .}}) (*{{responseName .}}, error) {
			return f(ctx, orig, req)
		}

		type {{requestName .}} struct {
			{{with .Request}}
				{{range .Fields -}}
					{{fieldName .}} {{.TypeRef}} ` + "`" + `json:"{{.Name}}" bson:"{{.Name}}"` + "`" + `
				{{end}}
			{{end}}
		}

		type {{responseName .}} struct {
			{{with .Response}}
				{{range .Fields -}}
					{{fieldName .}} {{.TypeRef}} ` + "`" + `json:"{{.Name}}" bson:"{{.Name}}"` + "`" + `
				{{end}}
			{{end}}
		}

		type {{handlerImplName .}} struct {
			impl {{handlerName .}}
		}

		func (h *{{handlerImplName .}}) HandleCommand(ctx context.Context, resp mongo.CommandResponseWriter, req *mongo.CommandRequest) error {
			var typedReq {{requestName .}}
			if err := bson.Unmarshal(req.Document, &typedReq); err != nil {
				return &mongo.Error{
					Code: mongo.CodeFailedToParse,
					CodeName: mongo.CodeToName(mongo.CodeFailedToParse),
					Message: "invalid {{.Name}} command",
					Cause: err,
				}
			}

			typedResp, err := h.impl.{{commandName .}}(ctx, req, &typedReq)
			if err != nil {
				return err
			}

			respDoc, err := bson.Marshal(typedResp)
			if err != nil {
				return &mongo.Error{
					Code: mongo.CodeInternalError,
					CodeName: mongo.CodeToName(mongo.CodeInternalError),
					Message: "failed marshaling {{.Name}} command output",
					Cause: err,
				}
			}

			return resp.WriteSingleDocument(respDoc)
		}

		func Register{{handlerName .}}(mux *mongo.CommandMux, h {{handlerName .}}) {
			mux.Handlers["{{.Name}}"] = &{{handlerImplName .}}{impl: h}
			{{- $cmd := . -}}
			{{range (.Attributes "" "alias")}}
				mux.Handlers["{{.}}"] = mux.Handlers["{{$cmd.Name}}"]
			{{- end}}
		}
	{{end}}

	{{range .Services}}
		type {{serviceName .}} interface {
			{{range .CommandRefs -}}
				{{serviceCommandHandlerName .}}
			{{end}}
		}

		func Register{{serviceName .}}(mux *mongo.CommandMux, svc {{serviceName .}}) {
			{{range .CommandRefs -}}
				Register{{serviceCommandHandlerName .}}(mux, svc)
			{{end -}}
		}
	{{end}}

	{{range .Structs}}
		type {{.Name}} struct {
			{{range .Fields -}}
				{{fieldName .}} {{.TypeRef}} ` + "`" + `json:"{{.Name}}" bson:"{{.Name}}"` + "`" + `
			{{end}}
		}
	{{end}}
`