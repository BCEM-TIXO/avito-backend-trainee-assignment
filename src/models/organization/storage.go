package organization

import "context"

type Repository interface {
	FindOne(ctx context.Context, id string) (Organization, error)
	FindAll(ctx context.Context) ([]Organization, error)
}
