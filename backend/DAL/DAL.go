package DAL

//---------------
//--           -- cursor parking --           --
//---------------

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/godror/godror"
)

//TODO: Create DAL
type OracleDB struct {
	*sql.DB
}



func NewOracleDB(connectionstring string) (*OracleDB,error) {
	db, err := sql.Open("godror", connectionstring)
	if err != nil {
		return nil,err
	}
	return &OracleDB{
		DB: db,
	},nil
}

//////////////////
//Receipt Action//
//////////////////

func (this *OracleDB)FetchAllReceipt() ([]Receipt, error) {
	var recette_list []Receipt // declare a slice of receipt to contains our rows data

	query := `SELECT *
		  FROM T_Recette`
	rows,err := this.Query(query) // try to query the db to get all rows in T_Recette
	if err != nil {
		// error encountered, log it and return it, 
		// along side a null pointer for the Receipt slice
		log.Printf("Could not query the db. Error: %s",err.Error())
		return nil, err
	}
	defer rows.Close() // close the connections to the rows (to the transactions) before leaving this scope

	recette_list = make([]Receipt, 5) // initialise the list with a slice of capacity=5

	for rows.Next()  { // if there is a next, continue to iterate
		recette := Receipt{} // init an Empty Receipt

		// try to scan the row and affect each by a struct field
		err := rows.Scan(&recette.Id, &recette.Total, &recette.DATE,&recette.Statut , &recette.Utilisateur_ID)
		if err != nil {
			return nil, err
		}

		recette_list := append(recette_list, recette) // insert the new Receipt at the end of the list
	}

	return recette_list, nil
}

func (this *OracleDB) CreateReceipt(recette Receipt) error {
	_,err:=this.Exec("Insert into T_Recette(RC_ID,REC_MONTANT,REC_DATE,REC_Status,UTI_ID) values(?,?,?,?,?)",recette.Id,
															   recette.Total,
															   recette.Date,
															   recette.Statut,
															   recette.Utilisateur_ID)
	if(err!= nil){
		return err
	}
	return nil
}

func (this *OracleDB) DeleteReceipt(recette Receipt) error {
	query := `DELETE FROM T_Rectte
		  WHERE RC_ID = ?`
	_,err:=this.Exec(query, recette.Id)
	if err!= nil {
		return err
	}
	return nil
}


func (this *OracleDB) ModifyReceipt(recette Receipt){
 	_,err:=this.Exec("UPDATE T_Recette SET(REC_MONTANT=?,REC_DATE=?,REC_Status=?,UTI_ID=?) Where(REC_ID=?)",recette.Total,
 																			  recette.Date,
 																			  recette.Statut,
 																			  recette.Utilisateur_ID,
																			  recette.Id)
	if(err!=nil){
		return err
	}
	return nil
}

func (db *OracleDB) FetchOneReceipt(id int) (*Receipt, error) {
	receipt := Receipt{}

	query := `SELECT *
		  FROM T_Recette
		  WHERE id = ?`
	row := db.QueryRow(query, id)
	
	if err := row.Scan(&recette.Id, &recette.Total, &recette.DATE,&recette.Statut , &recette.Utilisateur_ID); err != nil {
		log.Printf("Error while gettiing receipt with id %d", id)
		return nil, err
	}

	return &receipt, nil
}

//User Action
func (this *OracleDB)FetchOneUser() User {
	rows,err := this.Query("Select * From T_Recette") // FIXME: should query user table
	// TODO: check if there is an error

	for _, row := rows.Scan()  { // FIXME: use for rows.Next() to iterate
		// NOTE: rows.Scan() is used here, add the address of the field you want as outputs
	}

}



//TODO: Create Mapping
