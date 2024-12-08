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




//1. Authentication
    //POST /auth/register: Register a new user (role: user by default).
    //POST /auth/login: Authenticate a user and return a session token.
    //POST /auth/logout: Log out the user and invalidate their session.
//2. Users Management
    //GET /users: Retrieve all users (Admin only).
    //POST /users: Create a new user (Admin only).
    //PUT /users/:id: Update a user's details or role (Admin only).
    //DELETE /users/:id: Delete a user (Admin only).
//3. Rentable Entities
    //GET /entities: Retrieve all rentable entities with optional filters (e.g., category, availability).
    //POST /entities: Add a new rentable entity (Admin only).
    //PUT /entities/:id: Update an existing rentable entity (Admin only).
    //DELETE /entities/:id: Delete a rentable entity (Admin only).
//4. Rental Logs
    //GET /rental-logs: Retrieve all rental logs or filter by date, user, or entity (Clerk/Admin).
    //POST /rental-logs: Create a new rental log (Clerk/Admin).
    //DELETE /rental-logs/:id: Delete a rental log (Clerk/Admin).
//5. Receipts
    //GET /receipts: Retrieve all receipts or filter by status, user, or date range (Admin/Clerk).
    //POST /receipts: Create a new receipt (Clerk only).
    //PUT /receipts/:id: Approve or reject a receipt (Admin only).
    //DELETE /receipts/:id: Delete a receipt (Admin only).
//6. Reports
    //GET /reports/revenue: Generate revenue reports based on a time range (Admin/Clerk).
    //GET /reports/availability: Generate availability reports (Admin/Clerk).
