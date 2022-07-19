package handlers

import (
	"RMS/database"
	"RMS/helpers"
	"RMS/models"
	"RMS/utilities"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	user := models.UserModel{}
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("SignIn : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// hash the given password
	pwd, hashErr := utilities.HashPassword(user.Password)
	if hashErr != nil {
		log.Printf("SignUp : Error in hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user.Password = pwd

	// entering user detail in database
	//ID, registerErr := helpers.RegisterUser(user)
	var ID string
	var err error
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		ID, err = helpers.RegisterUser(&user)
		if err != nil {
			log.Printf("SignUp : Error in adding details to user table")
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		var AddRole models.AddRoleModel
		AddRole.ID = ID
		AddRole.Role = "user"
		AddRole.Username = user.Username
		err = helpers.AddRoleQuery(AddRole, tx)
		if err != nil {
			log.Printf("SignUp : Error in adding details to role table")
			w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		var AddLocation *models.UserLocation
		AddLocation.Username = user.Username
		AddLocation.UserLoc = user.Location
		err = helpers.AddUserLocation(AddLocation, ID, tx)
		if err != nil {
			log.Printf("SignUp : Error in creating user location")
			return err
		}
		return err
	})

	if txErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(ID)
	if jsonErr != nil {
		log.Printf("SignUp : Error in encoding to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("SignUp : Error in writing json body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//if err := json.NewEncoder(w).Encode(struct {
	//	ID string `json:"id"`
	//}{ID: ID}); err != nil {
	//}
}
