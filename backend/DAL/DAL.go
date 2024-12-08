package DAL

//---------------
//--           -- cursor parking --           --
//---------------

import (
	"database/sql"
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

	recette_list = make([]Receipt, 5) // initialise the list with a slice of capacity=5

	for rows.Next() { // if there is a next, continue to iterate
		recette := Receipt{} // init an Empty Receipt

		// try to scan the row and affect each by a struct field
		err := rows.Scan(&recette.ID, &recette.UserID, &recette.TotalAmount, &recette.Status, &recette.CreatedAt, &recette.ApprovedAt)
		if err != nil {
			return nil, err
		}

		recette_list = append(recette_list, recette) // insert the new Receipt at the end of the list
	}

	return recette_list, nil
}

func (this *OracleDB) CreateReceipt(recette Receipt) error {
	_, err := this.Exec("Insert into Receipts(id, UserID, TotalAmount, Status, CreatedAt, ApprovedAt) values(?,?,?,?,?,?)", recette.ID,
		recette.UserID,
		recette.TotalAmount,
		recette.Status,
		recette.CreatedAt,
		recette.ApprovedAt)
	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) DeleteReceipt(recette Receipt) error {
	query := `DELETE FROM Receipts
		  WHERE id = ?`
	_, err := this.Exec(query, recette.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyReceipt(recette Receipt) error {
	_, err := this.Exec("UPDATE Receipts SET(TotalAmount, Status, CreatedAt, ApprovedAt) Where(ID=?)", recette.TotalAmount,
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
		  WHERE id = ?`
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
		err = rows.Scan(&user.ID, &user.Email, &user.Role, &user.Password)
		if err != nil {
			return nil, err
		}

		user_list = append(user_list, user)
	}

	return user_list, nil
}

func (this *OracleDB) CreateUser(user User) error {

	stmt, err := this.Prepare(`INSERT INTO Users(email, role, password, notification_preference) VALUES(?,?,?,?,?)`)

	if err != nil {
		return err
	}
	if _, err := stmt.Exec(user.Email,user.Role,user.Password,user.NotificationPreference); err!=nil {
		return err
	}
	return nil
}

func (this *OracleDB) DeleteUser(user User) error {
	_, err := this.Exec(`DELETE * Users where id=?`, user.ID)

	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyUser(user User) error {
	_, err := this.Exec(`UPDATE * Users SET (role=?, password=?, notification_preference=?) WHERE email=?`, user.Role,
		user.Password,
		user.NotificationPreference,
		user.Email)

	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneUser(email string) (*User, error) {
	user := User{}
	query := `Select * Users Where email=?`
	row := this.QueryRow(query, email)
	if err := row.Scan(&user.ID, &user.Email, &user.Role, &user.Password, &user.NotificationPreference); err != nil {
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
		   SET(price=?,description=?,imagepath=?,isavailable=?)
		   Where id=?
	`

	if _, err := this.Exec(query, item.ID); err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneRentable(id int) (*RentableEntity, error) {
	// NOTE: fonction etait comment out mais c'utiliser dans le BLL
	Rentable := RentableEntity{}

	rows := this.QueryRow("Select * From Rentable_Entities Where id=?", id)

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

	query := "Select * from Rental_Logs"
	rows, err := this.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&Rental.ID, &Rental.EntityID, &Rental.RentalDate, &Rental.StartTime, &Rental.EndTime)
		if err != nil {
			return nil, err
		}

		Rental_list = append(Rental_list, Rental)
	}
	return Rental_list, nil
}

func (this *OracleDB) CreateRentalLog(Rental RentalLog) error {
	_, err := this.Exec("INSERT INTO Rental_Logs(EntityID, RentalDate, StartTime, EndTime) Values(?,?,?,?)", Rental.EntityID,
		Rental.RentalDate,
		Rental.StartTime,
		Rental.EndTime)

	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) DeleteRentalLog(Rental RentalLog) error {
	_, err := this.Exec("DELETE * Rental_Logs Where id=?", Rental.ID)

	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) ModifyRentalLog(Rental RentalLog) error {
	_, err := this.Exec("UPDATE Rental_Logs SET(EntityID=?, RentalDate=?, StartTime=?, EndTime=?) WHERE ID=?", Rental.EntityID,
		Rental.RentalDate,
		Rental.StartTime,
		Rental.EndTime,
		Rental.ID)
	if err != nil {
		return err
	}
	return nil
}

func (this *OracleDB) FindOneRentalLog(id int) (*RentalLog, error) {
	Rental := RentalLog{}

	rows := this.QueryRow("Select * from RentalLogs Where id=?", id)

	if err := rows.Scan(&Rental.ID, &Rental.EntityID, &Rental.RentalDate, &Rental.StartTime, &Rental.EndTime); err != nil {
		return nil, err
	}

	return &Rental, nil
}

//TODO: Create Mapping
