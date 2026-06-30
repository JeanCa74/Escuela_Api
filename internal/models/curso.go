package models

// Curso representa una materia ofertada por la institución.
type Curso struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Nombre   string `json:"nombre" gorm:"not null"`
	Codigo   string `json:"codigo" gorm:"uniqueIndex;not null"`
	Creditos int    `json:"creditos"`
}
