package user

import "context"

type Repository interface {
	FindOne(ctx context.Context, userName string) (User, error)
	FindOneId(ctx context.Context, id string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
}
