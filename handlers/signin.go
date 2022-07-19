package handlers

import (
	"RMS/helpers"
	"RMS/models"
	"RMS/utilities"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {

	cred := models.Credentials{}
	msg := json.NewDecoder(r.Body).Decode(&cred)
	if msg != nil {
		log.Printf("SignIn : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//retrieve password from database to compare
	pass, getPassErr := helpers.GetPassword(cred.Username)
	if getPassErr != nil {
		log.Printf("Signin : Error in retreiving the password from database.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//comparing the passwords
	if compareErr := bcrypt.CompareHashAndPassword([]byte(pass), []byte(cred.Password)); compareErr != nil {
		log.Printf("Signin : Error in comparing the passwords.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// fetching user role, id and created_by
	user, getIdErr := helpers.GetRole(cred.Username)
	if getIdErr != nil {
		log.Printf("Signin : Error in retreiving the User role.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	validToken, tokenErr := helpers.TokenGeneration(cred.Username, user)
	if tokenErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var tokenMap = make(map[string]string)
	tokenMap["Token"] = validToken

	jsonData, jsonErr := utilities.EncodeToJson(tokenMap)
	if jsonErr != nil {
		log.Printf("CreateUser : Error in creating json file.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("Write Error : Error in writing json data.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
