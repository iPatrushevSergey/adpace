package usecase

import (
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/dto"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/application/port"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

// AdvertiserUseCases groups advertiser use cases exposed to the presentation layer.
type AdvertiserUseCases struct {
	Create port.UseCase[dto.CreateAdvertiserInput, entity.Advertiser]
	Get    port.UseCase[dto.GetAdvertiserInput, entity.Advertiser]
	Patch  port.UseCase[dto.PatchAdvertiserInput, entity.Advertiser]
	Put    port.UseCase[dto.PutAdvertiserInput, struct{}]
	Delete port.UseCase[dto.DeleteAdvertiserInput, struct{}]
}

type CampaignUseCases struct {
	Create port.UseCase[dto.CreateCampaignInput, entity.Campaign]
	Get    port.UseCase[dto.GetCampaignInput, entity.Campaign]
	Patch  port.UseCase[dto.PatchCampaignInput, entity.Campaign]
	Put    port.UseCase[dto.PutCampaignInput, struct{}]
	Delete port.UseCase[dto.DeleteCampaignInput, struct{}]
	Pause  port.UseCase[dto.PauseCampaignInput, entity.Campaign]
	Resume port.UseCase[dto.ResumeCampaignInput, entity.Campaign]
}
