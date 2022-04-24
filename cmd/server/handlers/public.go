package handlers

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Public(fs embed.FS) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		publicFS := http.FS(fs)
		fileServer := http.FileServer(publicFS)
		r.URL.Path = fmt.Sprintf("/public%s", p.ByName("filepath"))
		fileServer.ServeHTTP(w, r)
	}
}
