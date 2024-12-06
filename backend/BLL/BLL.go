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
	database *DAL.OracleDB // NOTE: We can create an interface for our database and change database type to the interface
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

	// Serving static website file
	ui_basePath := "/"
	m.Handle("GET "+ui_basePath, http.FileServer(http.Dir("../frontend/")))

	// api routes
	api_basePath := "/api"
	// Receipt CRUD
	m.HandleFunc("GET "+api_basePath+"/receipts", m.getAllReceipts())
	m.HandleFunc("GET "+api_basePath+"/receipts/{id}", m.getOneReceipt())
	m.HandleFunc("POST "+api_basePath+"/receipts/create/{id}", m.createReceipt())
	m.HandleFunc("DELETE "+api_basePath+"/receipts/delete/{id}", m.deleteReceipt())

	// user path
	// TODO: login route
	
	// source path
	// TODO: GET route for each source

	// validation path
	// TODO: GET all validation route
	// TODO: GET validation route
	// TODO: POST validation route
	// TODO: PATCH validation route

	//proposal routes
		// Authentication Endpoints
		m.HandleFunc("/auth/register", registerUser)
		m.HandleFunc("/auth/login", loginUser)
		m.HandleFunc("/auth/logout", logoutUser)

		// User Endpoints
		m.HandleFunc("/users", getAllUsers)
		m.HandleFunc("/users/create", createUser)
		m.HandleFunc("/users/update/", updateUser) // "/users/update/:id"
		m.HandleFunc("/users/delete/", deleteUser) // "/users/delete/:id"

		// Rentable Entity Endpoints
		m.HandleFunc("/entities", getAllEntities)
		m.HandleFunc("/entities/create", createEntity)
		m.HandleFunc("/entities/update/", updateEntity) // "/entities/update/:id"
		m.HandleFunc("/entities/delete/", deleteEntity) // "/entities/delete/:id"

		// Rental Log Endpoints
		m.HandleFunc("/rental-logs", getAllRentalLogs)
		m.HandleFunc("/rental-logs/create", createRentalLog)
		m.HandleFunc("/rental-logs/delete/", deleteRentalLog) // "/rental-logs/delete/:id"

		// Receipt Endpoints
		m.HandleFunc("/receipts", getAllReceipts)
		m.HandleFunc("/receipts/create", createReceipt)
		m.HandleFunc("/receipts/update/", updateReceipt) // "/receipts/update/:id"
		m.HandleFunc("/receipts/delete/", deleteReceipt) // "/receipts/delete/:id"

}
func(m *mux) Login() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {}
}



func (m *mux) getAllReceipts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Println("GET request received. Requesting list of reciept")

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
			log.Printf("Error while fetching receipt with id: %d. Error: %d", id, err.Error())
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
		var receipt DAL.Receipt
		if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
			log.Printf("Error while decoding receipt. Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
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
		log.Println("DEETE receipt request received")

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
			log.Print("Error while encoding id: %d into json. Error: %s", id, err.Error())
			return
		}
	}
}



// proposal routes handlers

// ---- Authentication Endpoints ----

// POST /auth/register
func registerUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input DAL.User
		
		json.NewDecoder(r.Body).Decode(&input)
		userID := 1 // Simulated user ID
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Registration successful", "user_id": userID})
	}
}

// POST /auth/login
// TODO: check if vaild
func loginUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&input)
	json.NewEncoder(w).Encode(map[string]interface{}{"token": "mock_token", "user_role": "user"})
}

// POST /auth/logout
// TODO: check if vaild
func logoutUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// ---- User Endpoints ----

// GET /users
// TODO: check if vaild
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{{ID: 1, Email: "admin@example.com", Role: "admin"}}
	json.NewEncoder(w).Encode(users)
}

// POST /users/create
// TODO: check if vaild
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
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
	entities := []RentableEntity{{ID: 1, Name: "Room 1", Category: "Hotel", PricingModel: "per_day", Price: 100}}
	json.NewEncoder(w).Encode(entities)
}

// POST /entities/create
// TODO: check if vaild
func createEntity(w http.ResponseWriter, r *http.Request) {
	var entity RentableEntity
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
	logs := []RentalLog{{ID: 1, EntityID: 1, UserID: 1, RentalDate: "2024-12-06"}}
	json.NewEncoder(w).Encode(logs)
}

// POST /rental-logs/create
// TODO: check if vaild
func createRentalLog(w http.ResponseWriter, r *http.Request) {
	var log RentalLog
	json.NewDecoder(r.Body).Decode(&log)
	json.NewEncoder(w).Encode(map[string]string{"message": "Rental log created successfully"})
}

// DELETE /rental-logs/delete/:id
// TODO: check if vaild
func deleteRentalLog(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/rental-logs/delete/"):]
	json.NewEncoder(w).Encode(map[string]string{"message": "Rental log deleted successfully", "id": id})
}
