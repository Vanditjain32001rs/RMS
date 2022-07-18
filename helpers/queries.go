package helpers

import (
	"RMS/database"
	"RMS/models"
	"RMS/utilities"
	"github.com/elgris/sqrl"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"log"
)

func CreateUser(user *models.Users, createdBy string) (string, error) {
	query := `WITH my_data(name,email,username,password,created_by,role,latitude,longitude) AS
    		  (VALUES ($1,$2,$3,$4,$5::UUID,$6::user_roles,$7,$8)),
     				step_one as (
         				insert into users (name, email, username, password, created_by)
             			SELECT m.name, m.email, m.username, m.password, m.created_by FROM my_data m
             			returning id),
     				step_two as (
         				insert into roles (user_id, role, username)
             			select s1.id,m.role,m.username
             			from step_one s1,my_data m
             			returning user_id)
			  insert
			  into location(user_id, latitude, longitude)
			  select s2.user_id,m.latitude,m.longitude
			  from step_two s2,my_data m
			  RETURNING user_id`

	var userID string
	createErr := database.Data.Get(&userID, query, user.Name, user.Email, user.Username, user.Password, createdBy, user.Role[0], user.Location.Latitude, user.Location.Longitude)

	return userID, createErr
}

func CreateSubAdmins(user *models.SubAdminModel, signedUserID uuid.UUID) (string, error) {
	query := `WITH my_data(name,email,username,password,created_by,role) AS
			  (VALUES ($1,$2,$3,$4,$5::UUID,$6::user_roles)),
					step_one as (
						INSERT INTO users(name,email,username,password,created_by)
						SELECT m.name, m.email, m.username, m.password, m.created_by FROM my_data m
						returning id)
						insert into roles(user_id, role,username)
						select s1.id, m.role, m.username
						from step_one s1, my_data m
						returning user_id`

	var userID string
	createErr := database.Data.Get(&userID, query, user.Name, user.Email, user.Username, user.Password, signedUserID, user.Role[0])

	return userID, createErr
}

func RegisterUser(user *models.UserModel) (string, error) {
	query := `WITH my_data(name,email,username,password,latitude,longitude) AS
    		  (VALUES ($1,$2,$3,$4,$5,$6)),
     				step_one as (
         				insert into users (name, email, username, password)
             			SELECT m.name, m.email, m.username, m.password FROM my_data m
             			returning id),
     				step_two as (
         				insert into roles (user_id, username)
             			select s1.id,m.username
             			from step_one s1,my_data m
             			returning user_id)
			  insert into location(user_id, latitude, longitude)
			  select s2.user_id,m.latitude,m.longitude
			  from step_two s2,my_data m
			  RETURNING user_id`

	var userID string
	registerErr := database.Data.Get(&userID, query, user.Name, user.Email, user.Username, user.Password, user.Location.Latitude, user.Location.Longitude)

	return userID, registerErr
}

func AddUserLocation(userLocation *models.UserLocation, userID string) error {
	query := `INSERT INTO location(user_id,latitude,longitude)
			  VALUES($1,$2,$3)`
	_, addLocationErr := database.Data.Exec(query, userID, userLocation.UserLoc.Latitude, userLocation.UserLoc.Longitude)
	return addLocationErr
}

func GetPassword(username string) (string, error) {

	query := `SELECT password FROM users WHERE username=$1 and archived_at is not null`
	var hashPass string
	getPassErr := database.Data.Get(&hashPass, query, username)

	return hashPass, getPassErr
}

func FetchUsers(createdBy uuid.UUID, pageNo, taskSize int) ([]models.UserFetchModel, error) {

	query := `WITH UserTask AS (SELECT u.id,u.name, u.email, u.username, r.role
			  FROM users u
         		       INNER JOIN roles r on r.user_id = u.id
			  WHERE u.created_by = $1 AND u.archived_at is null)
			  SELECT * from UserTask LIMIT $2 OFFSET $3`

	var user []models.UserFetchModel
	fetchErr := database.Data.Get(&user, query, createdBy, taskSize, pageNo*taskSize)

	return user, fetchErr
}

func FetchUserRole(userIDs []string) ([]models.RoleStruct, error) {

	query := `SELECT user_id,role FROM roles WHERE user_id IN (?) and archived_at is not null`
	var roleArr []models.RoleStruct
	//var userIDs []string
	//for _, user := range user {
	//	userIDs = append(userIDs, user.ID)
	//}
	sqlQuery, args, err := sqlx.In(query, userIDs)
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery = database.Data.Rebind(sqlQuery)
	err = database.Data.Select(&roleArr, sqlQuery, args...)
	return roleArr, err
}

func FetchAllUsers(pageNo, taskSize int) ([]models.UserFetchModel, error) {

	query := `WITH UserTask AS (SELECT u.id,u.name, u.email, u.username
			  FROM users u
			  WHERE u.archived_at is null)
			  SELECT * from UserTask LIMIT $2 OFFSET $3`

	var user []models.UserFetchModel
	fetchErr := database.Data.Select(&user, query, taskSize, pageNo*taskSize)
	var userIDs []string
	for _, user := range user {
		userIDs = append(userIDs, user.ID)
	}
	roleArr, err := FetchUserRole(userIDs)
	if err != nil {
		return nil, err
	}
	var users []models.UserFetchModel
	for _, u := range user {
		for _, r := range roleArr {
			if u.ID == r.UserID {
				u.Role = append(u.Role, r.UserRole)
			}
		}
		users = append(users, u)
	}
	return users, fetchErr
}

func GetRole(username string) (*models.UserRoleID, error) {
	query := `SELECT u.id, r.role, u.created_by
			  FROM roles r
         			   INNER JOIN users u on u.id = r.user_id
			  WHERE u.username = $1 and u.archived_at is not null
			  ORDER BY u.created_at desc
			  LIMIT 1`
	var user models.UserRoleID
	getRoleErr := database.Data.Get(&user, query, username)

	return &user, getRoleErr
}

func GetLocation(users []models.UserFetchModel) ([]models.UsersLocations, error) {

	query := `SELECT user_id,latitude,longitude FROM location WHERE user_id IN (?) and archived_at is not null`
	var userLocation []models.UsersLocations
	var userIDs []string
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	sqlQuery, args, err := sqlx.In(query, userIDs)
	if err != nil {
		log.Fatal(err)
	}
	sqlQuery = database.Data.Rebind(sqlQuery)
	err = database.Data.Select(&userLocation, sqlQuery, args...)
	return userLocation, err
}

func GetAllSubAdmins(pageNo, taskSize int) ([]models.UserFetchModel, error) {

	query := `WITH UserTask AS (SELECT u.id,u.name,u.email,u.username FROM users u where archived_at is not null)
			  SELECT * FROM UserTask LIMIT $2 OFFSET $3`

	var user []models.UserFetchModel
	fetchErr := database.Data.Select(&user, query, taskSize, pageNo*taskSize)

	var userIDs []string
	for _, user := range user {
		userIDs = append(userIDs, user.ID)
	}

	roleArr, err := FetchUserRole(userIDs)
	if err != nil {
		return nil, err
	}
	var users []models.UserFetchModel
	for _, u := range user {
		for _, r := range roleArr {
			if u.ID == r.UserID {
				u.Role = append(u.Role, r.UserRole)
			}
		}
		if utilities.Contains(u.Role, "subadmin") {
			users = append(users, u)
		}
	}

	return users, fetchErr
}

func AddRoleQuery(user *models.AddRoleModel) error {

	query := `INSERT INTO roles(user_id,role,username)
			  VALUES($1,$2,$3)`

	_, addErr := database.Data.Exec(query, user.ID, user.Role, user.Username)

	return addErr
}

func CreateDishes(dishes []models.Dish, restID string, createdBy uuid.UUID, tx *sqlx.Tx) error {
	psql := sqrl.StatementBuilder.PlaceholderFormat(sqrl.Dollar)
	insertQuery := psql.Insert("dishes").Columns("name", "price", "restaurant_id", "created_by")
	for _, dish := range dishes {
		insertQuery.Values(dish.Name, dish.Price, restID, createdBy)
	}
	sql, args, err := insertQuery.ToSql()
	if err != nil {
		log.Printf("CreateDishes : Error in making the query")
		return err
	}
	_, err = tx.Exec(sql, args...)
	if err != nil {
		log.Printf("CreateDishes : Error in Adding Dishes")
		return err
	}

	return nil
}

func CreateRestaurant(rest *models.AddRestaurantModel, createdBy uuid.UUID, tx *sqlx.Tx) (string, error) {
	query := `INSERT INTO restaurant(name, created_by, latitude, longitude) 
			  VALUES ($1,$2,$3,$4) RETURNING id`

	var restID string
	restErr := tx.Get(&restID, query, rest.Name, createdBy, rest.Latitude, rest.Longitude)

	return restID, restErr
}

func FetchAllRestaurant(pageNo, taskSize int) ([]models.FetchRestaurantModel, error) {

	query := `WITH UserTask AS (SELECT name, latitude, longitude
			  FROM restaurant where archived_at is not null)
			  SELECT * from UserTask LIMIT $2 OFFSET $3`
	var restaurant []models.FetchRestaurantModel
	err := database.Data.Select(&restaurant, query, taskSize, pageNo*taskSize)

	return restaurant, err
}

func FetchAllDish(restID string, pageNo, taskSize int) ([]models.Dish, error) {

	query := `WITH UserTask AS (SELECT name,price
              FROM dishes
			  WHERE restaurant_id=$1 and archived_at is not null) 
			  SELECT * from UserTask LIMIT $2 OFFSET $3`
	var dishList []models.Dish
	err := database.Data.Select(&dishList, query, restID, taskSize, pageNo*taskSize)

	return dishList, err
}

func AddDishQuery(dish models.AddDishModel, userID uuid.UUID) error {

	query := `INSERT INTO dishes(name,price,restaurant_id,created_by)
			  VALUES($1,$2,$3,$4)`

	_, addErr := database.Data.Exec(query, dish.Name, dish.Price, dish.RestaurantID, userID)

	return addErr
}

func FetchSpecificRestaurantLocation(restID string) (models.Location, error) {

	query := `SELECT latitude,longitude
			  FROM restaurant
			  WHERE id=$1 and archived_at is not null`

	var restLoc models.Location
	err := database.Data.Get(&restLoc, query, restID)

	return restLoc, err
}

func FetchDistance(user models.DistanceModel) (float64, error) {

	restaurant, restErr := FetchSpecificRestaurantLocation(user.RestaurantID)
	if restErr != nil {
		return -1, restErr
	}
	query := `SELECT to_char(float8 (point($1,$2) <-> point($3,$4)), 'FM999999999.00')`

	var dist float64
	distErr := database.Data.Get(&dist, query, restaurant.Latitude, restaurant.Longitude, user.UserLat, user.UserLng)

	return dist, distErr
}

func UpdateUser(user models.UpdateUsersModel) error {
	query := `UPDATE users SET name=$1, email=$2, username=$3 where id=$4 and archived_at is not null`
	_, err := database.Data.Exec(query, user.Name, user.Email, user.Username, user.ID)

	return err
}

func UpdateDish(dish models.AddDishModel) error {
	query := `UPDATE dishes SET  price=$2 WHERE restaurant_id=$3 and name=$1 and archived_at is not null`
	_, err := database.Data.Exec(query, dish.Name, dish.Price, dish.RestaurantID)

	return err
}
