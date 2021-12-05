package documents

import (
	"net/http"

	"go.uber.org/zap"
)

func NewHandler(log *zap.Logger) Handler {
	return Handler{
		Log: log,
	}
}

type Handler struct {
	Log *zap.Logger
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("received request")

	// Render the HTML.
	td := TemplateData{
		Name: "World",
	}
	doc := Document(td)

	w.Header().Add("Content-Type", "application/pdf")
	if err := ToPDF(r.Context(), doc, w); err != nil {
		errmsg := "failed to render PDF"
		h.Log.Error(errmsg, zap.Error(err))
		http.Error(w, errmsg, http.StatusInternalServerError)
	}
}
