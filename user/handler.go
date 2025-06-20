package user

import (
	"github.com/acoshift/arpc/v2"
	"github.com/moonrhythm/httpmux"
)

func Mount(mux *httpmux.Mux, am *arpc.Manager) {
	mux.Handle("POST /api/users/register", am.Handler(Register))
	mux.Handle("POST /api/users/login", am.Handler(Login))
}
