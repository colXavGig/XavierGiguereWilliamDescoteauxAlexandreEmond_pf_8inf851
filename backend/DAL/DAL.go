package DAL

//---------------
//--           -- cursor parking --           --
//---------------

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/godror/godror"
)

// TODO: Create DAL
type OracleDB struct {
	*sql.DB
}

type Repos interface {
	ReceiptRepo
	UserRepo
	RentableRepo
	RentableLogRepo
}

func NewOracleDB(connectionstring string) (*OracleDB, error) {
	db, err := sql.Open("godror", connectionstring)
	if err != nil {
		return nil, err
	}
	return &OracleDB{
		DB: db,
	}, nil
}

// ////////////////
// Receipt Action//
// ////////////////
type ReceiptRepo interface {
	FetchAllReceipt() ([]Receipt, error)
	FindOneReceipt(int) (*Receipt, error)
	CreateReceipt(Receipt) error
	DeleteReceipt(Receipt) error
	ModifyReceipt(Receipt) error
}

func (this *OracleDB) FetchAllReceipt() ([]Receipt, error) {
	var recette_list []Receipt // declare a slice of receipt to contains our rows data

	query := `SELECT *
		  FROM Receipts`
	rows, err := this.Query(query) // try to query the db to get all rows in T_Recette
	if err != nil {
		// error encountered, log it and return it,
		// along side a null pointer for the Receipt slice
		log.Printf("Could not query the db. Error: %s", err.Error())
		return nil, err
	}
	defer rows.Close() // close the connections to the rows (to the transactions) before leaving this scope

	recette_list = []Receipt{}// initialise the list with a slice of capacity=5

	for rows.Next() { // if there is a next, continue to iterate
		recette := Receipt{} // init an Empty Receipt

		// try to scan the row and affect each by a struct field
		err := rows.Scan(&recette.ID, &recette.UserID, &recette.TotalAmount, &recette.Status, &recette.CreatedAt)
		if err != nil {
			return nil, err
		}

		recette_list = append(recette_list, recette) // insert the new Receipt at the end of the list
	}

	return recette_list, nil
}

func (this *OracleDB) CreateReceipt(recette Receipt) error {
	_, err := this.Exec("Insert into Receipts(user_id, total_amount, status) values(:1,:2,:3)", 
		recette.UserID,
		recette.TotalAmount,
		recette.Status,)
	if err != nil {
		return err
	}

	for _, item := range recette.LineItems {
		if err := this.AssociateLineItems(item); err != nil {
			return err
		}
	}
	return nil
}

func (db *OracleDB) AssociateLineItems(lineItems LineItem) error {
	// TODO:
	// TODO: implement see Receipt_Line_Items table
	// TODO:

	return errors.New("Association line items func not implemented") // NOTE: delete line once implemented
}

func (this *OracleDB) DeleteReceipt(recette Receipt) error {
	query := `DELETE FROM Receipts
		  WHERE id = :1`
	_, err := this.Exec(query, recette.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyReceipt(recette Receipt) error {
	_, err := this.Exec("UPDATE Receipts SET(TotalAmount, Status, CreatedAt, ApprovedAt) Where(ID=:1)", recette.TotalAmount,
		recette.Status,
		recette.CreatedAt,
		recette.ApprovedAt,
		recette.ID)
	if err != nil {
		return err
	}
	return nil
}

func (db *OracleDB) FindOneReceipt(id int) (*Receipt, error) {
	recette := Receipt{}

	query := `SELECT *
		  FROM Receipts
		  WHERE id = :1`
	row := db.QueryRow(query, id)

	if err := row.Scan(&recette.UserID, &recette.TotalAmount, &recette.Status, &recette.CreatedAt, &recette.ApprovedAt); err != nil {
		log.Printf("Error while getting receipt with id %d", id)
		return nil, err
	}
	return &recette, nil
}

///////////////
//    User  //
//////////////

type UserRepo interface {
	FetchallUser() ([]User, error)
	FindOneUser(string) (*User, error)
	CreateUser(User) error
	DeleteUser(User) error
	ModifyUser(User) error
	FindOneUserID(int) (*User, error)
}

func (this *OracleDB) FetchallUser() ([]User, error) {
	user := User{}
	user_list := []User{}
	rows, err := this.Query("Select * From Users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Email, &user.Role, &user.Password, &user.NotificationPreference)
		if err != nil {
			return nil, err
		}

		user_list = append(user_list, user)
	}

	return user_list, nil
}

func (this *OracleDB) CreateUser(user User) error {

	stmt, err := this.Prepare("INSERT INTO Users(email,password, role, notification_preference) VALUES(:1,:2,:3,:4)")

	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.Email, user.Password, user.Role, user.NotificationPreference); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) DeleteUser(user User) error {

	stmt, err := this.Prepare("DELETE FROM Users where id=:1")

	if err != nil {
		return err
	}
	log.Print(user.ID)
	if _, err := stmt.Exec(user.ID); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyUser(user User) error {
	stmt, err := this.Prepare(`UPDATE Users SET role=:1, password=:2, notification_preference=:3 WHERE id=:4`)

	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.Role, user.Password, user.NotificationPreference, user.ID); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneUser(email string) (*User, error) {
	user := User{}
	query := `Select * FROM Users Where email=:1`
	row := this.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.NotificationPreference); err != nil {
		return nil, err
	}
	return &user, nil
}
func (this *OracleDB) FindOneUserID(ID int) (*User, error) {
	user := User{}
	query := `Select * FROM Users Where id=:1`
	row := this.QueryRow(query, ID)
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.NotificationPreference); err != nil {
		return nil, err
	}
	return &user, nil
}

////////////////
//  Rentable  //
////////////////

type RentableRepo interface {
	FetchAllRentable() ([]RentableEntity, error)
	UpdateRentables(RentableEntity) error
	FindOneRentable(int) (*RentableEntity, error)
}

func (this *OracleDB) FetchAllRentable() ([]RentableEntity, error) {

	Rentable := RentableEntity{}
	Rentable_list := []RentableEntity{}

	rows, err := this.Query("Select * From Rentable_Entities")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		rows.Scan(&Rentable.ID, &Rentable.Name, &Rentable.Category, &Rentable.PricingModel, &Rentable.Price,
			&Rentable.Description, &Rentable.ImagePath)
		// TODO: calculate is_available

		Rentable_list = append(Rentable_list, Rentable)
	}
	return Rentable_list, nil
}

func (this *OracleDB) UpdateRentables(item RentableEntity) error {

	query := `UPDATE Rentable_Entities 
		   SET(price=:1,description=:2,imagepath=:3,isavailable=:4)
		   Where id=:5
	`

	if _, err := this.Exec(query, item.ID); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneRentable(id int) (*RentableEntity, error) {
	// NOTE: fonction etait comment out mais c'utiliser dans le BLL
	Rentable := RentableEntity{}

	rows := this.QueryRow("Select * From Rentable_Entities Where id=:1", id)

	if err := rows.Scan(&Rentable.ID, &Rentable.Name, &Rentable.Category, &Rentable.PricingModel, &Rentable.Price,
		&Rentable.Price, &Rentable.Description, &Rentable.ImagePath, &Rentable.IsAvailable); err != nil {
		return nil, err
	}
	return &Rentable, nil
}

///////////////////
//  Rentable Log //
///////////////////

type RentableLogRepo interface {
	FetchAllRentalLog() ([]RentalLog, error)
	CreateRentalLog(RentalLog) error
	DeleteRentalLog(RentalLog) error
	ModifyRentalLog(RentalLog) error
	FindOneRentalLog(int) (*RentalLog, error)
}

func (this *OracleDB) FetchAllRentalLog() ([]RentalLog, error) {
	Rental := RentalLog{}
	Rental_list := []RentalLog{}

	query := "Select id, entity_id, user_id, rental_date from Rental_Logs"
	rows, err := this.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&Rental.ID, &Rental.EntityID, &Rental.UserID, &Rental.RentalDate)
		if err != nil {
			return nil, err
		}

		Rental_list = append(Rental_list, Rental)
	}
	return Rental_list, nil
}

func (this *OracleDB) CreateRentalLog(Rental RentalLog) error {
	stmt, err := this.Prepare("INSERT INTO Rental_Logs(entity_id, rental_date, user_id) Values(:1,TO_DATE(:2, 'YYYY-MM-DD'),:3)") // NOTE: start_time and end_time not used, add if you need it

	if err != nil {
		return err
	}
	if _, err := stmt.Exec(Rental.EntityID, Rental.RentalDate, Rental.UserID); err != nil {
		return err
	}

	return nil
}

func (this *OracleDB) DeleteRentalLog(Rental RentalLog) error {
	stmt, err := this.Prepare("DELETE FROM Rental_Logs Where id=:1")

	if err != nil {
		return err
	}
	if _, err := stmt.Exec(Rental.ID); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyRentalLog(Rental RentalLog) error {
	_, err := this.Exec("UPDATE Rental_Logs SET(EntityID=:1, RentalDate=:2) WHERE ID=:3", Rental.EntityID, // NOTE: start_time and end_time not used, add if you need it
		Rental.RentalDate,
		Rental.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneRentalLog(id int) (*RentalLog, error) {
	Rental := RentalLog{}

	rows := this.QueryRow("Select * from RentalLogs Where id=:1", id)

	if err := rows.Scan(&Rental.ID, &Rental.EntityID, &Rental.RentalDate, &Rental.StartTime, &Rental.EndTime); err != nil {
		return nil, err
	}

	return &Rental, nil
}

//TODO: Create Mapping
