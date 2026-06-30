package models

// Estudiante representa un alumno matriculado en la institución.
type Estudiante struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Nombre    string `json:"nombre" gorm:"not null"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Matricula string `json:"matricula" gorm:"not null"`
	Activo    bool   `json:"activo" gorm:"default:true"`
}
