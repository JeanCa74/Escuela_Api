# escuela-api — Módulo de Gestión Escolar

**Aplicaciones Web II (TDI-601) · Semana 11 · Actividad C1**

Proyecto grupal. Dominio: **Escuela** (Estudiantes, Cursos e Inscripciones).  
Cada integrante implementó su propio módulo completo: model → repository → service → handler.

---

## Integrantes y módulos

| Integrante | Módulo | Entidad | Regla de negocio |
|------------|--------|---------|-----------------|
| **Jean Carlos** | Estudiantes | `Estudiante` | Nombre vacío → `ErrNombreVacio` (400) |
| **Jhon** | Cursos | `Curso` | Créditos ≤ 0 → `ErrCreditosInvalidos` (400) |
| **Maria José** | Inscripciones | `Inscripcion` | Calificación fuera de [0,10] → `ErrCalificacionInvalida` (400) |

---

## Arquitectura: model → repository → service → handler

```
internal/
├── models/
│   ├── estudiante.go      — entidad Estudiante (GORM tags)
│   ├── curso.go           — entidad Curso
│   ├── inscripcion.go     — entidad Inscripcion
│   └── usuario.go         — entidad Usuario (auth)
├── storage/
│   ├── repositorio.go     — interfaces: EstudianteRepository, CursoRepository,
│   │                         InscripcionRepository, Almacen, UserRepository
│   ├── memoria.go         — fake en RAM (usado en tests de handler)
│   ├── gorm.go            — implementación GORM + SQLite (producción y test :memory:)
│   ├── gorm_test.go       ✅ Jean Carlos — repositorio Estudiante contra SQLite :memory:
│   ├── gorm_curso_test.go ✅ Jhon        — repositorio Curso contra SQLite :memory:
│   └── gorm_inscripcion_test.go ✅ Maria José — repositorio Inscripcion contra SQLite :memory:
├── service/
│   ├── errores.go         — errores de dominio compartidos
│   ├── estudiante.go      — reglas de negocio: validarEstudiante
│   ├── estudiante_test.go ✅ Jean Carlos — service con mock
│   ├── curso.go           — reglas de negocio: validarCurso
│   ├── curso_test.go      ✅ Jhon        — service con mock
│   ├── inscripcion.go     — reglas de negocio: validarInscripcion
│   ├── inscripcion_test.go ✅ Maria José — service con mock
│   └── auth.go            — bcrypt + JWT
├── handlers/
│   ├── server.go          — Server con los tres servicios inyectados
│   ├── respond.go         — RespondJSON / RespondError / statusDeError
│   ├── estudiante.go      — CRUD de estudiantes
│   ├── estudiante_test.go ✅ Jean Carlos — handler con httptest + 401
│   ├── curso.go           — CRUD de cursos
│   ├── curso_test.go      ✅ Jhon        — handler con httptest + 401
│   ├── inscripcion.go     — CRUD de inscripciones
│   ├── inscripcion_test.go ✅ Maria José — handler con httptest + 401
│   └── auth.go            — register + login
└── middleware/
    ├── auth.go            — JWT middleware (Bearer token) → produce el 401
    └── cors.go            — CORS para desarrollo
```

---

## Reglas de negocio por módulo

### Jean Carlos — Estudiantes

| Regla | Error de dominio | HTTP |
|-------|-----------------|------|
| `Nombre` vacío o solo espacios | `ErrNombreVacio` | 400 |
| `Email` vacío | `ErrEmailVacio` | 400 |
| Recurso no encontrado | `ErrNoEncontrado` | 404 |

### Jhon — Cursos

| Regla | Error de dominio | HTTP |
|-------|-----------------|------|
| `Nombre` vacío o solo espacios | `ErrNombreVacio` | 400 |
| `Creditos` igual a cero o negativo | `ErrCreditosInvalidos` | 400 |
| Recurso no encontrado | `ErrNoEncontrado` | 404 |

### Maria José — Inscripciones

| Regla | Error de dominio | HTTP |
|-------|-----------------|------|
| `Calificacion` fuera del rango [0, 10] | `ErrCalificacionInvalida` | 400 |
| Recurso no encontrado | `ErrNoEncontrado` | 404 |

### Auth (compartido)

| Regla | Error de dominio | HTTP |
|-------|-----------------|------|
| Token ausente o inválido | — | 401 (middleware) |
| Email ya registrado | `ErrEmailEnUso` | 409 |
| Credenciales incorrectas | `ErrCredencialesInvalidas` | 401 |

---

## Tests — Suite en verde (18 tests)

```bash
cd escuela-api
go test ./... -cover
```

Salida esperada:
```
ok  escuela-api/internal/handlers   coverage: 35.5%
ok  escuela-api/internal/service    coverage: 35.4%
ok  escuela-api/internal/storage    coverage: 23.5%
```

---

## Módulo Jean Carlos — Estudiantes

### Test 1 — Service con mock (`service/estudiante_test.go`)

**Qué comprueba:**  
Que `validarEstudiante` rechaza datos inválidos **antes** de llamar al repositorio.

**Cómo funciona:**  
- Se crea un `estudianteRepoMock` (testify/mock) que registra llamadas.
- Se llama a `svc.Crear(...)` con nombre vacío → se espera `ErrNombreVacio`.
- Se verifica con `AssertNotCalled` que el repo **nunca** recibió `CrearEstudiante`.

**Qué se rompería:**  
Si se elimina la validación de `validarEstudiante`, el mock recibirá una llamada no esperada y el test falla con "unexpected call to CrearEstudiante".

### Test 2 — Handler con httptest (`handlers/estudiante_test.go`)

**Test 2a — `TestCrearEstudiante_Exitoso`:**  
POST `/api/v1/estudiantes` con token válido y cuerpo correcto → 201 Created.

**Test 2b — `TestRutaProtegida_SinToken` (el 401):**  
POST `/api/v1/estudiantes` **sin** header `Authorization` → 401 Unauthorized.

**Cómo funciona:**  
- `construirEntorno` arma el router completo con middleware Auth **real**.
- `registrarYObtenerToken` hace register + login para obtener un JWT real.
- El test del 401 envía la petición sin el header; el middleware corta antes de llegar al handler.

**Qué se rompería:**  
Si se elimina `r.Use(middleware.Auth(authSvc))`, la petición llegaría al handler y respondería 201 en vez de 401.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_test.go`)

**Qué comprueba:**  
Que `CrearEstudiante` persiste en SQLite y `BuscarEstudiantePorID` lo refleja.

**Cómo funciona:**  
- `NuevoAlmacenGORM(":memory:")` abre una DB que vive solo mientras dure el test.
- Se crea un estudiante → se busca por ID → se verifica nombre y email.

**Qué se rompería:**  
Si `AutoMigrate` no crea la tabla, el `db.Create` falla y el ID queda en 0; el test falla en `require.NotZero(t, creado.ID)`.

---

## Módulo Jhon — Cursos

### Test 1 — Service con mock (`service/curso_test.go`)

**Qué comprueba:**  
Que `validarCurso` rechaza créditos inválidos **antes** de llamar al repositorio.

**Cómo funciona:**  
- Se crea un `cursoRepoMock` (testify/mock) que registra llamadas.
- Se llama a `svc.Crear(...)` con `Creditos: 0` → se espera `ErrCreditosInvalidos`.
- Se verifica con `AssertNotCalled` que el repo **nunca** recibió `CrearCurso`.

**Qué se rompería:**  
Si se elimina la validación de créditos en `validarCurso`, el mock recibirá una llamada inesperada y el test falla.

### Test 2 — Handler con httptest (`handlers/curso_test.go`)

**Test 2a — `TestCrearCurso_Exitoso`:**  
POST `/api/v1/cursos` con token y créditos válidos → 201 Created.

**Test 2b — `TestRutaCursos_SinToken` (el 401):**  
POST `/api/v1/cursos` **sin** header `Authorization` → 401 Unauthorized.

**Qué se rompería:**  
Si se elimina `r.Use(middleware.Auth(authSvc))` del grupo protegido, la petición llegaría al handler y respondería 201 en lugar de 401.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_curso_test.go`)

**Qué comprueba:**  
Que `CrearCurso` persiste en SQLite y `BuscarCursoPorID` lo refleja.

**Cómo funciona:**  
- `NuevoAlmacenGORM(":memory:")` abre una DB temporal.
- Se crea un curso → se busca por ID → se verifica nombre y créditos.

**Qué se rompería:**  
Si `AutoMigrate` no crea la tabla `cursos`, el Create falla y el ID queda en 0.

---

## Módulo Maria José — Inscripciones

### Test 1 — Service con mock (`service/inscripcion_test.go`)

**Qué comprueba:**  
Que `validarInscripcion` rechaza calificaciones fuera de [0,10] **antes** de llamar al repositorio.

**Cómo funciona:**  
- Se crea un `inscripcionRepoMock` (testify/mock) que registra llamadas.
- Se llama a `svc.Crear(...)` con `Calificacion: 15` → se espera `ErrCalificacionInvalida`.
- Se verifica con `AssertNotCalled` que el repo **nunca** recibió `CrearInscripcion`.

**Qué se rompería:**  
Si se elimina la validación del rango, el mock recibe una llamada inesperada y el test falla.

### Test 2 — Handler con httptest (`handlers/inscripcion_test.go`)

**Test 2a — `TestCrearInscripcion_Exitosa`:**  
POST `/api/v1/inscripciones` con token y calificación válida (8.5) → 201 Created.

**Test 2b — `TestRutaInscripciones_SinToken` (el 401):**  
POST `/api/v1/inscripciones` **sin** header `Authorization` → 401 Unauthorized.

**Qué se rompería:**  
Si se elimina `r.Use(middleware.Auth(authSvc))`, la petición llega al handler y responde 201 en lugar de 401.

### Test 3 — Repositorio GORM :memory: (`storage/gorm_inscripcion_test.go`)

**Qué comprueba:**  
Que `CrearInscripcion` persiste en SQLite y `BuscarInscripcionPorID` lo refleja.

**Cómo funciona:**  
- `NuevoAlmacenGORM(":memory:")` abre una DB temporal.
- Se crea una inscripción → se busca por ID → se verifica EstudianteID, CursoID y Calificacion.

**Qué se rompería:**  
Si `AutoMigrate` no crea la tabla `inscripciones`, el Create falla y el ID queda en 0.

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

| Método | Ruta | Auth | Módulo |
|--------|------|------|--------|
| POST | `/api/v1/auth/register` | No | — |
| POST | `/api/v1/auth/login` | No | — |
| GET | `/api/v1/estudiantes` | Bearer JWT | Jean Carlos |
| POST | `/api/v1/estudiantes` | Bearer JWT | Jean Carlos |
| GET | `/api/v1/estudiantes/{id}` | Bearer JWT | Jean Carlos |
| PUT | `/api/v1/estudiantes/{id}` | Bearer JWT | Jean Carlos |
| DELETE | `/api/v1/estudiantes/{id}` | Bearer JWT | Jean Carlos |
| GET | `/api/v1/cursos` | Bearer JWT | Jhon |
| POST | `/api/v1/cursos` | Bearer JWT | Jhon |
| GET | `/api/v1/cursos/{id}` | Bearer JWT | Jhon |
| GET | `/api/v1/inscripciones` | Bearer JWT | Maria José |
| POST | `/api/v1/inscripciones` | Bearer JWT | Maria José |
| GET | `/api/v1/inscripciones/{id}` | Bearer JWT | Maria José |
