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

func AddNewUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
	signedUserID := ctx.UserID

	user := &models.Users{}
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("AddUser : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//v := validator.New()
	//vErr := v.Struct(user)
	//if vErr != nil {
	//	log.Printf("AddNewUser : Error in validating the details entered")
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	passHash, hashErr := utilities.HashPassword(user.Password)
	if hashErr != nil {
		log.Printf("AddNewUser : Error in hashing the password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = passHash

	var ID string
	var err error
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		ID, err = helpers.CreateUser(user, signedUserID, tx)
		if err != nil {
			log.Printf("AddNewUser: error in creating user")
			//w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		var AddRole models.AddRoleModel
		AddRole.ID = ID
		AddRole.Username = user.Username
		AddRole.Role = user.Role[0]
		err = helpers.AddRoleQuery(AddRole, tx)
		if err != nil {
			log.Printf("AddNewUser : Error in creating user role")
			//w.WriteHeader(http.StatusInternalServerError)
			return err
		}

		var AddLocation *models.UserLocation
		AddLocation.Username = user.Username
		AddLocation.UserLoc = user.Location
		err = helpers.AddUserLocation(AddLocation, ID, tx)
		if err != nil {
			log.Printf("AddNewUser : Error in creating user location")
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
	var addLocationErr error
	if user.UserRole == "user" {

		txErr := database.Tx(func(tx *sqlx.Tx) error {
			addLocationErr = helpers.AddUserLocation(userLocation, user.UserID, tx)
			if addLocationErr != nil {
				log.Printf("AddNewLocation : Error in add location query")
				w.WriteHeader(http.StatusInternalServerError)
				return addLocationErr
			}
			return addLocationErr
		})
		if txErr != nil {
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

func FetchUsersByAdmin(w http.ResponseWriter, r *http.Request) {

	var page models.PageModel
	PageID := r.URL.Query().Get("pageNo")
	page.PageNo, _ = strconv.Atoi(PageID)
	TaskLimit := r.URL.Query().Get("taskLimit")
	page.TaskSize, _ = strconv.Atoi(TaskLimit)
	if page.TaskSize == 0 {
		page.TaskSize = 5
	}

	users := make([]models.UsersDetail, 0)

	fetchUser, fetchErr := helpers.FetchAllUsers(page.PageNo-1, page.TaskSize)
	if fetchErr != nil {
		log.Printf("FetchUsersByAdmin : Error in fetching all the users %s", fetchErr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, userInfo := range fetchUser.User {
		var temp models.UsersDetail
		temp.ID = userInfo.ID
		temp.Name = userInfo.Name
		temp.Email = userInfo.Email
		temp.Username = userInfo.Username
		temp.Role = userInfo.Role
		users = append(users, temp)

	}

	userLocations, locationErr := helpers.GetLocation(fetchUser.User)
	if locationErr != nil {
		log.Printf("FetchUsersByAdmin : Error in fetching users locations")
	}

	returnUser := make([]models.UsersDetail, 0)
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

	var user models.UserFetchAdmin
	user.User = returnUser
	user.TotalCount = fetchUser.User[0].TotalCount
	jsonData, jsonErr := utilities.EncodeToJson(user)
	if jsonErr != nil {
		log.Printf("FetchUsersByAdmin : Error in encoding to json ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("FetchUsersByAdmin : Error in writing the json data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func FetchUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context().Value("User").(models.ContextMap)
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

	fetchUser, fetchErr := helpers.FetchUsers(signedUserID, page.PageNo-1, page.TaskSize)
	if fetchErr != nil {
		log.Printf("FetchUsersBySubadmin : Error in fetching the users made by subadmin")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, userInfo := range fetchUser.User {
		var temp models.UsersDetail
		temp.ID = userInfo.ID
		temp.Name = userInfo.Name
		temp.Email = userInfo.Email
		temp.Username = userInfo.Username
		temp.Role = userInfo.Role
		users = append(users, temp)

	}

	userLocations, locationErr := helpers.GetLocation(fetchUser.User)
	if locationErr != nil {
		log.Printf("FetchUsersBySubAdmin : Error in fetching users locations")
	}
	var returnUser []models.UsersDetail
	for _, user := range users {
		for _, userAdd := range userLocations {
			if user.ID == userAdd.UserID {
				var tmp models.Location
				tmp.Latitude = userAdd.Latitude
				tmp.Longitude = userAdd.Longitude
				user.Location = append(user.Location, tmp)
			}
		}
		returnUser = append(returnUser, user)
	}
	var user models.UserFetchAdmin
	user.TotalCount = fetchUser.User[0].TotalCount
	user.User = returnUser

	jsonData, jsonErr := utilities.EncodeToJson(user)
	if jsonErr != nil {
		log.Printf("FetchUsersBySubAdmin : Error in encoding to json ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("FetchUsersBySubadmin : Error in writing the json data")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func FetchAllSubAdmins(w http.ResponseWriter, r *http.Request) {

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
	subAdminDetails.TotalCount = subAdminDetails.User[0].TotalCount
	//totalCount := fmt.Sprintf("Total Count : %d", subAdminDetails.TotalCount)
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

	user := models.AddRoleModel{}
	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("AddRole : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//addErr := helpers.AddRoleQuery(user)
	var addErr error

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		addErr = helpers.AddRoleQuery(user, tx)
		if addErr != nil {
			log.Printf("AddNewLocation : Error in add location query")
			w.WriteHeader(http.StatusInternalServerError)
			return addErr
		}
		return addErr
	})
	if txErr != nil {
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

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user models.UpdateUsersModel

	msg := json.NewDecoder(r.Body).Decode(&user)
	if msg != nil {
		log.Printf("UpdateUser : Error in decoding the json body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := helpers.UpdateUser(user)
	if err != nil {
		log.Printf("UpdateUser(admin) : Error in updating the user. %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonData, jsonErr := utilities.EncodeToJson(fmt.Sprintf("Updated User"))
	if jsonErr != nil {
		log.Printf("UpdateUser : Error in encoding to json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, wErr := w.Write(jsonData)
	if wErr != nil {
		log.Printf("UpdateUser : Error in writing the jsonData")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
