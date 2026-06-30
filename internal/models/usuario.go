package models

import "time"

// Usuario representa una cuenta que puede autenticarse en la API.
// PasswordHash nunca se serializa a JSON (tag "-").
type Usuario struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	CreadoEn     time.Time `json:"creado_en"`
}
