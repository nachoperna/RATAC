package ui

import (
	"RATAC/views"
	"net/http"

	"github.com/a-h/templ"
)

type PageHandler struct {
	authHandler *AuthHandler
}

func NewPageHandler(AuthHandler *AuthHandler) *PageHandler {
	return &PageHandler{
		authHandler: AuthHandler,
	}
}

func (h *PageHandler) Render(w http.ResponseWriter, r *http.Request, content templ.Component, list_pacientes_page bool) {
	if r.Header.Get("HX-Request") == "true" {
		content.Render(r.Context(), w)
		return
	}
	sesion_activa := h.authHandler.SesionActiva(w,r)
	views.Page(content, list_pacientes_page, sesion_activa, "").Render(r.Context(), w)
}
