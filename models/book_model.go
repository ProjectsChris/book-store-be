package models

type Book struct {
	Id          int     `json:"id"`
	Titolo      string  `json:"titolo" validate:"required,max=255"` // TODO: fix length
	Autore      string  `json:"autore" validate:"required,max=64"`
	Prezzo      float32 `json:"prezzo" validate:"required"`
	Summary     string  `json:"summary" validate:"required,max=512"`
	Copertina   string  `json:"copertina" validate:"required,eq=Hard Cover|eq=Flexible Cover"`
	Genere      string  `json:"genere" validate:"required,eq=Action|eq=Adventure|eq=Business|eq=Cookbooks|eq=Drama|eq=Detective|eq=Fantasy|eq=Fiction|eq=History|eq=Horror|eq=Romance|eq=Psychology|eq=Science Fiction|eq=Short Stories|eq=Thriller|"`
	Quantita    int     `json:"quantita" validate:"required,gte=1,lte=5"`
	Categoria   string  `json:"categoria"  validate:"required,eq=Best Seller|eq=New Releases|eq=Best Offers"`
	IdCopertina int     `json:"id_copertina" validate:"required"`
}
