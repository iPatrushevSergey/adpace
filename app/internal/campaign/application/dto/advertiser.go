package dto

import "github.com/google/uuid"

type CreateAdvertiserInput struct {
	Name    string
	Country string
}

type GetAdvertiserInput struct {
	AdvertiserID uuid.UUID
}

type PatchAdvertiserInput struct {
	AdvertiserID uuid.UUID
	Name         *string
	Country      *string
}

type PutAdvertiserInput struct {
	AdvertiserID uuid.UUID
	Name         string
	Country      string
}

type DeleteAdvertiserInput struct {
	AdvertiserID uuid.UUID
}
