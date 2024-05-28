package estates

type CreateEstateDto struct {
	Description string `json:"description"`
	Images      []int  `json:"images"`
	Amenities   []int  `json:"amenities"`
	PriceNight  int    `json:"price_night"`
	PriceWeek   int    `json:"price_week"`
	Area        int    `json:"area"`
	Rooms       int    `json:"rooms"`
	Showers     int    `json:"showers"`
	BabyRooms   int    `json:"baby_rooms"`
	CategoryId  int    `json:"category_id"`
	Address     struct {
		Number   int    `json:"address_number" db:"address_number"`
		Street   string `json:"street" db:"street"`
		City     string `json:"city" db:"city"`
		District string `json:"district" db:"district"`
	} `json:"address"`
}
