package tender

import "time"

type CreateTenderDTO struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ServiceType     string `json:"serviceType"`
	OrganizationId  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type UpdateTenderDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}

func (t UpdateTenderDTO) ToTender() Tender {
	return Tender{
		Name:        t.Name,
		Description: t.Description,
		ServiceType: t.ServiceType,
	}
}

func (t CreateTenderDTO) ToTender() Tender {
	return Tender{
		Name:           t.Name,
		Description:    t.Description,
		ServiceType:    t.ServiceType,
		OrganizationId: t.OrganizationId,
	}
}

type TenderDTO struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
	Status      string `json:"status"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}

func ToTenderDTO(t Tender) TenderDTO {
	return TenderDTO{
		Id:          t.Id,
		Name:        t.Name,
		Description: t.Description,
		ServiceType: t.ServiceType,
		Status:      t.Status,
		Version:     t.Version,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
	}
}

func ToTenderDTOs(tenders []Tender) []TenderDTO {
	dtos := make([]TenderDTO, len(tenders))
	for i, tender := range tenders {
		dtos[i] = ToTenderDTO(tender)
	}
	return dtos
}
