package log

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yumenaka/comigo/pagoda/pkg/tests"
)

func TestCtxSet(t *testing.T) {
	ctx, _ := tests.NewContext(echo.New(), "/")
	logger := Ctx(ctx)
	assert.NotNil(t, logger)

	logger = logger.With("a", "b")
	Set(ctx, logger)

	got := Ctx(ctx)
	assert.Equal(t, got, logger)
}
