package metadata

import (
	"context"
	"errors"

	"github.com/ataha1/movie-app/metadata/internal/repository"
	"github.com/ataha1/movie-app/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface{
	Get(ctx context.Context, id string)(*model.Metadata, error)
}

type Controller struct{
	repo metadataRepository
}

func New(repo metadataRepository)*Controller{
	return &Controller{
		repo: repo,
	}
}

func(c *Controller)Get(ctx context.Context, id string)(*model.Metadata, error){
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound){
		return nil, err
	}
	return res, nil
}