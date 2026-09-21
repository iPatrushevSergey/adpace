package dto

type CreateAdvertiserInput struct {
	Name    string
	Country string
}

type GetAdvertiserInput struct {
	AdvertiserID string
}

type PatchAdvertiserInput struct {
	AdvertiserID string
	Name         *string
	Country      *string
}

type PutAdvertiserInput struct {
	AdvertiserID string
	Name         string
	Country      string
}

type DeleteAdvertiserInput struct {
	AdvertiserID string
}
