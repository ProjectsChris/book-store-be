package models

type Book struct {
	Id              int32   `json:"id" bson:"id"`
	Titolo          string  `json:"titolo" bson:"titolo"`
	Autore          string  `json:"autore" bson:"autore"`
	Prezzo          float32 `json:"prezzo" bson:"prezzo"`
	Summary         string  `json:"summary" bson:"summary"`
	CopertinaRigida bool    `json:"copertina_rigida" bson:"copertina_rigida"`
	Genere          string  `json:"genere" bson:"genere"`
	Quantita        int     `json:"quantita" bson:"quantita"`
}
