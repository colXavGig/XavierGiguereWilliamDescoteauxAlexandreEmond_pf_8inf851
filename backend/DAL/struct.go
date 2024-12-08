package DAL

// Custom types
type StatusReceipt string

// Enums

const ( // UserRole Enums
	UserRoleCommis = "clerk"
	UserRolePDG    = "admin"
	UserRoleClient = "user"
)

const ( // StatusReceipt Enums
	StatusReceiptEnAttente = "pending"
	StatusReceiptValide    = "approved"
	StatusReceiptRejetee   = "rejected"
)

// DTO structs

type Validation struct {
	// TODO: add field
}

//Structs proposal

type RentableEntity struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	PricingModel string  `json:"pricing_model"`
	Price        float64 `json:"price"`
	Description  string  `json:"description,omitempty"`
	ImagePath    string  `json:"image_path,omitempty"`
	IsAvailable  bool    `json:"is_available"` // NOTE: valeur calculer
}
type RentalLog struct {
	ID         int    `json:"id"`
	EntityID   int    `json:"entity_id"`
	UserID     int    `json:"user_id"`
	RentalDate string `json:"rental_date"`
	StartTime  string `json:"start_time,omitempty"` //-----> pas mieux d'avoir
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
	ID                     int    `json:"id,omitempty"`
	Email                  string `json:"email"`
	Password               string `json:"password,omitempty"` // Should be hashed
	Role                   string `json:"role,omitempty"`     // "user", "clerk", "admin"
	NotificationPreference int    `json:"notification_preference,omitempty"`
}
