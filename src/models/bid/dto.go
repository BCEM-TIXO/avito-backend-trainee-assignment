package bid

import "time"

type BidDTO struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType"`
	AuthorId   string `json:"authorId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}

type CreateBidDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenderId    string `json:"tenderId"`
	AuthorType  string `json:"authorType"`
	AuthorId    string `json:"authorId"`
}

func (b CreateBidDTO) ToBid() Bid {
	return Bid{
		Name:        b.Name,
		Description: b.Description,
		TenderId:    b.TenderId,
		AuthorType:  b.AuthorType,
		AuthorId:    b.AuthorId,
	}
}

func (b Bid) ToBidDTO() BidDTO {
	return BidDTO{
		Id:         b.Id,
		Name:       b.Name,
		Status:     b.Status,
		Version:    b.Version,
		AuthorType: b.AuthorType,
		AuthorId:   b.AuthorId,
		CreatedAt:  b.CreatedAt.Format(time.RFC3339),
	}
}

func ToBidDTOs(tenders []Bid) []BidDTO {
	dtos := make([]BidDTO, len(tenders))
	for i, tender := range tenders {
		dtos[i] = tender.ToBidDTO()
	}
	return dtos
}
