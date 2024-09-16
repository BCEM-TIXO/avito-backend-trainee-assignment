package bid

import (
	"context"
	"tender/models/tender"
)

type Repository interface {
	Create(ctx context.Context, b *Bid) error
	FindAll(ctx context.Context, tenderId string, qm *tender.FindAllQueryModifier) ([]Bid, error)
	FindOne(ctx context.Context, id string) (Bid, error)
	Update(ctx context.Context, t *Bid, fields map[string]interface{}) error
	// Create (ctx context.Context, tender Tender) (string, error)
	// Create (ctx context.Context, tender Tender) (string, error)
}
