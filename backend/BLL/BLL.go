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
		database: database,
	}
	multiplexer.setRoutes()

	return &multiplexer
}

func (m *mux) setRoutes() {
	log.Println("Setting route to handle...")

	// Serving static website file
	ui_basePath := "/dashboard"
	m.Handle("GET "+ui_basePath+"/", http.StripPrefix(ui_basePath+"/", http.FileServer(http.Dir("../frontend"))))
	m.Handle("GET "+ui_basePath+"/css/", http.StripPrefix(ui_basePath+"/css/", http.FileServer(http.Dir("../frontend/css"))))
	m.Handle("GET "+ui_basePath+"/js/", http.StripPrefix(ui_basePath+"/js/", http.FileServer(http.Dir("../frontend/js"))))

	// api routes
	api_basePath := "/api"
	// Receipt CRUD
	m.HandleFunc("GET "+api_basePath+"/receipts", m.getAllReceipts())
	m.HandleFunc("GET "+api_basePath+"/receipts/{id}", m.getOneReceipt())
	m.HandleFunc("POST "+api_basePath+"/receipts/create", m.createReceipt())
	m.HandleFunc("DELETE "+api_basePath+"/receipts/delete/{id}", m.deleteReceipt())

	//proposal routes
	// Authentication Endpoints
	m.HandleFunc("POST "+api_basePath+"/auth/register", m.registerUser()) //----->always returns id 1 when creation of new register
	m.HandleFunc("POST "+api_basePath+"/auth/login", m.Login())           //--->works

	// User Endpoints
	m.HandleFunc("GET "+api_basePath+"/users", m.getAllUsers())               //--->works
	m.HandleFunc("POST "+api_basePath+"/users/create", m.CreateUser())        //--->works
	m.HandleFunc("PUT "+api_basePath+"/users/update/{id}", m.ModifyUser())    // "/users/update/:id"    --->work
	m.HandleFunc("DELETE "+api_basePath+"/users/delete/{id}", m.DeleteUser()) // "/users/delete/:id"    --->work

	// Rentable Entity Endpoints
	m.HandleFunc("GET "+api_basePath+"/entities", m.getallEntities()) //-----> NOTE:we think it works but cant test frl speculated to only have one value in db
	//m.HandleFunc("/entities/create", createEntity)
	m.HandleFunc("PUT "+api_basePath+"/entities/update/{id}", m.updateEntities()) // "/entities/update/:id"
	// m.HandleFunc("/entities/delete/{id}", deleteEntity) // "/entities/delete/:id"

	// Rental Log Endpoints
	m.HandleFunc("GET "+api_basePath+"/rental-logs", m.getallRentalLogs()) //------> works
	m.HandleFunc("POST "+api_basePath+"/rental-logs/create", m.createRentalLog())
	m.HandleFunc("DELETE "+api_basePath+"/rental-logs/delete/{id}", m.deleteRentalLog()) // "/rental-logs/delete/:id"

	// m.HandleFunc("GET /", m.redirectToDashboard(ui_basePath))

}

// func (m *mux) redirectToDashboard(url string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		http.Redirect(w, r, url, http.StatusOK)
// 	}
// }

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
		input.NotificationPreference = 1

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
		log.Printf(user.Password)
		log.Printf(user_request.Password)
		if user.Password != user_request.Password {
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
			return
		}

		update := map[string]any{}

		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := m.database.FindOneUserID(id)
		if err != nil {
			http.Error(w, "could not find user with id", http.StatusBadRequest)
			log.Printf("could not find user with id: %d", id)
			return
		}

		for k, v := range update {
			switch k {
			case "email":
				user.Email = v.(string)
			case "password":
				user.Password = v.(string)
			case "role":
				user.Role = v.(string)
			case "notification_preference":
				user.NotificationPreference = int(v.(float64))
			}
		}

		log.Printf("User received")

		if err := m.database.ModifyUser(*user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Could not Modify user. Error: %s \n", err.Error())
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
		if err != nil {
			log.Printf("Error while identifying logged in user. Error: %s", err.Error())
			http.Error(w, "internal error identifying logged in user", http.StatusInternalServerError)
			return
		}

		var receipt DAL.Receipt
		if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
			log.Printf("Error while decoding receipt. Error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if receipt.UserID == 0 {
			receipt.UserID = user.ID
		}

		if receipt.Status == "" {
			receipt.Status = DAL.StatusReceiptEnAttente
		}

		log.Printf("%-v object received.\n", receipt)
		
		if err := m.database.CreateReceipt(receipt); err != nil {
			log.Printf("Error while creating receipt in db. Error: %s", err.Error())
			http.Error(w, "internal error, receipt was not created", http.StatusInternalServerError)
			return
		}

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

		update := map[string]any{}

		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Printf("Error while decoding json. Error: %s", err.Error())
			http.Error(w, "Error while decoding json. ", http.StatusBadRequest)
			return
		}

		for k, v := range update {
			switch k {
			case "name":
				Entities.Name = v.(string)
			case "category":
				Entities.Category = v.(string)
			case "description":
				Entities.Description = v.(string)
			case "image_path":
				Entities.ImagePath = v.(string)
			case "is_available":
				Entities.IsAvailable = v.(bool)
			case "price":
				Entities.Price = v.(float64)
			case "pricing_model":
				Entities.PricingModel = v.(string)
			}
		}

		if err := m.database.UpdateRentables(*Entities); err != nil {
			log.Printf("Error sending update to db. Error: %s", err.Error())
			http.Error(w, "Error while decoding json. ", http.StatusInternalServerError)
			return
		}
	}
}

func (m *mux) getallRentalLogs() http.HandlerFunc {
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

		list, err := m.database.FetchAllRentalLog()

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

func (m *mux) createRentalLog() http.HandlerFunc {
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
		var rental DAL.RentalLog

		if err := json.NewDecoder(r.Body).Decode(&rental); err != nil {
			log.Printf("Error while decoding json. Error: %s", err.Error())
			http.Error(w, "Error while decoding json. ", http.StatusBadRequest)
			return
		}

		if err := m.database.CreateRentalLog(rental); err != nil {
			log.Printf(" Error: %s", err.Error())
			http.Error(w, "Could not create RentalLog", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(rental); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error while encoding into json. Error: %s", err.Error())
			return
		}
	}
}

func (m *mux) deleteRentalLog() http.HandlerFunc {
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

		if err := m.database.DeleteRentalLog(DAL.RentalLog{ID: id}); err != nil {
			log.Printf("Error while deleting RentalLog with id: %d. Error: %s", id, err.Error())
			http.Error(w, "Could not delete RentalLog", http.StatusInternalServerError)
		}

		if err := json.NewEncoder(w).Encode(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Error while encoding id: %d into json. Error: %s", id, err.Error())
			return
		}

	}
}
