package repository

import (
	"context"
	"sync"
	"time"

	"github.com/Ycnik/suprise/internal/model"
)

type MemorySoldatRepository struct {
	mu       sync.Mutex
	nextID   int
	soldaten []model.Soldat
}

func NewMemorySoldatRepository() *MemorySoldatRepository {
	return &MemorySoldatRepository{nextID: 1000}
}

func (r *MemorySoldatRepository) List(ctx context.Context) ([]model.Soldat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	soldaten := make([]model.Soldat, len(r.soldaten))
	copy(soldaten, r.soldaten)
	return soldaten, nil
}

func (r *MemorySoldatRepository) FindByID(ctx context.Context, id int) (*model.Soldat, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, soldat := range r.soldaten {
		if soldat.ID == id {
			found := soldat
			return &found, nil
		}
	}
	return nil, ErrSoldatNotFound
}

func (r *MemorySoldatRepository) Create(ctx context.Context, soldat *model.Soldat) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	soldat.ID = r.nextID
	soldat.Version = 0
	soldat.Erzeugt = now
	soldat.Aktualisiert = now

	r.nextID++
	r.soldaten = append(r.soldaten, *soldat)
	return nil
}
