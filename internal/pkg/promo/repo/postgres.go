package repo

import (
	"context"
	"errors"
	"fmt"
	"ozon_replic/internal/models/models"
	"ozon_replic/internal/pkg/promo"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	uuid "github.com/satori/go.uuid"
)

const (
	readPromocode = "SELECT * FROM promocode WHERE name=$1;"
	usePromocode  = "UPDATE promocode SET leftover=leftover-1 WHERE name=$1 AND leftover>0 RETURNING *;"

	checkUniqPromoUser = `
	SELECT profile_id, promocode_id FROM order_info WHERE profile_id=$1 AND promocode_id=$2;`
)

type PromoRepo struct {
	db pgxtype.Querier // TODO: add logger
}

func NewPromoRepo(db pgxtype.Querier) *PromoRepo {
	return &PromoRepo{
		db: db,
	}
}

func (r *PromoRepo) ReadPromocode(ctx context.Context, promocodeName string) (*models.Promocode, error) {
	p := &models.Promocode{}
	if err := r.db.QueryRow(ctx, readPromocode, promocodeName).
		Scan(&p.Id, &p.Discount, &p.Name, &p.Leftover, &p.Deadline); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.Promocode{}, promo.ErrPromocodeNotFound
		}
		return &models.Promocode{}, fmt.Errorf("error happened in row.Scan: %w", err)
	}
	return p, nil
}

func (r *PromoRepo) CheckUniq(ctx context.Context, userID uuid.UUID, promocodeID int) error {
	res, err := r.db.Exec(ctx, checkUniqPromoUser, userID, promocodeID)
	if err != nil {
		return fmt.Errorf("error happened in CheckUniq sql exec: %w", err)
	}

	if res.RowsAffected() > 0 {
		return promo.ErrAlreadyUsed
	}

	return nil
}

func (r *PromoRepo) UsePromocode(ctx context.Context, promocodeName string) (*models.Promocode, error) {
	p := &models.Promocode{}
	if err := r.db.QueryRow(ctx, usePromocode, promocodeName).
		Scan(&p.Id, &p.Discount, &p.Name, &p.Leftover, &p.Deadline); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &models.Promocode{}, promo.ErrPromocodeNotFound
		}
		return &models.Promocode{}, fmt.Errorf("error happened in row.Scan: %w", err)
	}
	/*res, err := r.db.Exec(ctx, usePromocode, promocodeName)
	if err != nil {
		return &models.Promocode{}, fmt.Errorf("error happened in usePromocode sql exec: %w", err)
	}

	if res.RowsAffected() != 1 {
		return &models.Promocode{}, errors.New("failed update")
	}*/

	return p, nil
}
