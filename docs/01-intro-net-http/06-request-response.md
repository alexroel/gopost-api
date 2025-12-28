# Request y Response

## La Petición (Request)

La estructura `http.Request` contiene toda la información sobre la petición HTTP que el cliente envía al servidor.

### Estructura de http.Request

```go
type Request struct {
    Method     string              // GET, POST, PUT, DELETE, etc.
    URL        *url.URL            // URL parseada
    Proto      string              // "HTTP/1.1"
    Header     Header              // Cabeceras HTTP
    Body       io.ReadCloser       // Cuerpo de la petición
    Host       string              // Host del servidor
    RemoteAddr string              // IP del cliente
    RequestURI string              // URI sin parsear
    // ... más campos
}
```

### Campos Importantes de Request

#### 1. Method - Método HTTP

```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Method) // "GET", "POST", "PUT", etc.

    // Verificar método manualmente (si no usas rutas con método)
    if r.Method != http.MethodPost {
        http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
        return
    }
}
```

#### 2. URL - Información de la URL

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Ruta completa
    fmt.Println(r.URL.Path)     // "/api/users"

    // Query parameters
    fmt.Println(r.URL.RawQuery) // "page=1&limit=10"

    // Obtener query parameter específico
    page := r.URL.Query().Get("page")

    // Obtener todos los valores de un parámetro
    tags := r.URL.Query()["tag"]
}
```

#### 3. Header - Cabeceras HTTP

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Obtener una cabecera
    contentType := r.Header.Get("Content-Type")
    authorization := r.Header.Get("Authorization")
    userAgent := r.Header.Get("User-Agent")

    // Verificar si existe una cabecera
    if r.Header.Get("X-API-Key") == "" {
        http.Error(w, "API Key requerida", http.StatusUnauthorized)
        return
    }

    // Obtener todas las cabeceras
    for name, values := range r.Header {
        fmt.Printf("%s: %v\n", name, values)
    }
}
```

#### 4. Body - Cuerpo de la Petición

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Leer el body completo
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error leyendo body", http.StatusBadRequest)
        return
    }
    defer r.Body.Close() // ¡Importante! Siempre cerrar el body

    fmt.Println(string(body))
}
```

#### 5. PathValue - Parámetros de Ruta (Go 1.22+)

```go
// Ruta: GET /users/{id}
func handler(w http.ResponseWriter, r *http.Request) {
    userID := r.PathValue("id")
    fmt.Printf("Usuario ID: %s\n", userID)
}
```

#### 6. RemoteAddr - IP del Cliente

```go
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.RemoteAddr) // "127.0.0.1:54321"
}
```

### Leer JSON del Body

```go
import (
    "encoding/json"
    "net/http"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    // Decodificar JSON del body
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Usar los datos
    fmt.Printf("Usuario: %s (%s)\n", user.Name, user.Email)

    // Responder
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

**Probar:**

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Juan","email":"juan@example.com"}'
```

### Validar Datos de la Petición

```go
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Validaciones
    if user.Name == "" {
        http.Error(w, "El nombre es requerido", http.StatusBadRequest)
        return
    }

    if user.Email == "" {
        http.Error(w, "El email es requerido", http.StatusBadRequest)
        return
    }

    // Si todo está bien, procesar...
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
```

## La Respuesta (Response)

La interfaz `http.ResponseWriter` se usa para construir y enviar la respuesta al cliente.

### Interfaz ResponseWriter

```go
type ResponseWriter interface {
    Header() Header                // Obtener/modificar cabeceras
    Write([]byte) (int, error)     // Escribir el cuerpo
    WriteHeader(statusCode int)    // Establecer código de estado
}
```

### Orden de Operaciones

⚠️ **Importante:** El orden correcto es:

1. **Establecer cabeceras** con `w.Header().Set()`
2. **Escribir código de estado** con `w.WriteHeader()`
3. **Escribir cuerpo** con `w.Write()` o `json.NewEncoder()`

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // 1. Cabeceras primero
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Custom-Header", "valor")

    // 2. Código de estado
    w.WriteHeader(http.StatusOK)

    // 3. Cuerpo
    w.Write([]byte(`{"message": "ok"}`))
}
```

❌ **Incorrecto:**

```go
// Mal: escribir body antes del status code
w.Write([]byte("algo"))
w.WriteHeader(http.StatusOK) // Ya no tiene efecto
```

### 1. Establecer Cabeceras

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Establecer Content-Type
    w.Header().Set("Content-Type", "application/json")

    // Cabeceras personalizadas
    w.Header().Set("X-API-Version", "1.0")
    w.Header().Set("X-Request-ID", "abc123")

    // Cabeceras de caché
    w.Header().Set("Cache-Control", "no-cache")

    // CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
}
```

### 2. Códigos de Estado

```go
// Éxito
w.WriteHeader(http.StatusOK)                  // 200
w.WriteHeader(http.StatusCreated)             // 201
w.WriteHeader(http.StatusNoContent)           // 204

// Redirección
w.WriteHeader(http.StatusMovedPermanently)    // 301
w.WriteHeader(http.StatusFound)               // 302

// Errores del cliente
w.WriteHeader(http.StatusBadRequest)          // 400
w.WriteHeader(http.StatusUnauthorized)        // 401
w.WriteHeader(http.StatusForbidden)           // 403
w.WriteHeader(http.StatusNotFound)            // 404

// Errores del servidor
w.WriteHeader(http.StatusInternalServerError) // 500
```

### 3. Escribir el Cuerpo

#### Texto Plano

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Hola mundo"))

    // O usando fmt.Fprintf
    fmt.Fprintf(w, "Hola %s", "mundo")
}
```

#### JSON

```go
func handler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{
        "message": "Éxito",
        "status":  "ok",
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(data)
}
```

#### JSON con Structs

```go
type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Success: true,
        Message: "Operación exitosa",
        Data: map[string]int{
            "count": 42,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

### Helpers para Respuestas Comunes

#### 1. http.Error - Respuestas de Error

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Forma simple de enviar un error
    http.Error(w, "Recurso no encontrado", http.StatusNotFound)

    // Equivalente a:
    // w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    // w.WriteHeader(http.StatusNotFound)
    // w.Write([]byte("Recurso no encontrado\n"))
}
```

#### 2. http.Redirect - Redirecciones

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Redirigir a otra URL
    http.Redirect(w, r, "/nueva-ruta", http.StatusFound)

    // Redirección permanente
    http.Redirect(w, r, "/nueva-ruta", http.StatusMovedPermanently)
}
```

#### 3. http.NotFound - 404

```go
func handler(w http.ResponseWriter, r *http.Request) {
    http.NotFound(w, r)

    // Equivalente a:
    // http.Error(w, "404 page not found", http.StatusNotFound)
}
```

### Ejemplo Completo: API con Respuestas Bien Estructuradas

```go
package main

import (
    "encoding/json"
    "net/http"
)

// Estructuras de datos
type Post struct {
    ID      int    `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

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

// Handlers
func getPosts(w http.ResponseWriter, r *http.Request) {
    posts := []Post{
        {ID: 1, Title: "Primer post", Content: "Contenido 1"},
        {ID: 2, Title: "Segundo post", Content: "Contenido 2"},
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    posts,
    })
}

func getPost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    // Simular que no existe
    if id == "999" {
        sendError(w, http.StatusNotFound, "Post no encontrado")
        return
    }

    post := Post{
        ID:      1,
        Title:   "Mi Post",
        Content: "Contenido del post",
    }

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    post,
    })
}

func createPost(w http.ResponseWriter, r *http.Request) {
    var post Post

    // Leer y validar JSON
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

    // Simular creación
    post.ID = 123

    sendJSON(w, http.StatusCreated, SuccessResponse{
        Success: true,
        Data:    post,
    })
}

func updatePost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    var post Post
    if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
        sendError(w, http.StatusBadRequest, "JSON inválido")
        return
    }
    defer r.Body.Close()

    sendJSON(w, http.StatusOK, SuccessResponse{
        Success: true,
        Data:    map[string]string{"id": id, "message": "Post actualizado"},
    })
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    // 204 No Content - sin cuerpo en la respuesta
    w.WriteHeader(http.StatusNoContent)

    // Alternativamente, con confirmación:
    // sendJSON(w, http.StatusOK, SuccessResponse{
    //     Success: true,
    //     Data:    map[string]string{"id": id, "deleted": "true"},
    // })
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /posts", getPosts)
    mux.HandleFunc("GET /posts/{id}", getPost)
    mux.HandleFunc("POST /posts", createPost)
    mux.HandleFunc("PUT /posts/{id}", updatePost)
    mux.HandleFunc("DELETE /posts/{id}", deletePost)

    http.ListenAndServe(":8080", mux)
}
```

### Probar la API

```bash
# Obtener todos los posts
curl http://localhost:8080/posts

# Obtener un post
curl http://localhost:8080/posts/1

# Post no encontrado
curl http://localhost:8080/posts/999

# Crear un post
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"Nuevo Post","content":"Contenido"}'

# Crear con error (sin título)
curl -X POST http://localhost:8080/posts \
  -H "Content-Type: application/json" \
  -d '{"content":"Solo contenido"}'

# Actualizar
curl -X PUT http://localhost:8080/posts/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Actualizado","content":"Nuevo contenido"}'

# Eliminar
curl -X DELETE http://localhost:8080/posts/1 -v
```

## Mejores Prácticas

### 1. Siempre Establece Content-Type

```go
// JSON
w.Header().Set("Content-Type", "application/json")

// HTML
w.Header().Set("Content-Type", "text/html; charset=utf-8")

// Texto plano
w.Header().Set("Content-Type", "text/plain; charset=utf-8")
```

### 2. Usa Códigos de Estado Apropiados

```go
// Crear recurso
w.WriteHeader(http.StatusCreated) // 201

// Actualización exitosa
w.WriteHeader(http.StatusOK) // 200

// Eliminación exitosa sin contenido
w.WriteHeader(http.StatusNoContent) // 204

// Datos inválidos
w.WriteHeader(http.StatusBadRequest) // 400

// No encontrado
w.WriteHeader(http.StatusNotFound) // 404
```

### 3. Cierra el Body

```go
func handler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close() // ¡Siempre!

    // ... procesar body
}
```

### 4. Maneja Errores Consistentemente

```go
// Crea funciones helper para respuestas consistentes
func sendError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{
        "error": message,
    })
}
```

## Resumen

✅ `http.Request` contiene toda la información de la petición
✅ Accede a datos con: `r.Method`, `r.URL`, `r.Header`, `r.Body`, `r.PathValue()`
✅ `http.ResponseWriter` se usa para construir respuestas
✅ Orden correcto: Cabeceras → Status Code → Body
✅ Usa helpers como `sendJSON()` para respuestas consistentes
✅ Siempre cierra el `Body` con `defer r.Body.Close()`

---

**Anterior:** [Rutas](05-rutas.md) | **Siguiente:** [Mini API](07-mini-api.md)
