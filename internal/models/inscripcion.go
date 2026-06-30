package models

// Inscripcion registra la matrícula de un Estudiante en un Curso.
type Inscripcion struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	EstudianteID int     `json:"estudiante_id" gorm:"not null"`
	CursoID      int     `json:"curso_id" gorm:"not null"`
	Calificacion float64 `json:"calificacion"`
}
