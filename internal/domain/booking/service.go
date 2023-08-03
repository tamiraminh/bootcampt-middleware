package booking

import (
	"github.com/evermos/boilerplate-go/configs"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
)

type BookingService interface {
	Create(requestFormat BookingRequestFormat) (booking Booking, err error)
	ResolveByID(id uuid.UUID) (booking Booking, err error)
	SoftDelete(id uuid.UUID) (booking Booking, err error)
	Update(id uuid.UUID, requestFormat BookingRequestFormat) (booking Booking, err error)
}

type BookingServiceImpl struct {
	BookingRepository BookingRepository
	Config        *configs.Config
}

func ProvideBookingServiceImpl(bookingRepository BookingRepository, config *configs.Config) *BookingServiceImpl {
	s := new(BookingServiceImpl)
	s.BookingRepository = bookingRepository
	s.Config = config

	return s
}


func (s *BookingServiceImpl) Create(requestFormat BookingRequestFormat) (booking Booking, err error) {
	booking = booking.NewFromRequestFormat(requestFormat)

	if err != nil {
		return booking, failure.BadRequest(err)
	}

	err = s.BookingRepository.Create(booking)

	if err != nil {
		return
	}


	return
}


func (s *BookingServiceImpl) ResolveByID(id uuid.UUID) (booking Booking, err error) {
	booking, err = s.BookingRepository.ResolveByID(id)

	if booking.IsDeleted() {
		return booking, failure.NotFound("booking")
	}

	return
}

func (s *BookingServiceImpl) Update(id uuid.UUID, requestFormat BookingRequestFormat) (booking Booking, err error) {
	booking, err = s.BookingRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = booking.Update(requestFormat)
	if err != nil {
		return
	}

	err = s.BookingRepository.Update(booking)
	return
}


func (s *BookingServiceImpl) SoftDelete(id uuid.UUID) (booking Booking, err error) {
	booking, err = s.BookingRepository.ResolveByID(id)
	if err != nil {
		return
	}

	err = booking.SoftDelete()
	if err != nil {
		return
	}

	err = s.BookingRepository.Update(booking)
	return
}




