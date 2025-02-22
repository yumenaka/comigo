package middleware

import (
	"fmt"
	"testing"

	"github.com/yumenaka/comigo/pagoda/ent"
	"github.com/yumenaka/comigo/pagoda/pkg/context"
	"github.com/yumenaka/comigo/pagoda/pkg/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUser(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	ctx.SetParamNames("user")
	ctx.SetParamValues(fmt.Sprintf("%d", usr.ID))
	_ = tests.ExecuteMiddleware(ctx, LoadUser(c.ORM))
	ctxUsr, ok := ctx.Get(context.UserKey).(*ent.User)
	require.True(t, ok)
	assert.Equal(t, usr.ID, ctxUsr.ID)
}
