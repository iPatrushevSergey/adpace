package dto

type CreateCampaignInput struct {
	AdvertiserID string
	Name         string
	BudgetTotal  int64
	BudgetDaily  int64
}

type GetCampaignInput struct {
	CampaignID string
}

type PatchCampaignInput struct {
	CampaignID  string
	Name        *string
	BudgetTotal *int64
	BudgetDaily *int64
}

type PutCampaignInput struct {
	CampaignID  string
	Name        string
	BudgetTotal int64
	BudgetDaily int64
}

type DeleteCampaignInput struct {
	CampaignID string
}

type PauseCampaignInput struct {
	CampaignID string
	Reason     string
}

type ResumeCampaignInput struct {
	CampaignID string
}
