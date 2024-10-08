package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ataha1/movie-app/movie/internal/controller/movie"
	metadataGateway "github.com/ataha1/movie-app/movie/internal/gateway/metadata/http"
	ratingGateWay "github.com/ataha1/movie-app/movie/internal/gateway/rating/http"
	httpHandler "github.com/ataha1/movie-app/movie/internal/handler/http"
	"github.com/ataha1/movie-app/pkg/discovery"
	"github.com/ataha1/movie-app/pkg/discovery/consul"
)

const serviceName = "movies"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil{
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil{
		panic(err)
	}
	go func ()  {
		for{
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil{
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	metadataGateway := metadataGateway.New(registry)
	ratingGateWay := ratingGateWay.New(registry)
	svc := movie.New(metadataGateway, ratingGateWay)
	h := httpHandler.New(svc)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil{
		panic(err)
	}
}