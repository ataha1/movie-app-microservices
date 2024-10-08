package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ataha1/movie-app/rating/internal/controller/rating"
	"github.com/ataha1/movie-app/rating/pkg/model"
)

type Handler struct {
	ctrl *rating.Controller
}

func New (ctrl *rating.Controller) *Handler{
	return &Handler{
		ctrl:  ctrl,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request){
	recordID := model.RecordID(r.FormValue("id"))
	if recordID == ""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	recordType := model.RecordType(r.FormValue("type"))
	if recordType == ""{
		w.WriteHeader(http.StatusBadRequest)
	}
	switch r.Method{
	case http.MethodGet:
		v, err := h.ctrl.GetAggregatedRating(r.Context(), recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(v); err != nil{
			log.Printf("Response encode errors: %v\n", err)
		}

	case http.MethodPut:
		userID := model.UserID(r.FormValue("userId"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)
		if err != nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := h.ctrl.PutRating(r.Context(), recordID, recordType, &model.Rating{UserID: userID, Value: model.RatingValue(v)}); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}