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

//Structs proposal
/*
type RentableEntity struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	PricingModel string  `json:"pricing_model"`
	Price        float64 `json:"price"`
	Description  string  `json:"description,omitempty"`
	ImagePath    string  `json:"image_path,omitempty"`
	IsAvailable  bool    `json:"is_available"`
}
type RentalLog struct {
	ID         int    `json:"id"`
	EntityID   int    `json:"entity_id"`
	UserID     int    `json:"user_id"`
	RentalDate string `json:"rental_date"`
	StartTime  string `json:"start_time,omitempty"`
	EndTime    string `json:"end_time,omitempty"`
}
type Receipt struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"` // "pending", "approved", "rejected"
	CreatedAt   string  `json:"created_at"`
	ApprovedAt  string  `json:"approved_at,omitempty"`
	LineItems   []struct {
		EntityID int     `json:"entity_id"`
		Price    float64 `json:"price"`
	} `json:"line_items"`
}
type User struct {
	ID                   int    `json:"id"`
	Email                string `json:"email"`
	Password             string `json:"password,omitempty"` // Should be hashed
	Role                 string `json:"role"`              // "user", "clerk", "admin"
	NotificationPreference bool  `json:"notification_preference"`
}
*/
