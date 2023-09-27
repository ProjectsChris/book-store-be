package models

type Book struct {
	Titolo      string  `json:"titolo" bson:"titolo" validate:"required,max=255"`
	Autore      string  `json:"autore" bson:"autore" validate:"required"`
	Prezzo      float32 `json:"prezzo" bson:"prezzo" validate:"required"`
	Summary     string  `json:"summary" bson:"summary" validate:"required,max=512"`
	Copertina   string  `json:"copertina" bson:"copertina" validate:"required,eq=Hard Cover|eq=Flexible Cover"`
	Genere      string  `json:"genere" bson:"genere" validate:"required,eq=Novel|eq=Fantasy|eq=Business|eq=Psychology|eq=Design|eq=Fiction"`
	Quantita    int     `json:"quantita" bson:"quantita" validate:"required,gte=1,lte=5"`
	Categoria   string  `json:"categoria"  bson:"categoria" validate:"required,eq=Best Seller|eq=New Releases|eq=Best Offers"`
	IdCopertina string  `json:"id_copertina" bson:"id_copertina" validate:"required"`
}
