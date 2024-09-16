package responsible

import (
	"context"
	"fmt"
	"tender/models/responsible"
	postgresql "tender/pkg/client"

	sq "github.com/Masterminds/squirrel"
	// "tender/pkg/logging"
)

type repository struct {
	client postgresql.Client
	psql   sq.StatementBuilderType
}

// FindOne implements tender.Repository.
func (r *repository) FindOne(ctx context.Context, userId string, orgId string) (responsible.Responsible, error) {
	qb := sq.Select("id, organization_id, user_id").Distinct().From("organization_responsible").Where(sq.Eq{"user_id": userId, "organization_id": orgId})
	qb = qb.PlaceholderFormat(sq.Dollar)
	sql, args, _ := qb.ToSql()
	fmt.Println(sql)
	fmt.Println(args)
	row := r.client.QueryRow(ctx, sql, args...)
	var resp responsible.Responsible
	err := row.Scan(&resp.Id, &resp.OrganizationId, &resp.UserId)
	if err != nil {
		return responsible.Responsible{}, err
	}
	return resp, nil
}

func (r *repository) FindAll(ctx context.Context, userId string) ([]responsible.Responsible, error) {
	qb := r.psql.Select("id, organization_id, user_id").From("organization_responsible").Where(sq.Eq{"user_id": userId})
	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]responsible.Responsible, 0)
	for rows.Next() {
		var resp responsible.Responsible
		err = rows.Scan(&resp.Id, &resp.OrganizationId, &resp.UserId)
		if err != nil {
			return nil, err
		}
		res = append(res, resp)
	}
	return res, nil
}

func (r *repository) FindAlls(ctx context.Context) ([]responsible.Responsible, error) {
	qb := r.psql.Select("id, organization_id, user_id").From("organization_responsible")
	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]responsible.Responsible, 0)
	for rows.Next() {
		var resp responsible.Responsible
		err = rows.Scan(&resp.Id, &resp.OrganizationId, &resp.UserId)
		if err != nil {
			return nil, err
		}
		res = append(res, resp)
	}
	return res, nil
}

// TODO FINDALL FOR TENDERS MY  RETURN ALLS ORGS FOR USERID
func NewRepository(clinet postgresql.Client) responsible.Repository {
	return &repository{
		client: clinet,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
