package core

type Estate struct {
	Id             int      `json:"id" db:"id"`
	Description    string   `json:"description" db:"description"`
	Amenities      []int    `json:"amenities" db:"amenities"`
	AmenitiesNames []string `json:"amenities_names"`
	OwnerId        int      `json:"owner_id" db:"owner_id"`
	PriceNight     int      `json:"price_night" db:"price_night"`
	PriceWeek      int      `json:"price_week" db:"price_week"`
	Area           int      `json:"area" db:"area"`
	Rooms          int      `json:"rooms" db:"rooms"`
	Showers        int      `json:"showers" db:"showers"`
	BabyRooms      int      `json:"baby_rooms" db:"baby_rooms"`
	CategoryId     int      `json:"category_id" db:"category_id"`
	CreatedAt      string   `json:"created_at" db:"created_at"`
	Images         []string `json:"images" db:"images"`
	Address        Address  `json:"address"`
}

type Image struct {
	Id       int    `json:"image_id" db:"image_id"`
	EstateId int    `json:"estate_id" db:"estate_id"`
	Path     string `json:"path" db:"path"`
	Name     string `json:"name" db:"name"`
	IsTemp   bool   `json:"is_temp" db:"is_temp"`
	Priority string `json:"priority" db:"priority"`
}

type Address struct {
	Id       int    `json:"address_id" db:"address_id"`
	EstateId int    `json:"estate_id" db:"estate_id"`
	Number   int    `json:"address_number" db:"address_number"`
	Street   string `json:"street" db:"street"`
	City     string `json:"city" db:"city"`
	District string `json:"district" db:"district"`
}

type Amenity struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Category struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
