package tu

import (
	"testing"

	"github.com/acoshift/pgsql/pgctx"
	"github.com/moonrhythm/randid"
	"github.com/stretchr/testify/require"
)

func (ctx *Context) CreateUser(t *testing.T, username string) string {
	t.Helper()

	userID := randid.MustGenerate()
	_, err := pgctx.Exec(ctx.Ctx(), `
		insert into users (id, username)
		values ($1, $2)
	`, userID, username)
	require.NoError(t, err)
	return userID.String()
}
