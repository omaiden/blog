package comment

import (
	"github.com/acoshift/arpc/v2"
	"github.com/moonrhythm/httpmux"

	"blog/internal/auth"
)

func Mount(mux *httpmux.Mux, am *arpc.Manager) {
	authed := mux.Middleware(auth.AuthMiddleware())
	authed.Handle("POST /api/posts/comments/create", am.Handler(CreateComment))
	authed.Handle("POST /api/posts/comments/get", am.Handler(GetComments))
}
