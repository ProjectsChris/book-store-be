package models

type Book struct {
	Titolo          string  `json:"titolo" bson:"titolo" validate:"required"`
	Autore          string  `json:"autore" bson:"autore" validate:"required"`
	Prezzo          float32 `json:"prezzo" bson:"prezzo" validate:"required"`
	Summary         string  `json:"summary" bson:"summary" validate:"required"`
	CopertinaRigida bool    `json:"copertina_rigida" bson:"copertina_rigida" validate:"required"`
	Genere          string  `json:"genere" bson:"genere" validate:"required"`
	Quantita        int     `json:"quantita" bson:"quantita" validate:"required"`
	IdFoto          string  `json:"id_foto" bson:"id_foto" validate:"required"`
}
