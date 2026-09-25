package converter

//go:generate go run github.com/jmattheis/goverter/cmd/goverter@v1.11.0 gen .

import (
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/adapter/repository/postgres/sqlcgen"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain/entity"
)

// goverter:converter
// goverter:output:file advertiser_generated.go
// goverter:extend UUIDToString
// goverter:extend StringToUUID
// goverter:extend CopyTime
// goverter:extend CopyTimePtr
type AdvertiserConverter interface {
	ToEntityAdvertiser(source sqlcgen.Advertiser) entity.Advertiser
	ToCreateAdvertiserParams(source entity.Advertiser) (sqlcgen.CreateAdvertiserParams, error)
	ToSaveAdvertiserParams(source entity.Advertiser) (sqlcgen.SaveAdvertiserParams, error)
}
