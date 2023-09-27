package models

type Book struct {
	Titolo          string  `json:"titolo"`
	Autore          string  `json:"autore"`
	Prezzo          float32 `json:"prezzo"`
	Summary         string  `json:"summary"`
	CopertinaRigida bool    `json:"copertina_rigida"`
	Genere          string  `json:"genere"`
	Quantita        int     `json:"quantita"`
}
