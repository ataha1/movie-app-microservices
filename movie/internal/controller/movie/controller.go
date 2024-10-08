package movie

import (
	"context"
	"errors"

	metadatamodel "github.com/ataha1/movie-app/metadata/pkg/model"
	"github.com/ataha1/movie-app/movie/internal/gateway"
	"github.com/ataha1/movie-app/movie/pkg/model"
	ratingmodel "github.com/ataha1/movie-app/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type ratingGateWay interface{
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType)(float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface{
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

type Controller struct{
	ratingGateWay ratingGateWay
	metadataGateway metadataGateway
}

func New(metadataGateway metadataGateway, ratingGateWay ratingGateWay) *Controller{
	return &Controller{
		metadataGateway: metadataGateway,
		ratingGateWay: ratingGateWay,
	}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error){
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound){
		return nil, ErrNotFound
	} else if err != nil{
		return  nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateWay.GetAggregatedRating(ctx, ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)
	if err != nil && errors.Is(err, gateway.ErrNotFound){
		return nil, ErrNotFound
	} else if err != nil{
		return nil, err
	} else {
		details.Rating = rating
	}
	return details, nil
}

