package organization

import (
	"context"
	"tender/models/organization"
	postgresql "tender/pkg/client"
	// "tender/pkg/logging"
)

type repository struct {
	client postgresql.Client
}

// FindOne implements tender.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (organization.Organization, error) {
	q := `SELECT 
			id, name
		  FROM organization
		  WHERE id = $1
		  `
	row := r.client.QueryRow(ctx, q, id)
	var u organization.Organization
	err := row.Scan(&u.Id, &u.Name)
	if err != nil {
		return organization.Organization{}, err
	}
	return u, nil
}

func (r *repository) FindAll(ctx context.Context) ([]organization.Organization, error) {
	q := `SELECT 
			id, name
		  FROM organization
		  `
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]organization.Organization, 0)
	for rows.Next() {
		var resp organization.Organization
		err = rows.Scan(&resp.Id, &resp.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, resp)
	}
	return res, nil
}

func NewRepository(clinet postgresql.Client) organization.Repository {
	return &repository{
		client: clinet,
	}
}
