# escuela-api — Módulo de Gestión Escolar

**Aplicaciones Web II (TDI-601) · Semana 11 · Actividad C1**

Módulo propio del proyecto grupal. Dominio: **Escuela** (Estudiantes y Cursos).  
No es una copia de `cafeteria-uleam-api`; adapta los mismos patrones a un dominio diferente.

---

## Arquitectura: model → repository → service → handler

```
internal/
├── models/
│   ├── estudiante.go    — entidad Estudiante (GORM tags)
│   ├── curso.go         — entidad Curso
│   └── usuario.go       — entidad Usuario (auth)
├── storage/
│   ├── repositorio.go   — interfaces EstudianteRepository, CursoRepository, Almacen, UserRepository
│   ├── memoria.go       — fake en RAM (usado en tests de handler)
│   ├── gorm.go          — implementación GORM + SQLite (producción y test :memory:)
│   └── gorm_test.go     ✅ TEST — repositorio real contra SQLite :memory:
├── service/
│   ├── errores.go       — errores de dominio
│   ├── estudiante.go    — reglas de negocio (validarEstudiante)
│   ├── auth.go          — bcrypt + JWT
│   └── estudiante_test.go  ✅ TEST — service con mock
├── handlers/
│   ├── server.go        — Server con servicios inyectados
│   ├── respond.go       — RespondJSON / RespondError / statusDeError
│   ├── estudiante.go    — CRUD de estudiantes
│   ├── auth.go          — register + login
│   └── estudiante_test.go  ✅ TEST — handler con httptest + 401
└── middleware/
    ├── auth.go          — JWT middleware (Bearer token)
    └── cors.go          — CORS para desarrollo
```

---

## Reglas de negocio (dominio)

| Regla | Error de dominio | HTTP |
|-------|-----------------|------|
| `Nombre` vacío o solo espacios | `ErrNombreVacio` | 400 |
| `Email` vacío | `ErrEmailVacio` | 400 |
| Recurso no encontrado | `ErrNoEncontrado` | 404 |
| Email ya registrado | `ErrEmailEnUso` | 409 |
| Token ausente o inválido | `ErrCredencialesInvalidas` | 401 |

---

## Tests — Suite en verde

```bash
cd escuela-api
go test ./... -cover
```

Salida esperada:
```
ok  escuela-api/internal/handlers   coverage: 34.7%
ok  escuela-api/internal/service    coverage: 22.6%
ok  escuela-api/internal/storage    coverage: 12.9%
```

### Test 1 — Service con mock (`service/estudiante_test.go`)

**Qué comprueba:**  
Que `validarEstudiante` rechaza datos inválidos **antes** de llamar al repositorio.

**Cómo funciona:**  
- Se crea un `estudianteRepoMock` (testify/mock) que registra llamadas.
- Se llama a `svc.Crear(...)` con nombre vacío → se espera `ErrNombreVacio`.
- Se verifica con `AssertNotCalled` que el repo **nunca** recibió `CrearEstudiante`.

**Qué se rompería si la implementación fallara:**  
Si se elimina la validación de `validarEstudiante`, el mock recibirá una llamada no esperada y el test falla con "unexpected call to CrearEstudiante".

### Test 2 — Handler con httptest (`handlers/estudiante_test.go`)

**Test 2a — `TestCrearEstudiante_Exitoso`:**  
POST `/api/v1/estudiantes` con token válido y cuerpo correcto → 201 Created.

**Test 2b — `TestRutaProtegida_SinToken` (el 401):**  
POST `/api/v1/estudiantes` **sin** header `Authorization` → 401 Unauthorized.

**Cómo funciona:**  
- `construirEntorno` arma el router completo con middleware Auth **real**.
- `registrarYObtenerToken` hace register + login contra el propio router para obtener un JWT real.
- El test del 401 envía la petición **sin** el header; el middleware corta antes de llegar al handler.

**Qué se rompería:**  
Si se elimina `r.Use(middleware.Auth(authSvc))` del grupo protegido, la petición llegaría al handler y respondería 201 en vez de 401.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_test.go`)

**Qué comprueba:**  
Que `CrearEstudiante` persiste en SQLite y `BuscarEstudiantePorID` lo refleja.

**Cómo funciona:**  
- `NuevoAlmacenGORM(":memory:")` abre una DB que vive solo mientras dure el test.
- Se crea un estudiante → se busca por ID → se verifica nombre y email.

**Qué se rompería:**  
Si `AutoMigrate` no crea la tabla, el `db.Create` falla y el ID queda en 0; el test falla en `require.NotZero(t, creado.ID)`.

---

## Conceptos clave (glosario para la presentación)

| Término | Definición |
|---------|-----------|
| **Mock** | Doble de prueba que verifica que ciertas llamadas ocurrieron (o no). Usa testify/mock. |
| **Fake** | Implementación alternativa funcional (guarda datos en RAM). No verifica llamadas. |
| **httptest** | Paquete estándar de Go para probar handlers HTTP sin levantar un servidor real. |
| **:memory:** | DSN especial de SQLite que abre la base en RAM; se destruye al cerrar la conexión. |
| **401** | Lo genera el middleware Auth cuando el header Authorization está ausente o el token es inválido. El handler nunca lo produce. |
| **AutoMigrate** | GORM crea/actualiza la tabla automáticamente a partir del struct Go. |

---

## Ejecutar el servidor

```bash
go run ./cmd/escuela-api
# → http://localhost:8080
```

### Endpoints disponibles

| Método | Ruta | Auth | Descripción |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register` | No | Registrar usuario |
| POST | `/api/v1/auth/login` | No | Obtener JWT |
| GET | `/api/v1/estudiantes` | Bearer JWT | Listar estudiantes |
| POST | `/api/v1/estudiantes` | Bearer JWT | Crear estudiante |
| GET | `/api/v1/estudiantes/{id}` | Bearer JWT | Obtener por ID |
| PUT | `/api/v1/estudiantes/{id}` | Bearer JWT | Actualizar |
| DELETE | `/api/v1/estudiantes/{id}` | Bearer JWT | Eliminar |
