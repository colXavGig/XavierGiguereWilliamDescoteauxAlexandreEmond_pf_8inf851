const config = {
    // Base URL for the backend API
    apiBaseUrl: "/api", // Replace with the actual backend URL
  
    // Endpoints for the API
    endpoints: {
      // Authentication
      register: "/auth/register",
      login: "/auth/login",
      logout: "/auth/logout",
  
      // Users Management
      users: "/users", // Admin endpoint for managing users (GET, POST, PUT, DELETE)
  
      // Rentable Entities
      rentableEntities: "/entities", // Endpoint for CRUD operations on rentable entities
  
      // Rental Logs
      rentalLogs: "/rental-logs", // Endpoint for logging rentals (GET, POST, DELETE)
  
      // Receipts
      receipts: "/receipts", // Endpoint for receipt management (GET, POST, PUT, DELETE)
  
      // Reports
      reports: {
        revenue: "/reports/revenue", // Endpoint for generating revenue reports
        availability: "/reports/availability", // Endpoint for availability reports
      },
    },
  };
  
  export default config;



/*
  Static Website Routes
	• GET /dashboard/
		○ Serves static files from the ../frontend directory.
		○ No additional requirements for fetch.

API Routes
Receipts
	1. GET /api/receipts
		○ Purpose: Fetch all receipts.
		○ Headers: token (for authentication).
	2. GET /api/receipts/{id}
		○ Purpose: Fetch one receipt by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (receipt ID).
	3. POST /api/receipts/create
		○ Purpose: Create a new receipt.
		○ Headers: token (for authentication).
		○ Body: JSON object representing a receipt (e.g., { "user_id": 1, "status": "Pending", ... }).
	4. DELETE /api/receipts/delete/{id}
		○ Purpose: Delete a receipt by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (receipt ID).

Authentication
	5. POST /api/auth/register
		○ Purpose: Register a new user.
		○ Body: JSON object with user details (e.g., { "email": "user@example.com", "password": "password123" }).
	6. POST /api/auth/login
		○ Purpose: Log in a user.
		○ Body: JSON object with user credentials (e.g., { "email": "user@example.com", "password": "password123" }).

Users
	7. GET /api/users
		○ Purpose: Fetch all users.
		○ Headers: token (for authentication).
	8. POST /api/users/create
		○ Purpose: Create a new user.
		○ Headers: token (for authentication).
		○ Body: JSON object with user details (e.g., { "email": "user@example.com", "password": "password123", ... }).
	9. PUT /api/users/update/{id}
		○ Purpose: Update a user by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (user ID).
		○ Body: JSON object with updated fields (e.g., { "email": "updated@example.com", ... }).
	10. DELETE /api/users/delete/{id}
		○ Purpose: Delete a user by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (user ID).

Entities
	11. GET /api/entities
		○ Purpose: Fetch all rentable entities.
		○ Headers: token (for authentication).
	12. PUT /api/entities/update/{id}
		○ Purpose: Update an entity by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (entity ID).
		○ Body: JSON object with updated fields (e.g., { "name": "Updated Name", "price": 100.00, ... }).

Rental Logs
	13. GET /api/rental-logs
		○ Purpose: Fetch all rental logs.
		○ Headers: token (for authentication).
	14. POST /api/rental-logs/create
		○ Purpose: Create a new rental log.
		○ Headers: token (for authentication).
		○ Body: JSON object with rental log details (e.g., { "user_id": 1, "entity_id": 2, ... }).
	15. DELETE /api/rental-logs/delete/{id}
		○ Purpose: Delete a rental log by ID.
		○ Headers: token (for authentication).
		○ Path Parameter: id (rental log ID).
*/