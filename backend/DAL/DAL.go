package DAL

//---------------
//--           -- cursor parking --           --
//---------------

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/godror/godror"
)
type StatutReceipt string
const  (
	StatusReceiptEnAttente StatutReceipt = StatusReceipt("en_attente")
	StatusReceiptValide StatutReceipt = StatusReceipt("validee")
)


//TODO: Create DAL
type OracleDB struct {
	*sql.DB
}

type Receipt struct{
	Id int `db:"REC_ID"`
	Total float `db:"REC_MONTANT"`
	DATE time.Time `db:"REC_DATE"`
	Statut StatutReceipt  `db:"REC_Status"`
	Utilisateur_ID int `db:"UTI_ID"`
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


//Receipt Action
func (this *OracleDB)FetchAllReceipt() ([]Receipt,error) {
	var recette_list []Receipt
	
	rows,err := this.Query("Select * From T_Recette")
	if err != nil {
		return nil, err
	}

	recette_list = make([]Receipt,5)

	for rows.Next()  {
		recette := Receipt{}

		err := rows.Scan(&recette.Id, &recette.Total, &recette.DATE,&recette.Statut , &recette.Utilisateur_ID)
		if err != nil {
			return nil, err
		}
		recette_list := append(recette_list, recette)
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

func (this *OracleDB)ModifyReceipt(recette Receipt){
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
//User Action
func (this *OracleDB)FetchOneUser() []Receipt {
	rows,err := this.Query("Select * From T_Recette")

	for _, row := rows.Scan()  {
		
	}
	
	

}



//TODO: Create Mapping
