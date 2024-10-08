package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ataha1/movie-app/metadata/internal/controller/metadata"
	httpHandler "github.com/ataha1/movie-app/metadata/internal/handler/http"
	"github.com/ataha1/movie-app/metadata/internal/repository/memory"
	"github.com/ataha1/movie-app/pkg/discovery"
	"github.com/ataha1/movie-app/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "post", 8081, "API handler port")
	flag.Parse()
	log.Printf("starting the movie metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil{
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil{
		panic(err)
	}
	go func(){
		for{
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil{
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)
	repo := memory.New()
	svc := metadata.New(repo)
	h := httpHandler.New(svc)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil{
		panic(err)
	}
}