package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/evermos/boilerplate-go/internal/domain/booking"
	"github.com/evermos/boilerplate-go/shared"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/transport/http/middleware"
	"github.com/evermos/boilerplate-go/transport/http/response"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
)

type BookingHandler struct {
	BookingService booking.BookingService
	AuthMiddleware *middleware.Authentication
}

func ProvideBookingHandler(BookingService booking.BookingService, authMiddleware *middleware.Authentication) BookingHandler {
	return BookingHandler{
		BookingService: BookingService,
		AuthMiddleware: authMiddleware,
	}
}

func (h *BookingHandler) Router(r chi.Router) {
	r.Route("/bookings", func(r chi.Router) {
		// r.Use(h.AuthMiddleware.AuthMiddleware)
		r.Use(h.AuthMiddleware.CustomMiddleware)
		r.Post("/", h.CreateBooking)
		r.Get("/{id}", h.ResolveBookingByID)
		r.Get("/", h.ResolveBookingByParams)
		r.Delete("/{id}", h.SoftDeleteBooking)
		r.Put("/{id}", h.UpdateBooking)
		

	})
}


func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestFormat booking.BookingRequestFormat
	err := decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	err = shared.GetValidator().Struct(requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	booking, err := h.BookingService.Create(requestFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}
	w.Header().Set("RequestId", booking.Id.String())
	response.WithJSON(w, http.StatusCreated, booking)
}


func (h *BookingHandler) ResolveBookingByID(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}


	booking, err := h.BookingService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	w.Header().Set("RequestId", booking.Id.String())
	response.WithJSON(w, http.StatusOK, booking)
}


func (h *BookingHandler) ResolveBookingByParams(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}


	booking, err := h.BookingService.ResolveByID(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	w.Header().Set("RequestId", booking.Id.String())
	response.WithJSON(w, http.StatusOK, booking)
}




func (h *BookingHandler) SoftDeleteBooking(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}


	foo, err := h.BookingService.SoftDelete(id)
	if err != nil {
		response.WithError(w, err)
		return
	}
	w.Header().Set("RequestId", idString)
	response.WithJSON(w, http.StatusOK, foo)
}


func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	idString := chi.URLParam(r, "id")
	id, err := uuid.FromString(idString)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var requestFormat booking.BookingRequestFormat
	err = decoder.Decode(&requestFormat)
	if err != nil {
		response.WithError(w, failure.BadRequest(err))
		return
	}


	foo, err := h.BookingService.Update(id, requestFormat)
	if err != nil {
		response.WithError(w, err)
		return
	}
	w.Header().Set("RequestId", idString)
	response.WithJSON(w, http.StatusOK, foo)
}

