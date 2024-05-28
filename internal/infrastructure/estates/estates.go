package estates

import (
	"api/internal/core"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type estates struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *estates {
	return &estates{db: db}
}

func (r *estates) Add(e *core.Estate) (int, error) {
	var id int
	query := "INSERT INTO estates.estates (description, images, amenities, owner_id, price_night, price_week, area, rooms, showers, baby_rooms, category_id, is_public) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;"
	row := r.db.QueryRowx(query, e.Description, e.Images, e.Amenities, e.OwnerId, e.PriceNight, e.PriceWeek, e.Area, e.Rooms, e.Showers, e.BabyRooms, e.CategoryId, true)
	err := row.Scan(&id)
	return id, err
}

func (r *estates) AddTempImage(path string) (int, error) {
	var id int
	query := "INSERT INTO estates.temp_images (path) VALUES ($1) RETURNING id;"
	err := r.db.QueryRowx(query, path).Scan(&id)
	return id, err
}

func (r *estates) GetTempImages(ids []int) ([]string, error) {
	var output []string
	query := "SELECT path FROM estates.temp_images WHERE id = ANY($1);"
	rows, err := r.db.Queryx(query, ids)
	if err != nil {
		return output, err
	}

	for rows.Next() {
		var image string
		err := rows.Scan(&image)
		if err != nil {
			return output, err
		}
		output = append(output, image)
	}

	if err := r.db.QueryRow("DELETE FROM estates.temp_images WHERE id IN ($1);", ids); err != nil {
		log.Println(err)
	}

	return output, nil
}

func (r *estates) GetAmenities() ([]core.Amenity, error) {
	var output []core.Amenity

	query := "SELECT * FROM estates.amenities"
	rows, err := r.db.Queryx(query)
	if err != nil {
		return output, err
	}

	for rows.Next() {
		var item core.Amenity
		err := rows.StructScan(&item)
		if err != nil {
			return output, err
		}
		output = append(output, item)
	}

	return output, nil
}

func (r *estates) AddAddress(a *core.Address) error {
	var id int
	query := "INSERT INTO estates.addresses (estate_id, number, street, city, district) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	return r.db.QueryRowx(query, a.EstateId, a.Number, a.Street, a.City, a.District).Scan(&id)
}

func (r *estates) GetCategories() ([]core.Category, error) {
	var output []core.Category

	query := "SELECT * FROM estates.categories"
	rows, err := r.db.Queryx(query)
	if err != nil {
		return output, err
	}

	for rows.Next() {
		var item core.Category
		err := rows.StructScan(&item)
		if err != nil {
			return output, err
		}
		output = append(output, item)
	}

	return output, nil
}

func (r *estates) GetAll() (*[]core.Estate, error) {
	var output []core.Estate
	query := "SELECT id, images, price_night FROM estates.estates WHERE is_public = true;"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var estate core.Estate
		var arrayStr string
		if err := rows.Scan(&estate.Id, &arrayStr, &estate.PriceNight); err != nil {
			return nil, err
		}
		estate.Images = strings.Split(arrayStr[1:len(arrayStr)-1], ",")

		var addr core.Address
		err := r.db.QueryRowx("SELECT number as address_number, street, city, district FROM estates.addresses WHERE estate_id = $1 LIMIT 1", estate.Id).StructScan(&addr)
		if err != nil {
			log.Println(err)
		}
		estate.Address = addr

		output = append(output, estate)
	}

	return &output, nil
}

func (r *estates) GetOne(id int) (*core.Estate, error) {
	var output core.Estate

	var arrayStr string
	var arrayStr2 string

	query := "SELECT owner_id, description, (SELECT array_agg(name) FROM estates.amenities WHERE id = any(amenities)) as amenities, price_night, price_week, images, rooms, showers, baby_rooms FROM estates.estates WHERE id = $1 AND is_public = true;"
	err := r.db.QueryRowx(query, id).Scan(
		&output.OwnerId, &output.Description, &arrayStr2, &output.PriceNight,
		&output.PriceWeek, &arrayStr, &output.Rooms,
		&output.Showers, &output.BabyRooms,
	)
	if err != nil {
		return nil, err
	}

	output.Images = strings.Split(arrayStr[1:len(arrayStr)-1], ",")
	output.AmenitiesNames = strings.Split(arrayStr2[1:len(arrayStr2)-1], ",")

	var addr core.Address
	err = r.db.QueryRowx("SELECT number as address_number, street, city, district FROM estates.addresses WHERE estate_id = $1 LIMIT 1", id).StructScan(&addr)
	if err != nil {
		return nil, err
	}
	output.Address = addr

	return &output, err
}
