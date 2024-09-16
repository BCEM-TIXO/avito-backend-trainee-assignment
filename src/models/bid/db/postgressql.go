package bid

import (
	"context"
	"errors"
	"tender/models/bid"
	"tender/models/tender"
	postgresql "tender/pkg/client"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	psql   sq.StatementBuilderType
}

// Create implements tender.Repository.
func (r *repository) Create(ctx context.Context, b *bid.Bid) error {
	qb := r.psql.Insert("bid").Columns("name", "description", "tender_id", "author_type", "author_id").Values(b.Name, b.Description, b.TenderId, b.AuthorType, b.AuthorId)
	qb = qb.Suffix("RETURNING \"id\", \"status\", \"created_at\", \"version\"")
	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}
	err = r.client.QueryRow(ctx, sql, args...).Scan(&b.Id, &b.Status, &b.CreatedAt, &b.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "bid_tender_id_fkey" {
				return errors.New("Тендер не найден.")
			}
		}
	}
	return err
}

// FindAll implements tender.Repository.
func (r *repository) FindAll(ctx context.Context, tenderId string, qm *tender.FindAllQueryModifier) ([]bid.Bid, error) {
	qb := r.psql.Select("id, name, description, status, tenderId, author_type, author_id, version, created_at").From("bid").Where(sq.Eq{"tender_id": tenderId}).OrderBy("name")
	if qm != nil {
		qb = qb.Offset(uint64(qm.Pagination.Offset)).Limit(uint64(qm.Pagination.Limit))
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
	bids := make([]bid.Bid, 0)
	for rows.Next() {
		var b bid.Bid
		err = rows.Scan(&b.Id, &b.Name, &b.Description, &b.Status, &b.TenderId, &b.AuthorType, &b.AuthorId, &b.Version, &b.CreatedAt)
		if err != nil {
			return nil, err
		}
		bids = append(bids, b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bids, nil
}

// FindOne implements tender.Repository.
func (r *repository) FindOne(ctx context.Context, id string) (bid.Bid, error) {
	qb := r.psql.Select("id, name, description, status, tenderId, author_type, author_id, version, created_at").From("bid").Where(sq.Eq{"id": id})
	q, args, err := qb.ToSql()
	if err != nil {
		return bid.Bid{}, err
	}
	row := r.client.QueryRow(ctx, q, args...)
	var b bid.Bid
	err = row.Scan(&b.Id, &b.Name, &b.Description, &b.Status, &b.TenderId, &b.AuthorType, &b.AuthorId, &b.Version, &b.CreatedAt)
	if err != nil {
		return bid.Bid{}, err
	}
	return b, nil
}

func (r *repository) Update(ctx context.Context, b *bid.Bid, fields map[string]interface{}) error {
	qb := r.psql.Update("bid").Where(sq.Eq{"id": b.Id}).SetMap(fields).Suffix("RETURNING id, name, description, status, type, version, created_at, organization_id")
	sql, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	err = r.client.QueryRow(ctx, sql, args...).Scan(&b.Id, &b.Name, &b.Description, &b.Status, &b.TenderId, &b.AuthorType, &b.AuthorId, &b.Version, &b.CreatedAt)
	return err
}

func NewRepository(clinet postgresql.Client) bid.Repository {
	return &repository{
		client: clinet,
		psql:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
