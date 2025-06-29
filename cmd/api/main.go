package main

import (
	"blog/comment"
	"blog/post"
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/acoshift/arpc/v2"
	"github.com/acoshift/configfile"
	"github.com/acoshift/pgsql/pgctx"
	"github.com/moonrhythm/httpmux"
	"github.com/moonrhythm/parapet"
	"github.com/moonrhythm/parapet/pkg/cors"

	"blog/pkg/api"
	"blog/pkg/logs"
	"blog/pkg/ops"
	"blog/user"
)

func main() {
	cfg := configfile.NewEnvReader()

	port := cfg.StringDefault("port", "8081")

	db, err := sql.Open("postgres", cfg.MustString("db_url"))
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}
	defer db.Close()

	db.SetMaxIdleConns(cfg.IntDefault("db_max_idle_conns", 30))
	db.SetMaxOpenConns(cfg.IntDefault("db_max_open_conns", 50))
	db.SetConnMaxIdleTime(cfg.DurationDefault("db_conn_max_idle_time", 30*time.Second))

	ctx := context.Background()
	ctx = pgctx.NewContext(ctx, db)

	ops.Init(ctx)
	defer ops.Close()

	am := arpc.New()
	am.WrapError = api.WrapError
	am.OnOK(logs.ReportRPCOK)
	am.OnError(logs.ReportRPCError)

	m := httpmux.New()
	m.Handle("/", am.NotFoundHandler())
	user.Mount(m, am)
	post.Mount(m, am)
	comment.Mount(m, am)

	s := parapet.NewBackend()
	s.Addr = ":" + port
	s.Handler = m

	s.Use(logs.InjectRecord())
	s.Use(cors.CORS{
		AllowAllOrigins: true,
		AllowMethods:    []string{"POST"},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		MaxAge:          time.Hour,
	})
	s.Use(ops.Trace())
	s.Use(ops.Recovery())
	s.Use(ops.InjectRequestIDToSpan())
	s.Use(parapet.MiddlewareFunc(pgctx.Middleware(db)))

	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("cannot start api server; %v", err)
	}
}
