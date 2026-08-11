package ui

import (
	"RATAC/views"
	"net/http"

	"github.com/a-h/templ"
)

func render(w http.ResponseWriter, r *http.Request, content templ.Component, list_pacientes_page bool) {
	if r.Header.Get("HX-Request") == "true" {
		content.Render(r.Context(), w)
		return
	}
	views.Page(content, list_pacientes_page).Render(r.Context(), w)
}
