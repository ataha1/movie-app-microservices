package main

import (
	"log"
	"net/http"

	"github.com/ataha1/movie-app/movie/internal/controller/movie"
	metadataGateway "github.com/ataha1/movie-app/movie/internal/gateway/metadata/http"
	ratingGateWay "github.com/ataha1/movie-app/movie/internal/gateway/rating/http"
	httpHandler "github.com/ataha1/movie-app/movie/internal/handler/http"
)

func main() {
	log.Print("Starting the movie service")
	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateWay := ratingGateWay.New("localhost:8082")
	ctrl := movie.New(metadataGateway, ratingGateWay)
	h := httpHandler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(":8083", nil); err != nil{
		panic(err)
	}
}