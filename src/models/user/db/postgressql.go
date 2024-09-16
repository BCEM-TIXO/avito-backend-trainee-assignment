package user

import (
	"context"
	"tender/models/user"
	postgresql "tender/pkg/client"
	// "tender/pkg/logging"
)

type repository struct {
	client postgresql.Client
}

// FindOne implements tender.Repository.
func (r *repository) FindOne(ctx context.Context, userName string) (user.User, error) {
	q := `SELECT 
			id, username
		  FROM employee
		  WHERE username = $1
		  `
	row := r.client.QueryRow(ctx, q, userName)
	var u user.User
	err := row.Scan(&u.Id, &u.UserName)
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}

func (r *repository) FindAll(ctx context.Context) ([]user.User, error) {
	q := `SELECT 
			id, username
		  FROM employee
		  `
	rows, err := r.client.Query(ctx, q)
	defer rows.Close()
	if err = rows.Err(); err != nil {
		return nil, err
	}
	users := make([]user.User, 0)
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.Id, &u.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *repository) FindOneId(ctx context.Context, id string) (user.User, error) {
	q := `SELECT 
			id, username
		  FROM employee
		  WHERE id = $1
		  `
	row := r.client.QueryRow(ctx, q, id)
	var u user.User
	err := row.Scan(&u.Id, &u.UserName)
	if err != nil {
		return user.User{}, err
	}
	return u, nil
}

func NewRepository(clinet postgresql.Client) user.Repository {
	return &repository{
		client: clinet,
	}
}
