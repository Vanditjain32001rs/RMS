package models

type Dish struct {
	Name  string  `json:"name" db:"name"`
	Price float64 `json:"price" db:"price"`
}

type AddRestaurantModel struct {
	Name      string  `json:"name" db:"name"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	Dishes    []Dish  `json:"dishes"`
}

type FetchRestaurantModel struct {
	Name      string  `json:"name" db:"name"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
}

type AddDishModel struct {
	RestaurantID string  `json:"restaurantID" db:"restaurant_id"`
	Name         string  `json:"name" db:"name"`
	Price        float64 `json:"price" db:"price"`
}

type DistanceModel struct {
	RestaurantID string  `json:"restaurantID" db:"restaurant_id"`
	UserLat      float64 `json:"userLatitude" db:"latitude"`
	UserLng      float64 `json:"userLongitude" db:"longitude"`
}

type Restaurant struct {
	RestaurantID string `json:"restaurantID" db:"restaurant_id"`
}
