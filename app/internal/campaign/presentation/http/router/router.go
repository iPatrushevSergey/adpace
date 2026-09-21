package router

import (
	"github.com/gin-gonic/gin"
	appport "github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/usecase"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/presentation/http/handler"
)

func New(uc usecase.AdvertiserUseCases, log appport.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api/v1")
	protected := r.Group("/api/v1")

	RegisterAdvertiserRoutes(public, protected, uc, log)

	return r
}

// RegisterAdvertiserRoutes registers advertiser module HTTP routes.
func RegisterAdvertiserRoutes(
	public, protected *gin.RouterGroup,
	uc usecase.AdvertiserUseCases,
	log appport.Logger,
) {
	h := handler.NewAdvertiserHandler(uc, log)

	advertisers := protected.Group("/advertisers")
	advertisers.POST("", h.Create)
	advertisers.GET("/:advertiser_id", h.Get)
	advertisers.PATCH("/:advertiser_id", h.Patch)
	advertisers.PUT("/:advertiser_id", h.Put)
	advertisers.DELETE("/:advertiser_id", h.Delete)
}
