package BLL

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/colXavGig/XavierGiguereWilliamDescoteauxAlexandreEmond_pf_8inf851/DAL"
)

//TODO: create BLL

type mux struct {
	*http.ServeMux
	database DAL.Repos // NOTE: We can create an interface for our database and change database type to the interface
}

func NewServer(address string, db_connString string) *http.Server {
	log.Println("Setting up server...")

	multiplexer := newMux(db_connString)
	if multiplexer == nil {
		return nil
	}

	return &http.Server{
		Addr:    address,
		Handler: multiplexer,
	}

}

func newMux(db_connString string) *mux {
	log.Println("Setting up multiplexer...")

	database, err := DAL.NewOracleDB(db_connString)
	if err != nil {
		log.Fatalf("Could not create multiplexer. Error: %s\n", err.Error())
		return nil
	}

	multiplexer := mux{
		ServeMux: http.NewServeMux(),
		database: *&database,
	}
	multiplexer.setRoutes()

	return &multiplexer
}

func (m *mux) setRoutes() {
	log.Println("Setting route to handle...")

	m.HandleFunc("GET /", m.redirectToDashboard())

	// Serving static website file
	ui_basePath := "/dashboard"
	m.Handle("GET "+ui_basePath, http.FileServer(http.Dir("../frontend/")))

	// api routes
	api_basePath := "/api"
	// Receipt CRUD
	m.HandleFunc("GET "+api_basePath+"/receipts", m.getAllReceipts())
	m.HandleFunc("GET "+api_basePath+"/receipts/{id}", m.getOneReceipt())
	m.HandleFunc("POST "+api_basePath+"/receipts/create/{id}", m.createReceipt())
	m.HandleFunc("DELETE "+api_basePath+"/receipts/delete/{id}", m.deleteReceipt())

	// source path
	// TODO: GET route for each source

	// validation path
	// TODO: GET all validation route
	// TODO: GET validation route
	// TODO: POST validation route
	// TODO: PATCH validation route

	//proposal routes
	// Authentication Endpoints
	m.HandleFunc("/auth/register", m.registerUser())
	m.HandleFunc("/auth/login", m.Login())
	m.HandleFunc("/auth/logout", logoutUser())

	// User Endpoints
	m.HandleFunc("/users", m.getAllUsers())
	m.HandleFunc("/users/create", m.CreateUser())
	m.HandleFunc("/users/update/{id}", m.ModifyUser()) // "/users/update/:id"
	m.HandleFunc("/users/delete/{id}", m.DeleteUser()) // "/users/delete/:id"

	// Rentable Entity Endpoints
	m.HandleFunc("/entities", getAllEntities)
	//m.HandleFunc("/entities/create", createEntity)
	m.HandleFunc("/entities/update/{id}", updateEntity) // "/entities/update/:id"
	// m.HandleFunc("/entities/delete/{id}", deleteEntity) // "/entities/delete/:id"

	// Rental Log Endpoints
	m.HandleFunc("/rental-logs", getAllRentalLogs)
	m.HandleFunc("/rental-logs/create", createRentalLog)
	m.HandleFunc("/rental-logs/delete/{id}", deleteRentalLog) // "/rental-logs/delete/:id"

}

func (m *mux) redirectToDashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusOK)
	}
}

/////////////////////////
//  Authentification   //
/////////////////////////

// POST /auth/register
func (m *mux) registerUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input DAL.User

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "could not decode json", http.StatusBadRequest)
			return
		}

		input.Role = DAL.UserRoleClient
		input.NotificationPreference = true

		if err := m.database.CreateUser(input); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userID := 1 // Simulated user ID
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Registration successful", "user_id": userID})
	}
}

func (m *mux) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Println("Create request received")
		var user_request DAL.User

		if err := json.NewDecoder(r.Body).Decode(&user_request); err != nil {
			log.Printf("Error while decoding receipt. Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("GET request received. Requesting login")

		user, err := m.database.FindOneUser(user_request.Email)
		if err != nil {
			log.Printf("Error Username invalid err:\n %s", err.Error())
			http.Error(w, "error while fetching one user", http.StatusInternalServerError)
			return
		}
		if user.Password != user.Password {
			log.Printf("Invalid Password err:\n")

			json.NewEncoder(w).Encode(map[string]string{"error": "invalid credential"})
			// TODO: return http error 401 unauthorized
			return
		} else {
			token, err := createToken(user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := map[string]interface{}{
				"token":     token,
				"user_id":   user.ID,
				"user_role": user.Role,
			}
			json.NewEncoder(w).Encode(data)
			return
		}
	}
}

// ///////////////////
//
//	User       //
//
// ///////////////////
func (m *mux) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println("GET request received. Requesting List of User")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		users, err := m.database.FetchallUser()

		if err != nil {
			log.Printf("Error while fetching all list. Error: %s\n", err.Error())
			http.Error(w, "error while fetching", http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&users); err != nil {
			log.Printf("Could not encode User list Error %s\n", err.Error())
			http.Error(w, "error while encoding", http.StatusInternalServerError)
			return
		}
	}
}

func (m *mux) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		var user DAL.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			log.Printf("Could not decode Body")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("User received")
		m.database.CreateUser(user)

		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Could not create user. Error: %s \n", err.Error())
			return
		}
	}
}

func (m *mux) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "{ \"error\" : \"could not parse id\"}", http.StatusBadRequest)
			return
		}

		log.Printf("User received")
		if err = m.database.DeleteUser(DAL.User{ID: id}); err != nil {
			http.Error(w, "{ \"error\" : \"\"}", http.StatusBadRequest)
		}

	}
}

func (m *mux) ModifyUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "could not parse id", http.StatusBadRequest)
		}

		log.Printf("User received")

		if err := m.database.ModifyUser(DAL.User{ID: id}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Could not Modify user. Error: %s \n", err.Error())
			return
		}
	}
}

// ///////////////
//
//	Receipts  //
//
// ///////////////
func (m *mux) getAllReceipts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println("GET request received. Requesting list of recipt")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		receipts, err := m.database.FetchAllReceipt()
		if err != nil {
			log.Printf("Error while fetching all list. Error: %s\n", err.Error())
			http.Error(w, "error while fetching", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(&receipts); err != nil {
			log.Printf("Could not encode receipts list. Error: %s\n", err.Error())
			http.Error(w, "error while encoding", http.StatusInternalServerError)
			return
		}
	}
}

func (m *mux) getOneReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println("GET request received for")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		list, err := m.database.FetchAllReceipt()
		if err != nil {
			log.Printf("Error while trying to get receipt list. Error: %s\n", err.Error())
			return
		} else if len(list) <= 0 {
			http.Error(w, "No receipt in database", http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Could not parse id: "+r.PathValue("id"), http.StatusBadRequest)
			return
		}

		receipt, err := m.database.FindOneReceipt(id)
		if err != nil {
			log.Printf("Error while fetching receipt with id: %d. Error: %s", id, err.Error())
			http.Error(w, "error while fetching receipt with id: %d", id)
			return
		}

		if err := json.NewEncoder(w).Encode(receipt); err != nil {
			log.Printf("Error while encoding receipt with id: %d. Error: %s", id, err.Error())
			http.Error(w, "error while encoding receipt with id: "+strconv.Itoa(id), http.StatusInternalServerError) //Itoa is the same as FormatInt(int64(i), 10)
			return
		}
	}
}

func (m *mux) createReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // TODO: Implement create receipt request handler
		w.Header().Set("Content-Type", "application/json")

		log.Println("Create request received")

		email, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "could not verify header", http.StatusForbidden)
			return
		}

		user, err := m.database.FindOneUser(email)

		var receipt DAL.Receipt
		if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
			log.Printf("Error while decoding receipt. Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if receipt.UserID != user.ID {
			json.NewEncoder(w).Encode(map[string]string{"error": "user id and token doesnt match"})
			return
		}

		log.Printf("%-v object received.\n", receipt)
		m.database.CreateReceipt(receipt)

		if err := json.NewEncoder(w).Encode(receipt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err.Error())
			return
		}
	}
}

func (m *mux) deleteReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("DELETE receipt request received")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		list, err := m.database.FetchAllReceipt()
		if err != nil {
			log.Printf("Error while trying to get receipt list. Error: %s\n", err.Error())
			return
		} else if len(list) <= 0 {
			http.Error(w, "No receipt in database", http.StatusNotFound)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Printf("An error occured while parsing id: %s\n", err.Error())
			return
		}

		if err := m.database.DeleteReceipt(DAL.Receipt{ID: id}); err != nil {
			log.Printf("Error while deleting receipt with id: %d. Error: %s", id, err.Error())
			http.Error(w, "Could not delete receipt", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error while encoding id: %d into json. Error: %s", id, err.Error())
			return
		}
	}
}

// //////////////
// Entites   //
// ////////////
func (m *mux) getallEntities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Printf("Get request received")

		_, err := verifyToken(r.Header.Get("token"))
		if err != nil {
			http.Error(w, "token verification failed", http.StatusForbidden)
			return
		}

		list, err := m.database.FetchAllRentable()

		if err != nil {
			log.Print("Error couldn't fetch entities Error: &s", err.Error())
			http.Error(w, "error ehile fetching", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(&list); err != nil {
			log.Printf("Could not encode entites list. Error: %s\n", err.Error())
			http.Error(w, "error while encoding", http.StatusInternalServerError)
			return
		}

	}

}

func (m *mux) updateEntities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		email, err := verifyToken(r.Header.Get("Token"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		_, err = m.database.FindOneUser(email)
		if err != nil {
			log.Printf("Error while finding user. Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Error while parsing id. Error: %s", err.Error())
			http.Error(w, "could not parse id", http.StatusBadRequest)
			return
		}

		Entities, err := m.database.FindOneRentable(id)
		if err != nil {
			log.Printf("Error while finding rentable with id: %d. Error: %s", id, err.Error())
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		json.NewDecoder(r.Body).Decode(&Entities)
	}
}

// proposal routes handlers

// ---- Authentication Endpoints ----

// POST /auth/logout
// TODO: check if vaild
func logoutUser() http.HandlerFunc {
	// TODO: implement
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
	}
}

// ---- User Endpoints ----

// GET /users
// TODO: check if vaild
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []DAL.User{{ID: 1, Email: "admin@example.com", Role: "admin"}}
	json.NewEncoder(w).Encode(users)
}

// POST /users/create
// TODO: check if vaild
func createUser(w http.ResponseWriter, r *http.Request) {
	var user DAL.User
	json.NewDecoder(r.Body).Decode(&user)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// PUT /users/update/:id
// TODO: check if vaild
func updateUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/update/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully", "id": id})
}

// DELETE /users/delete/:id
// TODO: check if vaild
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/delete/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully", "id": id})
}

// ---- Rentable Entity Endpoints ----

// GET /entities
// TODO: check if vaild
func getAllEntities(w http.ResponseWriter, r *http.Request) {
	entities := []DAL.RentableEntity{{ID: 1, Name: "Room 1", Category: "Hotel", PricingModel: "per_day", Price: 100}}
	json.NewEncoder(w).Encode(entities)
}

// POST /entities/create
// TODO: check if vaild
func createEntity(w http.ResponseWriter, r *http.Request) {
	var entity DAL.RentableEntity
	json.NewDecoder(r.Body).Decode(&entity)
	json.NewEncoder(w).Encode(map[string]string{"message": "Entity created successfully"})
}

// PUT /entities/update/:id
// TODO: check if vaild
func updateEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/entities/update/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "Entity updated successfully", "id": id})
}

// DELETE /entities/delete/:id
// TODO: check if vaild
func deleteEntity(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/entities/delete/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "Entity deleted successfully", "id": id})
}

// ---- Rental Log Endpoints ----

// GET /rental-logs
// TODO: check if vaild
func getAllRentalLogs(w http.ResponseWriter, r *http.Request) {
	logs := []DAL.RentalLog{{ID: 1, EntityID: 1, UserID: 1, RentalDate: "2024-12-06"}}
	json.NewEncoder(w).Encode(logs)
}

// POST /rental-logs/create
// TODO: check if vaild
func createRentalLog(w http.ResponseWriter, r *http.Request) {
	var log DAL.RentalLog

	json.NewDecoder(r.Body).Decode(&log)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rental log created successfully"})
}

// DELETE /rental-logs/delete/:id
// TODO: check if vaild
func deleteRentalLog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/rental-logs/delete/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "Rental log deleted successfully", "id": id})
}
