package responsible

import "context"

type Repository interface {
	FindOne(ctx context.Context, userId string, orgId string) (Responsible, error)
	FindAll(ctx context.Context, userId string) ([]Responsible, error)
	FindAlls(ctx context.Context) ([]Responsible, error)
}
