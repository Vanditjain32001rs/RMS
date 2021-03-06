package handlers

import (
	"RMS/database"
	"RMS/helpers"
	"RMS/models"
	"RMS/utilities"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"strconv"
)

func AddRestaurant(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
	signedUser := ctx.UserID
	signedUserID, _ := uuid.Parse(signedUser)

	rest := &models.AddRestaurantModel{}
	msg := json.NewDecoder(r.Body).Decode(rest)
	if msg != nil {
		log.Printf("AddRestaurant : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var restID string
	var err error
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		restID, err = helpers.CreateRestaurant(rest, signedUserID, tx)
		if err != nil {
			log.Printf("AddRestaurant: error in creating restaurant")
			//w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		err = helpers.CreateDishes(rest.Dishes, restID, signedUserID, tx)
		if err != nil {
			log.Printf("AddRestaurant : Error in creating the dishes")
			//w.WriteHeader(http.StatusInternalServerError)
			return err
		}
		return err
	})

	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(restID)
	if jsonErr != nil {
		log.Printf("AddRestaurant : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("AddRestaurant : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetRestaurantList(w http.ResponseWriter, r *http.Request) {

	var page models.PageModel
	PageID := r.URL.Query().Get("pageNo")
	page.PageNo, _ = strconv.Atoi(PageID)
	TaskLimit := r.URL.Query().Get("taskLimit")
	page.TaskSize, _ = strconv.Atoi(TaskLimit)

	if page.TaskSize == 0 {
		page.TaskSize = 5
	}

	fetchRestaurant := make([]models.FetchRestaurantModel, 0)
	var fetchErr error

	fetchRestaurant, fetchErr = helpers.FetchAllRestaurant(page.PageNo-1, page.TaskSize)
	if fetchErr != nil {
		log.Printf("GetRestaurantList : Error in fetching all the restaurant")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(fetchRestaurant)
	if jsonErr != nil {
		log.Printf("GetRestaurantList : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("GetRestaurantList : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func GetRestaurantDishList(w http.ResponseWriter, r *http.Request) {

	var page models.PageModel
	PageID := r.URL.Query().Get("pageNo")
	page.PageNo, _ = strconv.Atoi(PageID)
	TaskLimit := r.URL.Query().Get("taskLimit")
	page.TaskSize, _ = strconv.Atoi(TaskLimit)

	if page.TaskSize == 0 {
		page.TaskSize = 5
	}

	//var restaurantID string
	var restID models.Restaurant
	msg := json.NewDecoder(r.Body).Decode(&restID)
	if msg != nil {
		log.Printf("GetRestaurantDishList : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var fetchRestaurantDish []models.Dish
	var fetchDishErr error
	fetchRestaurantDish, fetchDishErr = helpers.FetchAllDish(restID.RestaurantID, page.PageNo-1, page.TaskSize)
	if fetchDishErr != nil {
		log.Printf("GetRestaurantDishList : Error in fetching all the dish")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(fetchRestaurantDish)
	if jsonErr != nil {
		log.Printf("GetRestaurantDishList : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("GetRestaurantDishList : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AddDish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value("User").(models.ContextMap)
	signedUser := ctx.UserID
	signedUserID, _ := uuid.Parse(signedUser)

	var dish models.AddDishModel
	msg := json.NewDecoder(r.Body).Decode(&dish)
	if msg != nil {
		log.Printf("AddDish : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addErr := helpers.AddDishQuery(dish, signedUserID)
	if addErr != nil {
		log.Printf("AddDish : Error in adding the dish. %s", addErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(fmt.Sprintf("Dish Added"))
	if jsonErr != nil {
		log.Printf("AddDish : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("AddDish : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetDistance(w http.ResponseWriter, r *http.Request) {

	var user models.DistanceModel
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("GetDistance : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dist, fetchDistErr := helpers.FetchDistance(user)
	if fetchDistErr != nil {
		log.Printf("GetDistance : Error in fetching the distance. %s", fetchDistErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(dist)
	if jsonErr != nil {
		log.Printf("GetDistance : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("GetDistance : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func UpdateDish(w http.ResponseWriter, r *http.Request) {

	var restaurant models.AddDishModel
	msg := json.NewDecoder(r.Body).Decode(&restaurant)
	if msg != nil {
		log.Printf("UpdateDish : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := helpers.UpdateDish(restaurant)
	if err != nil {
		log.Printf("UpdateDish : Error in updating the dish. %s", err)
	}

	jsonData, jsonErr := utilities.EncodeToJson(fmt.Sprintf("Updated Dish"))
	if jsonErr != nil {
		log.Printf("UpdateDish : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("UpdateDish : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
