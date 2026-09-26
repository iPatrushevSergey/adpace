package converter

//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v1.11.0 gen .

import (
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

// goverter:converter
// goverter:output:file campaign_generated.go
// goverter:extend CopyTime
type CampaignConverter interface {
	ToEntityCampaign(source sqlcgen.Campaign) entity.Campaign
	ToCreateCampaignParams(source entity.Campaign) sqlcgen.CreateCampaignParams
	ToSaveCampaignParams(source entity.Campaign) sqlcgen.SaveCampaignParams
}
