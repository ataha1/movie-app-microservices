package main

import (
	"log"
	"net/http"

	"github.com/ataha1/movie-app/rating/internal/controller/rating"
	httphandler "github.com/ataha1/movie-app/rating/internal/handler/http"
	"github.com/ataha1/movie-app/rating/internal/repository/memory"
)

func main() {
	log.Println("Starting rating serviece")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("GET", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil{
		panic(err)
	}
}