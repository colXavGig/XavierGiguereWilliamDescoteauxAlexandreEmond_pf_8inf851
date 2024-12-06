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
	m.HandleFunc("GET "+api_basePath+"/receipt", m.getAllReceipts())
	m.HandleFunc("GET "+api_basePath+"/receipt/{id}", m.getOneReceipt())
	m.HandleFunc("POST "+api_basePath+"/receipt", m.createReceipt())
	m.HandleFunc("PATCH "+api_basePath+"/receipt/{id}", m.modifyReceipt())
	m.HandleFunc("DELETE "+api_basePath+"/receipt/{id}", m.deleteReceipt())

	// user path
	// TODO: login route

	// source path
	// TODO: GET route for each source

	// validation path
	// TODO: GET all validation route
	// TODO: GET validation route
	// TODO: POST validation route
	// TODO: PATCH validation route

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

		receipt, err := m.database.FetchOneReceipt(id)
		if err != nil {
			log.Printf("Error while fetching receipt with id: %d. Error: %d", id, err.Error())
			http.Error(w, "error while fetching receipt with id: %d", id)
			return
		}

		if err := json.NewEncoder(w).Encode(receipt); err != nil {
			log.Printf("Error while encoding receipt with id: %d. Error: %s", id, err.Error())
			http.Error(w, "error while encodinf receipt with id: "+id, http.StatusInternalServerError)
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

func (m *mux) modifyReceipt() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // TODO: Implement modify receipt request handler
		http.Error(w, "not implemented", http.StatusNotImplemented)
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

		if err := m.database.DeleteReceipt(DAL.Receipt{Id: id}); err != nil {
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
