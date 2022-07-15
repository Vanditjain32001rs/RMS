package handlers

import (
	"RMS/helpers"
	"RMS/models"
	"RMS/utilities"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strconv"
)

func AddNewUser(w http.ResponseWriter, r *http.Request) {
	//log.Printf("AddNewUser : reached")
	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserRole := ctx.UserRole
	signedUserID := ctx.UserID

	user := &models.Users{}
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("AddUser : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if (signedUserRole == "subadmin" || signedUserRole == "user") && (utilities.Contains(user.Role, "subadmin") || utilities.Contains(user.Role, "admin")) {
		log.Printf("AddNewUser : %s cannot make %s", signedUserRole, user.Role)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passHash, hashErr := utilities.HashPassword(user.Password)
	if hashErr != nil {
		log.Printf("AddNewUser : Error in hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = passHash

	ID, addErr := helpers.CreateUser(user, signedUserID)
	if addErr != nil {
		log.Printf("AddUser : Error in creating the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(ID)
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

func AddNewLocation(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserRole := ctx.UserRole
	signedUser := ctx.UserID
	signedUserID, uuidErr := uuid.Parse(signedUser)
	if uuidErr != nil {
		log.Printf("AddNewLocation : Error in converting the userid string to uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userLocation := &models.UserLocation{}
	msg := json.NewDecoder(r.Body).Decode(&userLocation)
	if msg != nil {
		log.Printf("AddNewLocation : Error in decoding the json body")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, getIdErr := helpers.GetRole(userLocation.Username)
	if getIdErr != nil {
		log.Printf("AddNewLocation : Error in retrieving user role an id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user.UserRole == "user" {
		if signedUserRole == "subadmin" && user.CreatedBy != signedUserID {
			log.Printf("AddNewLocation : subadmin's can only add location for the user's they created")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		addLocationErr := helpers.AddUserLocation(userLocation, user.UserID)
		if addLocationErr != nil {
			log.Printf("AddNewLocation : Error in add location query")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonData, jsonErr := utilities.EncodeToJson(fmt.Sprintf("New Location Added"))
		if jsonErr != nil {
			log.Printf("AddNewLocation : Error in encoding to json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, wErr := w.Write(jsonData)
		if wErr != nil {
			log.Printf("AddNewLocation : Error in writing the json data")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := fmt.Sprint("Admins and subadmin's does not have location")
		jsonData, jsonErr := utilities.EncodeToJson(err)
		if jsonErr != nil {
			log.Printf("AddNewLocation : Error in encoding %s to json ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, wErr := w.Write(jsonData)
		if wErr != nil {
			log.Printf("AddNewLocation : Error in writing the json data %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func AddSubAdmins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserRole := ctx.UserRole
	signedUser := ctx.UserID
	signedUserID, uuidErr := uuid.Parse(signedUser)
	if uuidErr != nil {
		log.Printf("AddSubAdmins : Error in converting string to uuid for the signed user")
	}

	if signedUserRole == "subadmin" {
		log.Printf("AddSubAdmins : Subadmins cannot create other subadmins")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &models.SubAdminModel{}
	msg := json.NewDecoder(r.Body).Decode(user)
	if msg != nil {
		log.Printf("AddSubAdmins : Error in decoding the json body %s", msg)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passHash, hashErr := utilities.HashPassword(user.Password)
	if hashErr != nil {
		log.Printf("AddSubAdmins : Error in hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = passHash

	userID, createErr := helpers.CreateSubAdmins(user, signedUserID)
	if createErr != nil {
		log.Printf("AddSubAdmins : Error in creating subadmins %s", createErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := utilities.EncodeToJson(userID)
	if jsonErr != nil {
		log.Printf("AddSubAdmins : Error in encoding to json ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("AddSubAdmins : Error in writing the json data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func FetchUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserRole := ctx.UserRole
	signedUser := ctx.UserID
	signedUserID, _ := uuid.Parse(signedUser)

	var page models.PageModel
	PageID := r.URL.Query().Get("pageNo")
	page.PageNo, _ = strconv.Atoi(PageID)
	TaskLimit := r.URL.Query().Get("taskLimit")
	page.TaskSize, _ = strconv.Atoi(TaskLimit)
	if page.TaskSize == 0 {
		page.TaskSize = 5
	}

	users := make([]models.UsersDetail, 0)
	fetchUser := make([]models.UserFetchModel, 0)
	var fetchErr error

	if signedUserRole == "subadmin" {

		fetchUser, fetchErr = helpers.FetchUsers(signedUserID, page.PageNo-1, page.TaskSize)
		if fetchErr != nil {
			log.Printf("FetchUsers : Error in fetching the users made by subadmin")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if signedUserRole == "admin" {
		fetchUser, fetchErr = helpers.FetchAllUsers(page.PageNo-1, page.TaskSize)
		if fetchErr != nil {
			log.Printf("FetchUsers : Error in fetching all the users %s", fetchErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//log.Printf("FetchUsers : Starting mapping arrays.")
	for _, userInfo := range fetchUser {
		var temp models.UsersDetail
		temp.ID = userInfo.ID
		temp.Name = userInfo.Name
		temp.Email = userInfo.Email
		temp.Username = userInfo.Username
		temp.Role = userInfo.Role
		users = append(users, temp)

	}
	//log.Printf("FetchUsers : Mapped simple details of user.")
	userLocations, locationErr := helpers.GetLocation(fetchUser)
	if locationErr != nil {
		log.Printf("FetchUsers : Error in fetching users locations")
	}
	var returnUser []models.UsersDetail
	for _, user := range users {
		for _, userAdd := range userLocations {
			if user.ID == userAdd.UserID {
				//log.Printf(user.ID)
				var tmp models.Location
				tmp.Latitude = userAdd.Latitude
				tmp.Longitude = userAdd.Longitude
				user.Location = append(user.Location, tmp)
			}
		}
		returnUser = append(returnUser, user)
	}

	jsonData, jsonErr := utilities.EncodeToJson(returnUser)
	if jsonErr != nil {
		log.Printf("AddNewLocation : Error in encoding to json ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("AddNewLocation : Error in writing the json data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func FetchAllSubAdmins(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserRole := ctx.UserRole

	if signedUserRole == "subadmin" || signedUserRole == "user" {
		log.Printf("FetchAllSubAdmins : subadmins or users cannot fetch all the subadmins")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var page models.PageModel
	PageID := r.URL.Query().Get("pageNo")
	page.PageNo, _ = strconv.Atoi(PageID)
	TaskLimit := r.URL.Query().Get("taskLimit")
	page.TaskSize, _ = strconv.Atoi(TaskLimit)

	if page.TaskSize == 0 {
		page.TaskSize = 5
	}

	subAdminDetails, fetchErr := helpers.GetAllSubAdmins(page.PageNo-1, page.TaskSize)
	if fetchErr != nil {
		log.Printf("FetchAllSubAdmins : Error in fetching all the subadmins. %s", fetchErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, jsonErr := utilities.EncodeToJson(subAdminDetails)
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

func AddRole(w http.ResponseWriter, r *http.Request) {

	user := &models.AddRoleModel{}
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("AddRole : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	addErr := helpers.AddRoleQuery(user)
	if addErr != nil {
		log.Printf("AddRole : Error in adding role. %s", addErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(fmt.Sprintf("Role added"))
	if jsonErr != nil {
		log.Printf("AddRole : Error in encoding to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("AddRole : Error in writing the jsonData")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
