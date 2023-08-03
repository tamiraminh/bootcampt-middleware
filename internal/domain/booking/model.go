package booking

import (
	"encoding/json"
	"time"

	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/gofrs/uuid"
	"github.com/guregu/null"
)

type Booking struct {
	Id               uuid.UUID 	`db:"ID"`
	ClientName       string		`db:"ClientName"`
	PhotographerName string		`db:"PhotographerName"`
	Package          string		`db:"Package"`
	DateTime         string		`db:"DateTime"`
	Location         string		`db:"Location"`
	Status           string		`db:"Status"`
	CreatedAt       time.Time   `db:"createdAt"`
	UpdatedAt       null.Time   `db:"updatedAt"`
	DeletedAt       null.Time   `db:"deletedAt"`

}

func (b *Booking) IsDeleted() (deleted bool) {
	return b.DeletedAt.Valid 
}

func (b Booking) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.ToResponseFormat())
}


func (b *Booking) Update(req BookingRequestFormat) (err error) {


	b.ClientName = req.ClientName
	b.PhotographerName = req.PhotographerName
	b.Package = req.Package
	b.DateTime = req.DateTime
	b.Location = req.Location
	b.Status = req.Status
	b.UpdatedAt = null.TimeFrom(time.Now())
	

	return
}



func (b *Booking) SoftDelete() (err error) {
	if b.IsDeleted() {
		return failure.Conflict("softDelete", "foo", "already marked as deleted")
	}

	b.DeletedAt = null.TimeFrom(time.Now())

	return
}


func (f Booking) NewFromRequestFormat(req BookingRequestFormat) (newBooking Booking) {
	bookingID, _ := uuid.NewV4()
	newBooking = Booking{
		Id: bookingID,
		ClientName: req.ClientName,
		PhotographerName: req.PhotographerName,
		Package: req.Package,
		DateTime: req.DateTime,
		Location: req.Location,
		Status: req.Status,
		CreatedAt: time.Now(),
	}

	return
}


func (b Booking) ToResponseFormat() BookingResponseFormat {
	resp := BookingResponseFormat{
		Id: b.Id,
		ClientName: b.ClientName,
		PhotographerName: b.PhotographerName,
		Package: b.Package,
		DateTime: b.DateTime,
		Location: b.Location,
		Status: b.Status,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
		DeletedAt: b.DeletedAt,
	}

	return resp
}



type BookingRequestFormat struct {
	ClientName       string `json:"clientName"`
	PhotographerName string `json:"photographerName"`
	Package          string `json:"package"`
	DateTime         string `json:"dateTime"`
	Location         string `json:"location"`
	Status           string `json:"status"`
}

type BookingResponseFormat struct {
	Id               uuid.UUID 	`json:"id"`
	ClientName       string 	`json:"clientName"`
	PhotographerName string 	`json:"photographerName"`
	Package          string 	`json:"package"`
	DateTime         string 	`json:"dateTime"`
	Location         string 	`json:"location"`
	Status           string 	`json:"status"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       null.Time   `json:"updatedAt,omitempty"`
	DeletedAt       null.Time   `json:"deletedAt,omitempty"`
}

