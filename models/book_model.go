package models

type Book struct {
	Id          int     `json:"id"`
	Titolo      string  `json:"titolo" validate:"required,max=255" example:"Il silenzio di un mare in tempesta"`
	Autore      string  `json:"autore" validate:"required,max=64" example:"Ruben Fabrizi"`
	Prezzo      float32 `json:"prezzo" validate:"required" example:"15.90"`
	Summary     string  `json:"summary" validate:"required,max=512" example:"Lorem Ipsum is simply dummy text of the printing and typesetting industry."`
	Copertina   string  `json:"copertina" validate:"required,eq=Hard Cover|eq=Flexible Cover" example:"Hard Cover"`
	Genere      string  `json:"genere" validate:"required,eq=Action|eq=Adventure|eq=Business|eq=Cookbooks|eq=Drama|eq=Detective|eq=Fantasy|eq=Fiction|eq=History|eq=Horror|eq=Romance|eq=Psychology|eq=Science Fiction|eq=Short Stories|eq=Thriller" example:"Romance"`
	Quantita    int     `json:"quantita" validate:"required,gte=1,lte=5" example:"2"`
	Categoria   string  `json:"categoria"  validate:"required,eq=Best Seller|eq=New Releases|eq=Best Offers" example:"New Releases"`
	IdCopertina int     `json:"id_copertina" validate:"required" example:"564"`
}
