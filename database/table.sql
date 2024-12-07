-- Table 4: Users
CREATE TABLE Users (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    email VARCHAR2(100) UNIQUE NOT NULL,
    password VARCHAR2(255) NOT NULL,
    role VARCHAR2(20) CHECK (role IN ('user', 'clerk', 'admin')) NOT NULL,
    notification_preference NUMBER(1) DEFAULT 1 CHECK (notification_preference IN (0, 1))
);

-- Insert sample data for users
INSERT INTO Users (email, password, role, notification_preference)
VALUES 
('admin@test.com', '123', 'admin', 0),
('clerk@test.com', '123', 'clerk', 0),
('user@test.com', '123', 'user', 1);

-- Table 1: Rentable Entities
CREATE TABLE Rentable_Entities (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    name VARCHAR2(100) NOT NULL,
    category VARCHAR2(50) CHECK (category IN ('Hotel', 'Magasin', 'WC', 'Douche')) NOT NULL,
    pricing_model VARCHAR2(20) CHECK (pricing_model IN ('per_day', 'per_month', 'per_use')) NOT NULL,
    price NUMBER(10, 2) NOT NULL,
    description VARCHAR2(500), -- Optional for display purposes
    image_path VARCHAR2(255)   -- Optional for UI presentation
);

-- Table 2: Rental Logs
CREATE TABLE Rental_Logs (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    entity_id NUMBER NOT NULL,
    user_id NUMBER NOT NULL,
    rental_date DATE NOT NULL,
    start_time TIMESTAMP, -- Optional for per_use rentals
    end_time TIMESTAMP,   -- Optional for per_use rentals
    CONSTRAINT fk_rental_entity FOREIGN KEY (entity_id) REFERENCES Rentable_Entities(id),
    CONSTRAINT fk_rental_user FOREIGN KEY (user_id) REFERENCES Users(id),
    CONSTRAINT unique_daily_rentals UNIQUE (entity_id, rental_date), 
    CONSTRAINT unique_per_use_rentals UNIQUE (entity_id, rental_date, start_time)
);



-- Table 3: Receipts
CREATE TABLE Receipts (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    user_id NUMBER NOT NULL,
    total_amount NUMBER(10, 2) NOT NULL,
    status VARCHAR2(20) CHECK (status IN ('pending', 'approved', 'rejected')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_receipt_user FOREIGN KEY (user_id) REFERENCES Users(id)
);

-- Table 3.1: Receipt Line Items Subtable
CREATE TABLE Receipt_Line_Items (
    id NUMBER GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    receipt_id NUMBER NOT NULL,
    entity_id NUMBER NOT NULL,
    price NUMBER(10, 2) NOT NULL,
    CONSTRAINT fk_line_item_receipt FOREIGN KEY (receipt_id) REFERENCES Receipts(id),
    CONSTRAINT fk_line_item_entity FOREIGN KEY (entity_id) REFERENCES Rentable_Entities(id)
);



-- Insert statements for the hotel rooms
INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Chambre Ile Prince Edouard', 'Hotel', 'per_day', 50.00, 'Room with modern amenities', 'images/room1.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Chambre Québec City', 'Hotel', 'per_day', 60.00, 'Room with modern amenities', 'images/room2.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Chambre Ottawa', 'Hotel', 'per_day', 80.00, 'Room with modern amenities', 'images/room3.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Chambre St John', 'Hotel', 'per_day', 70.00, 'Room with modern amenities', 'images/room4.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Chambre New Brunswick', 'Hotel', 'per_day', 100.00, 'Room with modern amenities', 'images/room5.jpg');

-- Insert statements for the stores
INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin Ile Prince Edouard', 'Magasin', 'per_month', 100.00, 'Retail space with high visibility', 'images/magasin1.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin Québec City', 'Magasin', 'per_month', 90.00, 'Retail space with high visibility', 'images/magasin2.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin Chambre Ottawa', 'Magasin', 'per_month', 80.00, 'Retail space with high visibility', 'images/magasin3.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin Chambre St John', 'Magasin', 'per_month', 70.00, 'Retail space with high visibility', 'images/magasin4.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin New Brunswick', 'Magasin', 'per_month', 100.00, 'Retail space with high visibility', 'images/magasin5.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Magasin Halifax', 'Magasin', 'per_month', 100.00, 'Retail space with high visibility', 'images/magasin6.jpg');

-- Insert statements for the restaurant spaces
INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Espace Ile Prince Edouard', 'Hotel', 'per_month', 100.00, 'Spacious area for restaurant setup', 'images/restaurant1.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Espace Québec City', 'Hotel', 'per_month', 90.00, 'Spacious area for restaurant setup', 'images/restaurant2.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Espace Chambre Ottawa', 'Hotel', 'per_month', 80.00, 'Spacious area for restaurant setup', 'images/restaurant3.jpg');


-- Insert statements for the public WCs
INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('WC 1', 'WC', 'per_use', 2.00, 'Clean public restroom', 'images/wc1.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('WC 2', 'WC', 'per_use', 2.00, 'Clean public restroom', 'images/wc2.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('WC 3', 'WC', 'per_use', 2.00, 'Clean public restroom', 'images/wc3.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('WC 4', 'WC', 'per_use', 2.00, 'Clean public restroom', 'images/wc4.jpg');

-- Insert statements for the public showers
INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Shower 1', 'Douche', 'per_use', 5.00, 'Clean public shower', 'images/shower1.jpg');

INSERT INTO Rentable_Entities (name, category, pricing_model, price, description, image_path)
VALUES ('Shower 2', 'Douche', 'per_use', 5.00, 'Clean public shower', 'images/shower2.jpg');






                                    --       View         --
/*
Create view V_Chambre(
CHA_ID INT PRIMARY KEY,
CHA_NOM VARCHAR(75),
CHA_PRIX float,
Constraint FK_CHA_ID FOREIGN KEY (CHA_ID) REFERENCES Rentable_Entities(id)
)

Create View V_Magasin(
MAG_ID INT PRIMARY KEY,
MAG_NOM VARCHAR(75),
MAG_PRIX float,
Constraint FK_MAG_ID FOREIGN KEY (MAG_ID) REFERENCES Rentable_Entities(id)
)

Create View V_Restaurant(
RES_ID INT PRIMARY KEY,
RES_NOM VARCHAR(75),
RES_PRIX float,
Constraint FK_RES_ID FOREIGN KEY (RES_ID) REFERENCES Rentable_Entities(id)
)

Create View V_WC(
WC_ID INT PRIMARY KEY,
WC_NOM VARCHAR(75),
Constraint FK_WC_ID FOREIGN KEY (WC_ID) REFERENCES Rentable_Entities(id)
)

Create View V_Douche(
DOU_ID INT PRIMARY KEY,
DOU_NOM VARCHAR(75),
)
*/