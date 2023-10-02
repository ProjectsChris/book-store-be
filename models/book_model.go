package models

type Book struct {
	Id          int     `json:"id"`
	Titolo      string  `json:"titolo" validate:"required,max=255"`
	Autore      string  `json:"autore" validate:"required"`
	Prezzo      float32 `json:"prezzo" validate:"required"`
	Summary     string  `json:"summary" validate:"required,max=512"`
	Copertina   string  `json:"copertina" validate:"required,eq=Hard Cover|eq=Flexible Cover"`
	Genere      string  `json:"genere" validate:"required,eq=Novel|eq=Fantasy|eq=Business|eq=Psychology|eq=Design|eq=Fiction"`
	Quantita    int     `json:"quantita" validate:"required,gte=1,lte=5"`
	Categoria   string  `json:"categoria"  validate:"required,eq=Best Seller|eq=New Releases|eq=Best Offers"`
	IdCopertina string  `json:"id_copertina" validate:"required"`
}
