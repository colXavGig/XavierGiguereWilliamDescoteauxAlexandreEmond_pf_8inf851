package DAL

import (
	"time"
)

// Custom types
type StatusReceipt string

// Enums

const ( // UserRole Enums
	UserRoleCommis = "Commis"
	UserRolePDG    = "PDG"
	UserRoleClient = "Client"
)

const ( // StatusReceipt Enums
	StatusReceiptEnAttente StatusReceipt = StatusReceipt("en_attente")
	StatusReceiptValide    StatusReceipt = StatusReceipt("validee")
)

// DTO structs
type User struct {
	Id       int    `json:"ID"`
	Nom      string `json:"nom"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Receipt struct { // NOTE: we can remove db: annotation since we are not using it
	Id             int           `db:"REC_ID" json:"ID"`
	Total          float64       `db:"REC_MONTANT" json:"total"`
	DATE           time.Time     `db:"REC_DATE" json:"date"`
	Statut         StatusReceipt `db:"REC_Status" json:"statut"`
	Utilisateur_ID int           `db:"UTI_ID" json:"utilisateurID"`
}

type Validation struct {
	// TODO: add field
}
