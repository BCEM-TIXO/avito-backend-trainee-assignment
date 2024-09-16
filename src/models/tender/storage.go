package tender

import (
	"context"
	repeatable "tender/pkg/utils"
)

type FindAllQueryModifier struct {
	Pagination      repeatable.Pagination
	ServiceTypes    []string
	Organization_id []string
}

type Repository interface {
	Create(ctx context.Context, t *Tender) error
	FindAll(ctx context.Context, qm *FindAllQueryModifier) ([]Tender, error)
	FindOne(ctx context.Context, id string) (Tender, error)
	Update(ctx context.Context, t *Tender, fields map[string]interface{}) error
	// Create (ctx context.Context, tender Tender) (string, error)
	// Create (ctx context.Context, tender Tender) (string, error)
}
