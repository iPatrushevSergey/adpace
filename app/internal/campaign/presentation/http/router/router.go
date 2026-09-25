package router

import (
	"github.com/gin-gonic/gin"
	appport "github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/usecase"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/presentation/http/handler"
)

func New(
	advUC usecase.AdvertiserUseCases,
	camUC usecase.CampaignUseCases,
	log appport.Logger,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	public := r.Group("/api/v1")
	protected := r.Group("/api/v1")

	RegisterAdvertiserRoutes(public, protected, advUC, log)
	RegisterCampaignRoutes(public, protected, camUC, log)

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

func RegisterCampaignRoutes(
	public, protected *gin.RouterGroup,
	uc usecase.CampaignUseCases,
	log appport.Logger,
) {
	h := handler.NewCampaignHandler(uc, log)

	campaigns := protected.Group("/campaigns")
	campaigns.POST("", h.Create)
	campaigns.GET("/:campaign_id", h.Get)
	campaigns.PATCH("/:campaign_id", h.Patch)
	campaigns.PUT("/:campaign_id", h.Put)
	campaigns.DELETE("/:campaign_id", h.Delete)
	campaigns.POST("/:campaign_id/pause", h.Pause)
	campaigns.POST("/:campaign_id/resume", h.Resume)
}
