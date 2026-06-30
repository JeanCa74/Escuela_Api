package service

import "errors"

// Errores de dominio. El handler los traduce a códigos HTTP:
//
//	ErrNombreVacio, ErrEmailVacio, ErrCreditosInvalidos -> 400 Bad Request
//	ErrNoEncontrado                                     -> 404 Not Found
//	ErrEmailEnUso                                       -> 409 Conflict
//	ErrCredencialesInvalidas                            -> 401 Unauthorized
var (
	ErrNombreVacio           = errors.New("el campo nombre es obligatorio")
	ErrEmailVacio            = errors.New("el campo email es obligatorio")
	ErrCreditosInvalidos     = errors.New("los creditos deben ser mayores a cero")
	ErrNoEncontrado          = errors.New("recurso no encontrado")
	ErrEmailEnUso            = errors.New("el email ya esta registrado")
	ErrCredencialesInvalidas = errors.New("email o contrasena incorrectos")
)
