package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Ycnik/suprise/internal/model"
	"gorm.io/gorm"
)

var ErrSoldatNotFound = errors.New("soldat not found")

type SoldatRepository interface {
	List(ctx context.Context) ([]model.Soldat, error)
	FindByID(ctx context.Context, id int) (*model.Soldat, error)
	Create(ctx context.Context, soldat *model.Soldat) error
}

type GormSoldatRepository struct {
	db *gorm.DB
}

func NewGormSoldatRepository(db *gorm.DB) *GormSoldatRepository {
	return &GormSoldatRepository{db: db}
}

func (r *GormSoldatRepository) List(ctx context.Context) ([]model.Soldat, error) {
	var soldaten []model.Soldat
	err := r.db.WithContext(ctx).
		Preload("Ausruestung").
		Preload("Verletzungen").
		Order("id asc").
		Find(&soldaten).
		Error
	return soldaten, err
}

func (r *GormSoldatRepository) FindByID(ctx context.Context, id int) (*model.Soldat, error) {
	var soldat model.Soldat
	err := r.db.WithContext(ctx).
		Preload("Ausruestung").
		Preload("Verletzungen").
		First(&soldat, "id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrSoldatNotFound
	}
	if err != nil {
		return nil, err
	}
	return &soldat, nil
}

func (r *GormSoldatRepository) Create(ctx context.Context, soldat *model.Soldat) error {
	now := time.Now().UTC()
	soldat.Version = 0
	soldat.Erzeugt = now
	soldat.Aktualisiert = now

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(soldat).Error
	})
}
