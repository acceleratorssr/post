package web

import (
	"context"
	"github.com/gin-gonic/gin"
)

type handler interface {
	RegisterRoutes(ctx context.Context, engine *gin.Engine)
}
