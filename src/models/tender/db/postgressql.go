package tender

import (
	"context"
	// "strings"
	// "errors"
	"tender/models/tender"
	postgresql "tender/pkg/client"

	sq "github.com/Masterminds/squirrel"
)

type repository struct {
	client postgresql.Client
	psql   sq.StatementBuilderType
}

// Create implements tender.Repository.
func (r *repository) Create(ctx context.Context, t *tender.Tender) error {
	q := `INSERT INTO tender 
			(name, description, type, organization_id) 
		  VALUES
			($1, $2, $3, $4)
		  RETURNING id, status, created_at, version`
	err := r.client.QueryRow(
		ctx, q, t.Name, t.Description,
		t.ServiceType, t.OrganizationId).Scan(&t.Id, &t.Status, &t.CreatedAt, &t.Version)
	return err
}

// FindAll implements tender.Repository.
func (r *repository) FindAll(ctx context.Context, qm *tender.FindAllQueryModifier) ([]tender.Tender, error) {
	qb := r.psql.Select("id, name, description, status, type, version, created_at, organization_id").From("tender").OrderBy("name")
	if qm != nil {
		qb = qb.Offset(uint64(qm.Pagination.Offset)).Limit(uint64(qm.Pagination.Limit))
		if qm.Organization_id != nil {
			qb = qb.Where("organization_id = ANY($1)", qm.Organization_id)
		}
		if qm.ServiceTypes != nil {
			qb = qb.Where("service_type = ANY($1)", qm.ServiceTypes)
		}
	}
	q, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.client.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tenders := make([]tender.Tender, 0)
	for rows.Next() {
		var t tender.Tender
		err = rows.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.Version, &t.CreatedAt, &t.OrganizationId)
		if err != nil {
			return nil, err
		}
		tenders = append(tenders, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tenders, nil
}

// FindOne implements tender.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (tender.Tender, error) {
	q := `SELECT 
			id, name, description, status, type, version, created_at, organization_id
		  FROM tender
		  WHERE id = $1
		  `
	row := r.client.QueryRow(ctx, q, id)
	var t tender.Tender
	err := row.Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.Version, &t.CreatedAt, &t.OrganizationId)
	if err != nil {
		return tender.Tender{}, err
	}
	return t, nil
}

func (r *repository) Update(ctx context.Context, t *tender.Tender, fields map[string]interface{}) error {
	qb := r.psql.Update("tender").Where(sq.Eq{"id": t.Id}).SetMap(fields).Suffix("RETURNING id, name, description, status, type, version, created_at, organization_id")
	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	err = r.client.QueryRow(ctx, sql, args...).Scan(&t.Id, &t.Name, &t.Description, &t.Status, &t.ServiceType, &t.Version, &t.CreatedAt, &t.OrganizationId)
	return err
}

func NewRepository(clinet postgresql.Client) tender.Repository {
	return &repository{
		client: clinet,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
