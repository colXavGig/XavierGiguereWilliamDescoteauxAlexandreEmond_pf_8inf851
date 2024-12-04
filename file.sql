Create table T_Utilisateur(
UTI_ID INT PRIMARY KEY,
UTI_NOM VARCHAR(100),
UTI_ROLE ENUM('Commis','PDG','Client'),
UTI_EMAIL VARCHAR(150),
UTI_MDP VARCHAR(255),
)

Create table T_Recette(
REC_ID INT PRIMARY KEY,
REC_Montant float,
REC_DATE DATE,
REC_Statut ENUM('en_attente','validee','rejetee'),
UTI_ID INT
Constraint PK_UTI_ID FOREIGN KEY UTI_ID REFERENCES T_Utilisateur(UTI_ID)
)

Create table T_Validation(
VAL_ID INT PRIMARY KEY,
REC_ID float,
REC_DATE DATE,
VAL_Statut ENUM('en_attente','validee','rejetee'),
VAL_Commentaire TEXT
Constraint PK_REC_ID FOREIGN KEY REC_ID REFERENCES T_Recette(REC_ID)
)

Create table T_Chambre(
CHA_ID INT PRIMARY KEY,
CHA_NOM VARCHAR(75),
CHA_PRIX float,
)

Create table T_Magasin(
MAG_ID INT PRIMARY KEY,
MAG_NOM VARCHAR(75),
MAG_PRIX float,
)

Create table T_Restaurant(
RES_ID INT PRIMARY KEY,
RES_NOM VARCHAR(75),
RES_PRIX float,
)

Create table T_WC(
WC_ID INT PRIMARY KEY,
WC_NOM VARCHAR(75),
)

Create table T_Douche(
DOU_ID INT PRIMARY KEY,
DOU_NOM VARCHAR(75),
)