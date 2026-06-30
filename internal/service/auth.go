package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"escuela-api/internal/models"
	"escuela-api/internal/storage"
)

// secretoJWT firma y verifica los tokens.
var secretoJWT = []byte("escuela-api-secreto-demo")

const duracionToken = 24 * time.Hour

// Claims es el contenido del JWT.
type Claims struct {
	UsuarioID int `json:"uid"`
	jwt.RegisteredClaims
}

// AuthService concentra toda la lógica de autenticación: bcrypt + JWT.
type AuthService struct {
	repo storage.UserRepository
}

func NuevoAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || strings.TrimSpace(password) == "" {
		return models.Usuario{}, ErrCredencialesInvalidas
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}
	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
	})
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}
	return s.generarToken(u)
}

func (s *AuthService) generarToken(u models.Usuario) (string, error) {
	claims := Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretoJWT)
}

func (s *AuthService) ValidarToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretoJWT, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrCredencialesInvalidas
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidas
	}
	return claims.UsuarioID, nil
}
