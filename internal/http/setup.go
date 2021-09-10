package http

import (
	_ "embed"
	"net/http"
	"text/template"
)

//go:embed setup.sh
var setupScript string

var setupScriptTemplate = template.Must(template.New("setupScript").Parse(setupScript))

func setupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	templateData := struct{ HostName string }{r.Host}
	_ = setupScriptTemplate.Execute(w, templateData)
	return
}
