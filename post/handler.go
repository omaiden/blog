package post

import (
	"blog/internal/auth"

	"github.com/acoshift/arpc/v2"
	"github.com/moonrhythm/httpmux"
)

func Mount(mux *httpmux.Mux, am *arpc.Manager) {
	authed := mux.Middleware(auth.AuthMiddleware())
	authed.Handle("/api/posts", am.Handler(ListPosts))
	authed.Handle("POST /api/posts", am.Handler(CreatePost))
	authed.Handle("POST /api/posts/get", am.Handler(GetPost))
	authed.Handle("POST /api/posts/update", am.Handler(UpdatePost))
	authed.Handle("POST /api/posts/delete", am.Handler(DeletePost))
}
