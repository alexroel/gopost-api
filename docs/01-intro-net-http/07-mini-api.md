# Mini API - Proyecto Práctico

## Objetivo

Construir una mini API funcional que maneje posts con operaciones CRUD, separando el código en archivos organizados y manteniendo los datos en memoria.

## Estructura del Proyecto

```
gopost-mini/
├── go.mod
├── main.go          # Punto de entrada, configuración del servidor
└── handlers.go      # Lógica de los handlers
```

## Paso 1: Inicializar el Proyecto

```bash
mkdir gopost-mini
cd gopost-mini
go mod init github.com/tuusuario/gopost-mini
```

## Paso 2: Crear el Archivo handlers.go

Este archivo contendrá toda la lógica de los handlers y el almacenamiento en memoria.

```go
// handlers.go
package main

import (
    "encoding/json"
    "net/http"
    "strconv"
    "sync"
)

// Modelo de datos
type Post struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

// Almacenamiento en memoria
var (
    posts   = make(map[int]Post)
    nextID  = 1
    postsMu sync.RWMutex // Mutex para concurrencia segura
)

// Respuestas estándar
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}

type SuccessResponse struct {
    Success bool `json:"success"`
    Data    any  `json:"data,omitempty"`
}

// Helper para enviar JSON
func sendJSON(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}

// Helper para enviar errores
func sendError(w http.ResponseWriter, statusCode int, message string) {
    sendJSON(w, statusCode, ErrorResponse{
        Error:   http.StatusText(statusCode),
        Message: message,
    })
}

// GET /health - Health check
func healthHandler(w http.ResponseWriter, r *http.Request) {
    sendJSON(w, http.StatusOK, map[string]string{
        "status": "ok",
    })
}

// GET /posts - Obtener todos los posts
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
    postsMu.RLock()
    defer postsMu.RUnlock()

    // Convertir el map a slice
    postList := make([]Post, 0, len(posts))
    for _, post := range posts {
        postList = append(postList, post)
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    postList,
    })
}

// GET /posts/{id} - Obtener un post específico
func getPostHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        sendError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    postsMu.RLock()
    post, exists := posts[id]
    postsMu.RUnlock()

    if !exists {
        sendError(w, http.StatusNotFound, "Post no encontrado")
        return
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    post,
    })
}

// POST /posts - Crear un nuevo post
func createPostHandler(w http.ResponseWriter, r *http.Request) {
    var post Post

    // Decodificar JSON
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        sendError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    defer r.Body.Close()

    // Validaciones
    if post.Title == "" {
        sendError(w, http.StatusBadRequest, "El título es requerido")
        return
    }

    if post.Content == "" {
        sendError(w, http.StatusBadRequest, "El contenido es requerido")
        return
    }

    // Asignar ID y guardar
    postsMu.Lock()
    post.ID = nextID
    nextID++
    posts[post.ID] = post
    postsMu.Unlock()

    sendJSON(w, http.StatusCreated, SuccessResponse{
        Success: true,
        Data:    post,
    })
}

// PUT /posts/{id} - Actualizar un post
func updatePostHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        sendError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    var updatedPost Post
    if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
        sendError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    defer r.Body.Close()

    // Validaciones
    if updatedPost.Title == "" {
        sendError(w, http.StatusBadRequest, "El título es requerido")
        return
    }

    if updatedPost.Content == "" {
        sendError(w, http.StatusBadRequest, "El contenido es requerido")
        return
    }

    postsMu.Lock()
    defer postsMu.Unlock()

    // Verificar si existe
    if _, exists := posts[id]; !exists {
        sendError(w, http.StatusNotFound, "Post no encontrado")
        return
    }

    // Actualizar manteniendo el ID
    updatedPost.ID = id
    posts[id] = updatedPost

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    updatedPost,
    })
}

// DELETE /posts/{id} - Eliminar un post
func deletePostHandler(w http.ResponseWriter, r *http.Request) {
    idStr := r.PathValue("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        sendError(w, http.StatusBadRequest, "ID inválido")
        return
    }

    postsMu.Lock()
    defer postsMu.Unlock()

    if _, exists := posts[id]; !exists {
        sendError(w, http.StatusNotFound, "Post no encontrado")
        return
    }

    delete(posts, id)

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    map[string]string{"message": "Post eliminado"},
    })
}
```

## Paso 3: Crear el Archivo main.go

Este archivo configura el servidor y registra las rutas.

```go
// main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    // Crear el router
    mux := http.NewServeMux()

    // Registrar rutas
    mux.HandleFunc("GET /health", healthHandler)
    mux.HandleFunc("GET /posts", getPostsHandler)
    mux.HandleFunc("GET /posts/{id}", getPostHandler)
    mux.HandleFunc("POST /posts", createPostHandler)
    mux.HandleFunc("PUT /posts/{id}", updatePostHandler)
    mux.HandleFunc("DELETE /posts/{id}", deletePostHandler)

    // Configurar el servidor
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Mensaje de inicio
    fmt.Println("🚀 Servidor iniciando en http://localhost:8080")
    fmt.Println("\n📍 Rutas disponibles:")
    fmt.Println("  GET    /health          - Health check")
    fmt.Println("  GET    /posts           - Listar todos los posts")
    fmt.Println("  GET    /posts/{id}      - Obtener un post")
    fmt.Println("  POST   /posts           - Crear un post")
    fmt.Println("  PUT    /posts/{id}      - Actualizar un post")
    fmt.Println("  DELETE /posts/{id}      - Eliminar un post")
    fmt.Println("\n✨ Presiona Ctrl+C para detener el servidor")

    // Iniciar el servidor
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("❌ Error al iniciar servidor: %v", err)
    }
}
```

## Paso 4: Ejecutar la API

```bash
go run .
```

**Salida:**

```
🚀 Servidor iniciando en http://localhost:8080

📍 Rutas disponibles:
  GET    /health          - Health check
  GET    /posts           - Listar todos los posts
  GET    /posts/{id}      - Obtener un post
  POST   /posts           - Crear un post
  PUT    /posts/{id}      - Actualizar un post
  DELETE /posts/{id}      - Eliminar un post

✨ Presiona Ctrl+C para detener el servidor
```

## Paso 5: Probar la API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Respuesta:**

```json
{
  "status": "ok"
}
```

### 2. Listar Posts (Vacío Inicialmente)

```bash
curl http://localhost:8080/posts
```

**Respuesta:**

```json
{
  "success": true,
  "data": []
}
```

### 3. Crear un Post

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aprendiendo Go",
    "content": "Go es un lenguaje increíble para construir APIs"
  }'
```

**Respuesta:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Aprendiendo Go",
    "content": "Go es un lenguaje increíble para construir APIs"
  }
}
```

### 4. Crear Más Posts

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "REST APIs",
    "content": "Construyendo APIs REST con net/http"
  }'

curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Concurrencia en Go",
    "content": "Goroutines y channels son poderosos"
  }'
```

### 5. Listar Todos los Posts

```bash
curl http://localhost:8080/posts
```

**Respuesta:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Aprendiendo Go",
      "content": "Go es un lenguaje increíble para construir APIs"
    },
    {
      "id": 2,
      "title": "REST APIs",
      "content": "Construyendo APIs REST con net/http"
    },
    {
      "id": 3,
      "title": "Concurrencia en Go",
      "content": "Goroutines y channels son poderosos"
    }
  ]
}
```

### 6. Obtener un Post Específico

```bash
curl http://localhost:8080/posts/1
```

**Respuesta:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Aprendiendo Go",
    "content": "Go es un lenguaje increíble para construir APIs"
  }
}
```

### 7. Actualizar un Post

```bash
curl -X PUT http://localhost:8080/posts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Aprendiendo Go - Actualizado",
    "content": "Go es el mejor lenguaje para construir APIs REST"
  }'
```

**Respuesta:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Aprendiendo Go - Actualizado",
    "content": "Go es el mejor lenguaje para construir APIs REST"
  }
}
```

### 8. Eliminar un Post

```bash
curl -X DELETE http://localhost:8080/posts/2
```

**Respuesta:**

```json
{
  "success": true,
  "data": {
    "message": "Post eliminado"
  }
}
```

### 9. Probar Errores

**Post no encontrado:**

```bash
curl http://localhost:8080/posts/999
```

**Respuesta:**

```json
{
  "error": "Not Found",
  "message": "Post no encontrado"
}
```

**JSON inválido:**

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d 'JSON INVÁLIDO'
```

**Respuesta:**

```json
{
  "error": "Bad Request",
  "message": "JSON inválido"
}
```

**Falta el título:**

```bash
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"content": "Solo contenido"}'
```

**Respuesta:**

```json
{
  "error": "Bad Request",
  "message": "El título es requerido"
}
```

## Explicación del Código

### Almacenamiento en Memoria

```go
var (
    posts   = make(map[int]Post)  // Map para almacenar posts
    nextID  = 1                    // Contador para IDs
    postsMu sync.RWMutex           // Mutex para seguridad en concurrencia
)
```

**¿Por qué un Mutex?**

El servidor maneja cada petición en una goroutine separada. Si dos peticiones modifican `posts` simultáneamente, puede haber condiciones de carrera (race conditions). El mutex garantiza acceso seguro:

- `RLock()` / `RUnlock()` - Para lectura (múltiples goroutines pueden leer)
- `Lock()` / `Unlock()` - Para escritura (solo una goroutine puede escribir)

### Funciones Helper

```go
func sendJSON(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
```

Estas funciones evitan repetir código y garantizan respuestas consistentes en toda la API.

### Validaciones

Cada handler valida los datos antes de procesarlos:

```go
if post.Title == "" {
    sendError(w, http.StatusBadRequest, "El título es requerido")
    return
}
```

Esto previene datos inconsistentes en el almacenamiento.

## Características Implementadas

✅ **CRUD completo** - Crear, Leer, Actualizar, Eliminar
✅ **Validaciones** - Datos requeridos, IDs válidos
✅ **Manejo de errores** - Respuestas consistentes
✅ **Concurrencia segura** - Uso de mutex
✅ **Código organizado** - Separación de responsabilidades
✅ **Status codes correctos** - 200, 201, 400, 404
✅ **Respuestas JSON** - Formato estándar

## Mejoras Posibles

### 1. Añadir Timestamps

```go
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}
```

### 2. Añadir Paginación

```go
func getPostsHandler(w http.ResponseWriter, r *http.Request) {
    page := r.URL.Query().Get("page")
    limit := r.URL.Query().Get("limit")

    // Implementar lógica de paginación
}
```

### 3. Añadir Búsqueda

```go
mux.HandleFunc("GET /posts/search", searchPostsHandler)
```

### 4. Añadir Logging

```go
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}
```

## Ejercicios

### Ejercicio 1: Añadir Campos

Añade los campos `author` y `tags` al modelo Post.

### Ejercicio 2: Filtrado

Implementa `GET /posts?author=nombre` para filtrar posts por autor.

### Ejercicio 3: Validación Avanzada

Añade validación para que el título tenga mínimo 5 caracteres y máximo 100.

### Ejercicio 4: PATCH

Implementa `PATCH /posts/{id}` para actualización parcial (sin requerir todos los campos).

## Resumen

✅ Creaste una API REST funcional con Go
✅ Implementaste CRUD completo con almacenamiento en memoria
✅ Organizaste el código en múltiples archivos
✅ Manejaste concurrencia de forma segura
✅ Validaste datos y manejaste errores apropiadamente
✅ Usaste status codes HTTP correctos

¡Has construido una API real y funcional! En la siguiente sección construiremos una API más robusta con base de datos, arquitectura en capas y características avanzadas.

---

**Anterior:** [Request y Response](06-request-response.md) | **Siguiente:** [Resumen](08-resumen.md)
