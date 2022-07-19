package handlers

import (
	"RMS/helpers"
	"RMS/models"
	"RMS/utilities"
	"encoding/json"
	"log"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	user := &models.UserModel{}
	msg := json.NewDecoder(r.Body).Decode(user)
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
	ID, registerErr := helpers.RegisterUser(user)
	if registerErr != nil {
		log.Printf("SignUp : Error in registering the user")
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
