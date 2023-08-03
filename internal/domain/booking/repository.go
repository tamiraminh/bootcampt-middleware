package booking

import (
	"database/sql"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)


var (
	bookingQueries = struct {
		selectBooking                   string
		insertBooking					string
		updateBooking					string
	}{
		selectBooking: `
			SELECT
				*
			FROM booking where deletedAt is null `,
		insertBooking: `
			INSERT INTO booking (
				ID,
				ClientName,
				PhotographerName,
				Package,
				DateTime,
				Location,
				Status,
				createdAt,
				updatedAt,
				deletedAt
			) VALUES (
				:ID,
				:ClientName,
				:PhotographerName,
				:Package,
				:DateTime,
				:Location,
				:Status,
				:createdAt,
				:updatedAt,
				:deletedAt)`,

		updateBooking: `
		UPDATE booking
		SET
			ClientName = :ClientName,
			PhotographerName = :PhotographerName,
			Package = :Package,
			DateTime = :DateTime,
			Location = :Location,
			Status = :Status,
			createdAt = :createdAt,
			updatedAt = :updatedAt,
			deletedAt = :deletedAt
		WHERE ID = :ID
		`,
	}
)

type BookingRepository interface {
	Create(booking Booking) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (booking Booking, err error)
	Update(booking Booking) (err error)
}


type BookingRepositoryMySQL struct {
	DB *infras.MySQLConn
}

func ProvideBookingRepositoryMySQL(db *infras.MySQLConn) *BookingRepositoryMySQL {
	s := new(BookingRepositoryMySQL)
	s.DB = db
	return s
}


func (r *BookingRepositoryMySQL) Create(booking Booking) (err error) {
	exists, err := r.ExistsByID(booking.Id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "Booking", "already exists")
		logger.ErrorWithStack(err)
		return
	}

	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txCreate(tx, booking); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}

func (r *BookingRepositoryMySQL) ResolveByID(id uuid.UUID) (booking Booking, err error) {
	err = r.DB.Read.Get(
		&booking,
		bookingQueries.selectBooking+" and id = ?",
		id.String())
	if err != nil && err == sql.ErrNoRows {
		err = failure.NotFound("booking")
		logger.ErrorWithStack(err)
		return
	}
	return
}


func (r *BookingRepositoryMySQL) Update(booking Booking) (err error) {
	exists, err := r.ExistsByID(booking.Id)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if !exists {
		err = failure.NotFound("booking")
		logger.ErrorWithStack(err)
		return
	}


	return r.DB.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if err := r.txUpdate(tx, booking); err != nil {
			e <- err
			return
		}

		e <- nil
	})
}



func (r *BookingRepositoryMySQL) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(ID) FROM booking WHERE id = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *BookingRepositoryMySQL) txCreate(tx *sqlx.Tx, booking Booking) (err error) {
	stmt, err := tx.PrepareNamed(bookingQueries.insertBooking)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(booking)
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}


func (r *BookingRepositoryMySQL) txUpdate(tx *sqlx.Tx, booking Booking) (err error) {
	stmt, err := tx.PrepareNamed(bookingQueries.updateBooking)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(booking)
	if  err != nil {
		logger.ErrorWithStack(err)
	}

	return
}







