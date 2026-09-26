package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iPatrushevSergey/adpace/app/internal/campaign/domain"
)

type AdvertiserOption func(*Advertiser)

type Advertiser struct {
	AdvertiserID uuid.UUID
	Name         string
	Country      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewAdvertiser(
	advertiserID uuid.UUID,
	name, country string,
	updatedAt time.Time,
	opts ...AdvertiserOption,
) (Advertiser, error) {
	a := Advertiser{
		AdvertiserID: advertiserID,
		UpdatedAt:    updatedAt,
	}

	if err := a.SetName(&name); err != nil {
		return Advertiser{}, err
	}
	if err := a.SetCountry(&country); err != nil {
		return Advertiser{}, err
	}

	for _, opt := range opts {
		opt(&a)
	}
	return a, nil
}

func WithAdvertiserCreatedAt(t time.Time) AdvertiserOption {
	return func(a *Advertiser) { a.CreatedAt = t }
}

func (a *Advertiser) SetName(name *string) error {
	if name == nil {
		return nil
	}
	if !IsValidAdvertiserName(*name) {
		return fmt.Errorf("%w: name cannot be empty", domain.ErrBadInput)
	}
	a.Name = *name
	return nil
}

func (a *Advertiser) SetCountry(country *string) error {
	if country == nil {
		return nil
	}
	if !IsValidCountry(*country) {
		return fmt.Errorf("%w: invalid country %q", domain.ErrBadInput, *country)
	}
	a.Country = *country
	return nil
}

func IsValidAdvertiserName(name string) bool {
	return name != ""
}

func IsValidCountry(c string) bool {
	_, ok := validCountries[c]
	return ok
}

var validCountries = map[string]struct{}{
	"AD": {}, "AE": {}, "AF": {}, "AG": {}, "AI": {}, "AL": {}, "AM": {}, "AO": {}, "AQ": {},
	"AR": {}, "AS": {}, "AT": {}, "AU": {}, "AW": {}, "AX": {}, "AZ": {}, "BA": {}, "BB": {},
	"BD": {}, "BE": {}, "BF": {}, "BG": {}, "BH": {}, "BI": {}, "BJ": {}, "BL": {}, "BM": {},
	"BN": {}, "BO": {}, "BQ": {}, "BR": {}, "BS": {}, "BT": {}, "BV": {}, "BW": {}, "BY": {},
	"BZ": {}, "CA": {}, "CC": {}, "CD": {}, "CF": {}, "CG": {}, "CH": {}, "CI": {}, "CK": {},
	"CL": {}, "CM": {}, "CN": {}, "CO": {}, "CR": {}, "CU": {}, "CV": {}, "CW": {}, "CX": {},
	"CY": {}, "CZ": {}, "DE": {}, "DJ": {}, "DK": {}, "DM": {}, "DO": {}, "DZ": {}, "EC": {},
	"EE": {}, "EG": {}, "EH": {}, "ER": {}, "ES": {}, "ET": {}, "FI": {}, "FJ": {}, "FK": {},
	"FM": {}, "FO": {}, "FR": {}, "GA": {}, "GB": {}, "GD": {}, "GE": {}, "GF": {}, "GG": {},
	"GH": {}, "GI": {}, "GL": {}, "GM": {}, "GN": {}, "GP": {}, "GQ": {}, "GR": {}, "GS": {},
	"GT": {}, "GU": {}, "GW": {}, "GY": {}, "HK": {}, "HM": {}, "HN": {}, "HR": {}, "HT": {},
	"HU": {}, "ID": {}, "IE": {}, "IL": {}, "IM": {}, "IN": {}, "IO": {}, "IQ": {}, "IR": {},
	"IS": {}, "IT": {}, "JE": {}, "JM": {}, "JO": {}, "JP": {}, "KE": {}, "KG": {}, "KH": {},
	"KI": {}, "KM": {}, "KN": {}, "KP": {}, "KR": {}, "KW": {}, "KY": {}, "KZ": {}, "LA": {},
	"LB": {}, "LC": {}, "LI": {}, "LK": {}, "LR": {}, "LS": {}, "LT": {}, "LU": {}, "LV": {},
	"LY": {}, "MA": {}, "MC": {}, "MD": {}, "ME": {}, "MF": {}, "MG": {}, "MH": {}, "MK": {},
	"ML": {}, "MM": {}, "MN": {}, "MO": {}, "MP": {}, "MQ": {}, "MR": {}, "MS": {}, "MT": {},
	"MU": {}, "MV": {}, "MW": {}, "MX": {}, "MY": {}, "MZ": {}, "NA": {}, "NC": {}, "NE": {},
	"NF": {}, "NG": {}, "NI": {}, "NL": {}, "NO": {}, "NP": {}, "NR": {}, "NU": {}, "NZ": {},
	"OM": {}, "PA": {}, "PE": {}, "PF": {}, "PG": {}, "PH": {}, "PK": {}, "PL": {}, "PM": {},
	"PN": {}, "PR": {}, "PS": {}, "PT": {}, "PW": {}, "PY": {}, "QA": {}, "RE": {}, "RO": {},
	"RS": {}, "RU": {}, "RW": {}, "SA": {}, "SB": {}, "SC": {}, "SD": {}, "SE": {}, "SG": {},
	"SH": {}, "SI": {}, "SJ": {}, "SK": {}, "SL": {}, "SM": {}, "SN": {}, "SO": {}, "SR": {},
	"SS": {}, "ST": {}, "SV": {}, "SX": {}, "SY": {}, "SZ": {}, "TC": {}, "TD": {}, "TF": {},
	"TG": {}, "TH": {}, "TJ": {}, "TK": {}, "TL": {}, "TM": {}, "TN": {}, "TO": {}, "TR": {},
	"TT": {}, "TV": {}, "TW": {}, "TZ": {}, "UA": {}, "UG": {}, "UM": {}, "US": {}, "UY": {},
	"UZ": {}, "VA": {}, "VC": {}, "VE": {}, "VG": {}, "VI": {}, "VN": {}, "VU": {}, "WF": {},
	"WS": {}, "YE": {}, "YT": {}, "ZA": {}, "ZM": {}, "ZW": {},
}
