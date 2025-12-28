# Resumen - Introducción a net/http

## ¿Qué Aprendiste?

¡Felicidades! Has completado la primera sección del curso. Vamos a repasar todo lo que aprendiste.

## 1. Fundamentos de HTTP y REST

### HTTP (Hypertext Transfer Protocol)

- ✅ Es el protocolo de comunicación de la web
- ✅ Funciona con petición-respuesta entre cliente y servidor
- ✅ Es stateless (sin estado)
- ✅ Usa métodos (GET, POST, PUT, DELETE, etc.)

### REST (Representational State Transfer)

- ✅ Estilo arquitectónico para diseñar APIs
- ✅ Todo es un recurso con una URI única
- ✅ Usa métodos HTTP para operaciones CRUD
- ✅ Usa códigos de estado para comunicar resultados

### Métodos HTTP y CRUD

| Método | Operación | Uso                 |
| ------ | --------- | ------------------- |
| GET    | Read      | Obtener recursos    |
| POST   | Create    | Crear recursos      |
| PUT    | Update    | Actualizar completo |
| PATCH  | Update    | Actualizar parcial  |
| DELETE | Delete    | Eliminar recursos   |

### Códigos de Estado

- **2xx** - Éxito (200 OK, 201 Created, 204 No Content)
- **3xx** - Redirección (301, 302)
- **4xx** - Error del cliente (400, 401, 404)
- **5xx** - Error del servidor (500, 503)

## 2. El Paquete net/http

### Componentes Principales

```go
// 1. Server - El servidor HTTP
server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
}

// 2. ServeMux - El enrutador
mux := http.NewServeMux()

// 3. Handler - Función que procesa peticiones
func handler(w http.ResponseWriter, r *http.Request) {
    // Lógica aquí
}

// 4. Request - Información de la petición
r.Method, r.URL, r.Header, r.Body

// 5. ResponseWriter - Para construir la respuesta
w.Header(), w.WriteHeader(), w.Write()
```

### Ventajas de net/http

✅ Parte de la biblioteca estándar
✅ Alto rendimiento
✅ Concurrencia nativa con goroutines
✅ Producción ready
✅ Sin dependencias externas

## 3. Creación de Servidores

### Servidor Básico

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status":"ok"}`)
    })

    http.ListenAndServe(":8080", mux)
}
```

### Servidor con Configuración

```go
server := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}

server.ListenAndServe()
```

## 4. Manejo de Rutas

### Rutas con Métodos (Go 1.22+)

```go
mux := http.NewServeMux()

// Especificar método directamente
mux.HandleFunc("GET /users", getUsers)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("PUT /users/{id}", updateUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)
```

### Parámetros de Ruta

```go
// Definir parámetro
mux.HandleFunc("GET /users/{id}", getUser)

// Obtener parámetro
func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "Usuario ID: %s", id)
}
```

### Query Parameters

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // /search?q=go&limit=10
    query := r.URL.Query().Get("q")
    limit := r.URL.Query().Get("limit")
}
```

## 5. Request y Response

### Leer la Petición

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Método HTTP
    method := r.Method

    // Ruta
    path := r.URL.Path

    // Cabeceras
    contentType := r.Header.Get("Content-Type")

    // Parámetros de ruta
    id := r.PathValue("id")

    // Query parameters
    page := r.URL.Query().Get("page")

    // Body (JSON)
    var data struct{ Name string }
    json.NewDecoder(r.Body).Decode(&data)
    defer r.Body.Close()
}
```

### Escribir la Respuesta

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Cabeceras
    w.Header().Set("Content-Type", "application/json")

    // 2. Status code
    w.WriteHeader(http.StatusOK)

    // 3. Body
    json.NewEncoder(w).Encode(data)
}
```

## 6. Mini API Práctica

### Estructura

```
gopost-mini/
├── go.mod
├── main.go       # Servidor y rutas
└── handlers.go   # Lógica de handlers
```

### Características Implementadas

✅ CRUD completo
✅ Almacenamiento en memoria con `map[int]Post`
✅ Concurrencia segura con `sync.RWMutex`
✅ Validaciones de datos
✅ Manejo de errores
✅ Respuestas JSON consistentes
✅ Status codes apropiados

### Rutas Implementadas

```
GET    /health        - Health check
GET    /posts         - Listar posts
GET    /posts/{id}    - Obtener post
POST   /posts         - Crear post
PUT    /posts/{id}    - Actualizar post
DELETE /posts/{id}    - Eliminar post
```

## Conceptos Clave

### 1. Concurrencia

```go
var (
    posts   = make(map[int]Post)
    postsMu sync.RWMutex  // ¡Importante para concurrencia!
)

// Lectura
postsMu.RLock()
post := posts[id]
postsMu.RUnlock()

// Escritura
postsMu.Lock()
posts[id] = newPost
postsMu.Unlock()
```

### 2. Validaciones

```go
if post.Title == "" {
    sendError(w, http.StatusBadRequest, "Título requerido")
    return
}
```

### 3. Helpers

```go
func sendJSON(w http.ResponseWriter, code int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(data)
}
```

## Mejores Prácticas Aprendidas

### ✅ Hacer

1. **Especifica métodos HTTP en rutas**: `"GET /users"`
2. **Establece Content-Type**: `w.Header().Set("Content-Type", "application/json")`
3. **Usa códigos de estado correctos**: 200, 201, 400, 404, 500
4. **Cierra el Body**: `defer r.Body.Close()`
5. **Valida datos de entrada**: Antes de procesar
6. **Usa mutex para concurrencia**: `sync.RWMutex`
7. **Crea helpers**: Para respuestas consistentes
8. **Ordena respuesta correctamente**: Headers → Status → Body

### ❌ Evitar

1. No usar verbos en URIs: `/createUser` ❌ → `/users` ✅
2. No mezclar métodos manualmente si puedes especificarlos en la ruta
3. No escribir body antes del status code
4. No olvidar cerrar el body
5. No acceder a variables compartidas sin protección

## Comparación: Antes y Después

### Antes de Esta Sección

```go
// Código básico sin estructura
func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hola")
    })
    http.ListenAndServe(":8080", nil)
}
```

### Después de Esta Sección

```go
// API profesional con:
// - Rutas organizadas con métodos
// - Validaciones
// - Manejo de errores
// - Respuestas JSON estructuradas
// - Concurrencia segura
// - Código modular

mux := http.NewServeMux()
mux.HandleFunc("GET /posts", getPostsHandler)
mux.HandleFunc("POST /posts", createPostHandler)
// ... más rutas

server := &http.Server{
    Addr:    ":8080",
    Handler: mux,
    // ... timeouts
}
```

## ¿Qué Sigue?

En la próxima sección construiremos una API más robusta y profesional con:

### 🚀 Arquitectura en Capas

- **Handlers** - Lógica de presentación
- **Services** - Lógica de negocio
- **Repositories** - Acceso a datos
- **Models** - Estructuras de datos

### 🗄️ Base de Datos

- PostgreSQL con `pgx`
- Migraciones de base de datos
- Conexiones y transacciones
- Consultas optimizadas

### 🔐 Autenticación y Seguridad

- JWT (JSON Web Tokens)
- Middleware de autenticación
- Protección de rutas
- Hashing de contraseñas

### ⚡ Características Avanzadas

- Middleware personalizado
- Logging estructurado
- Manejo de errores centralizado
- Validaciones avanzadas
- Variables de entorno
- Configuración por ambiente

### 🧪 Testing

- Unit tests
- Integration tests
- Mocks y stubs
- Test coverage

### 📦 Deployment

- Docker
- Docker Compose
- Variables de entorno
- Best practices de producción

## Recursos para Profundizar

### Documentación Oficial

- [net/http package](https://pkg.go.dev/net/http)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)
- [Effective Go](https://go.dev/doc/effective_go)

### Lecturas Recomendadas

- [Go by Example: HTTP Servers](https://gobyexample.com/http-servers)
- [REST API Design Best Practices](https://restfulapi.net/)
- [HTTP Status Codes](https://httpstatuses.com/)

## Desafío Final

Antes de continuar, intenta estos desafíos para consolidar lo aprendido:

### Desafío 1: API de Tareas (TODO)

Crea una API completa para gestionar tareas con:

- CRUD de tareas
- Marcar como completada
- Filtrar por estado (completadas/pendientes)
- Buscar por título

### Desafío 2: API de Usuarios y Posts

Extiende la mini API para incluir:

- Usuarios (con CRUD)
- Relación: Un usuario tiene muchos posts
- Rutas: `/users/{id}/posts`
- Validar que el usuario existe al crear un post

### Desafío 3: Middleware de Logging

Implementa un middleware que:

- Registre cada petición (método, ruta, duración)
- Muestre el status code de respuesta
- Use colores para diferentes tipos de respuesta

## Reflexión Final

Has dado un paso importante en tu camino como desarrollador Go. Ahora tienes los fundamentos sólidos para construir APIs REST. Los conceptos que aprendiste aquí son la base de aplicaciones web modernas y escalables.

**¿Te sientes listo para el siguiente nivel?**

En la próxima sección, transformaremos este conocimiento en una API profesional de producción con arquitectura limpia, base de datos y todas las características que esperarías de una aplicación real.

---

## Resumen de Comandos

```bash
# Inicializar proyecto
go mod init github.com/usuario/proyecto

# Ejecutar servidor
go run main.go
go run .

# Probar endpoints
curl http://localhost:8080/health
curl http://localhost:8080/posts
curl -X POST http://localhost:8080/posts -H "Content-Type: application/json" -d '{"title":"Post"}'

# Ver documentación
go doc net/http
go doc net/http.Handler
```

## Checklist de Conocimientos

Marca lo que dominas:

- [ ] Entiendo qué es HTTP y REST
- [ ] Conozco los métodos HTTP y cuándo usarlos
- [ ] Sé usar códigos de estado apropiadamente
- [ ] Puedo crear un servidor HTTP básico
- [ ] Sé definir rutas con métodos específicos
- [ ] Puedo capturar parámetros de ruta
- [ ] Sé trabajar con query parameters
- [ ] Entiendo Request y ResponseWriter
- [ ] Puedo leer y escribir JSON
- [ ] Sé manejar errores apropiadamente
- [ ] Comprendo la concurrencia y uso de mutex
- [ ] Puedo crear una API CRUD completa

Si marcaste todas, ¡estás listo para la siguiente sección! 🎉

---

**Anterior:** [Mini API](07-mini-api.md) | **Siguiente Sección:** Construyendo gopost-api
