package memory

import (
	"context"
	"sync"

	"github.com/ataha1/movie-app/rating/internal/repository"
	"github.com/ataha1/movie-app/rating/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository{
	return &Repository{
		data: map[model.RecordType]map[model.RecordID][]model.Rating{},
	}
}

func(r *Repository)Get(_ context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error){
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.data[recordType]; !ok{
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0{
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordId], nil
}

func(r *Repository)Put(_ context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating)error{
	r.Lock()
	defer r.Unlock()
	if _, ok := r.data[recordType]; !ok{
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}
	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)
	return nil
}