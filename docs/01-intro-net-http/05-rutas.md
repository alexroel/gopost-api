# Rutas en Go con http.ServeMux

## Introducción

En Go 1.22 se introdujeron mejoras significativas en el manejo de rutas con `http.ServeMux`. Ahora podemos especificar métodos HTTP directamente en la definición de rutas y capturar parámetros de forma nativa.

## Rutas Básicas

### Forma Tradicional (Antes de Go 1.22)

```go
mux := http.NewServeMux()
mux.HandleFunc("/health", healthHandler)
mux.HandleFunc("/posts", postsHandler)
```

**Problema:** No distingue entre métodos HTTP. El mismo handler procesa GET, POST, PUT, DELETE.

### Nueva Forma (Go 1.22+)

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /health", healthHandler)
mux.HandleFunc("GET /posts", getPostsHandler)
mux.HandleFunc("POST /posts", createPostHandler)
```

**Ventaja:** El método HTTP es parte de la ruta. Más claro y específico.

## Especificando Métodos HTTP

### Sintaxis

```
"MÉTODO /ruta"
```

### Ejemplos

```go
mux := http.NewServeMux()

// GET - Obtener recursos
mux.HandleFunc("GET /users", getUsers)
mux.HandleFunc("GET /posts", getPosts)

// POST - Crear recursos
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("POST /posts", createPost)

// PUT - Actualizar recursos completos
mux.HandleFunc("PUT /users/{id}", updateUser)
mux.HandleFunc("PUT /posts/{id}", updatePost)

// DELETE - Eliminar recursos
mux.HandleFunc("DELETE /users/{id}", deleteUser)
mux.HandleFunc("DELETE /posts/{id}", deletePost)

// PATCH - Actualizar parcialmente
mux.HandleFunc("PATCH /users/{id}", patchUser)
```

### ¿Qué pasa si usas el método incorrecto?

```bash
# Si defines: "GET /users"
# Y haces:
curl -X POST http://localhost:8080/users

# Resultado: 405 Method Not Allowed
```

Go automáticamente devuelve `405 Method Not Allowed` si el método no coincide. ¡No necesitas verificarlo manualmente!

## Parámetros de Ruta (Path Parameters)

### Capturar Parámetros con {}

```go
// Definir ruta con parámetro
mux.HandleFunc("GET /users/{id}", getUserByID)
mux.HandleFunc("GET /posts/{postId}/comments/{commentId}", getComment)
```

### Obtener Parámetros con PathValue

```go
func getUserByID(w http.ResponseWriter, r *http.Request) {
    // Obtener el parámetro 'id' de la ruta
    id := r.PathValue("id")

    fmt.Fprintf(w, "Usuario ID: %s", id)
}
```

### Ejemplo Completo

```go
package main

import (
    "fmt"
    "net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) {
    // Extraer el parámetro 'id' de la URL
    userID := r.PathValue("id")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id": "%s", "name": "Usuario %s"}`, userID, userID)
}

func getPost(w http.ResponseWriter, r *http.Request) {
    postID := r.PathValue("id")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id": "%s", "title": "Post %s"}`, postID, postID)
}

func main() {
    mux := http.NewServeMux()

    // Rutas con parámetros
    mux.HandleFunc("GET /users/{id}", getUser)
    mux.HandleFunc("GET /posts/{id}", getPost)

    http.ListenAndServe(":8080", mux)
}
```

**Probar:**

```bash
curl http://localhost:8080/users/123
# Resultado: {"id": "123", "name": "Usuario 123"}

curl http://localhost:8080/posts/456
# Resultado: {"id": "456", "title": "Post 456"}
```

## Parámetros Múltiples

```go
func getComment(w http.ResponseWriter, r *http.Request) {
    postID := r.PathValue("postId")
    commentID := r.PathValue("commentId")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"postId": "%s", "commentId": "%s"}`, postID, commentID)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /posts/{postId}/comments/{commentId}", getComment)

    http.ListenAndServe(":8080", mux)
}
```

**Probar:**

```bash
curl http://localhost:8080/posts/10/comments/25
# Resultado: {"postId": "10", "commentId": "25"}
```

## Query Parameters

Los query parameters son diferentes a los path parameters. Van después del `?` en la URL.

### Obtener Query Parameters

```go
func searchUsers(w http.ResponseWriter, r *http.Request) {
    // Obtener query parameters
    name := r.URL.Query().Get("name")
    age := r.URL.Query().Get("age")

    if name == "" {
        name = "sin especificar"
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"buscar": {"name": "%s", "age": "%s"}}`, name, age)
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /users/search", searchUsers)

    http.ListenAndServe(":8080", mux)
}
```

**Probar:**

```bash
curl "http://localhost:8080/users/search?name=Juan&age=30"
# Resultado: {"buscar": {"name": "Juan", "age": "30"}}

curl "http://localhost:8080/users/search"
# Resultado: {"buscar": {"name": "sin especificar", "age": ""}}
```

### Query Parameters Múltiples

```go
func filterPosts(w http.ResponseWriter, r *http.Request) {
    // Obtener todos los valores de un parámetro
    tags := r.URL.Query()["tag"]

    // Obtener un único valor
    author := r.URL.Query().Get("author")
    limit := r.URL.Query().Get("limit")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"author": "%s", "limit": "%s", "tags": %v}`,
                author, limit, tags)
}
```

**Probar:**

```bash
curl "http://localhost:8080/posts?author=Juan&limit=10&tag=go&tag=api"
# Resultado: {"author": "Juan", "limit": "10", "tags": ["go", "api"]}
```

## Diferencia entre Path Parameters y Query Parameters

| Aspecto       | Path Parameters      | Query Parameters              |
| ------------- | -------------------- | ----------------------------- |
| **Sintaxis**  | `/users/{id}`        | `/users?id=123`               |
| **Uso**       | Identificar recursos | Filtros y opciones            |
| **Requerido** | Generalmente sí      | Generalmente opcional         |
| **Ejemplo**   | `/posts/123`         | `/posts?author=Juan&limit=10` |
| **Obtención** | `r.PathValue("id")`  | `r.URL.Query().Get("id")`     |

### Cuándo Usar Cada Uno

**Path Parameters:**

```go
GET /users/123           // ✅ Identificar un usuario específico
GET /posts/456/comments  // ✅ Recursos anidados
DELETE /posts/789        // ✅ Identificar recurso a eliminar
```

**Query Parameters:**

```go
GET /users?role=admin&active=true  // ✅ Filtrar usuarios
GET /posts?page=2&limit=20         // ✅ Paginación
GET /search?q=golang&sort=date     // ✅ Búsqueda
```

## Patrones de Rutas Avanzados

### Wildcard (Catch-all)

```go
// Captura cualquier ruta que empiece con /static/
mux.HandleFunc("GET /static/{path...}", serveStatic)

func serveStatic(w http.ResponseWriter, r *http.Request) {
    // path captura todo después de /static/
    path := r.PathValue("path")
    fmt.Fprintf(w, "Archivo: %s", path)
}
```

**Ejemplos:**

```bash
GET /static/css/style.css     → path = "css/style.css"
GET /static/js/app.js         → path = "js/app.js"
GET /static/images/logo.png   → path = "images/logo.png"
```

### Ruta Raíz

```go
// Coincide SOLO con /
mux.HandleFunc("GET /", homeHandler)

// Coincide con / y cualquier subruta
mux.HandleFunc("GET /{path...}", catchAllHandler)
```

## Ejemplo Completo: API de Posts

```go
package main

import (
    "fmt"
    "net/http"
)

// Handlers
func listPosts(w http.ResponseWriter, r *http.Request) {
    // Query parameters para paginación
    page := r.URL.Query().Get("page")
    limit := r.URL.Query().Get("limit")

    if page == "" {
        page = "1"
    }
    if limit == "" {
        limit = "10"
    }

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"page": "%s", "limit": "%s", "posts": []}`, page, limit)
}

func getPost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id": "%s", "title": "Post %s", "content": "Contenido"}`, id, id)
}

func createPost(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated) // 201
    fmt.Fprintf(w, `{"id": "123", "message": "Post creado"}`)
}

func updatePost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"id": "%s", "message": "Post actualizado"}`, id)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")

    w.WriteHeader(http.StatusNoContent) // 204
    fmt.Fprintf(w, "") // Sin contenido
}

func health(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "ok"}`)
}

func main() {
    mux := http.NewServeMux()

    // Ruta de salud
    mux.HandleFunc("GET /health", health)

    // CRUD de posts
    mux.HandleFunc("GET /posts", listPosts)          // Listar
    mux.HandleFunc("GET /posts/{id}", getPost)       // Obtener uno
    mux.HandleFunc("POST /posts", createPost)        // Crear
    mux.HandleFunc("PUT /posts/{id}", updatePost)    // Actualizar
    mux.HandleFunc("DELETE /posts/{id}", deletePost) // Eliminar

    fmt.Println("🚀 Servidor en http://localhost:8080")
    fmt.Println("📍 Rutas disponibles:")
    fmt.Println("  GET    /health")
    fmt.Println("  GET    /posts")
    fmt.Println("  GET    /posts/{id}")
    fmt.Println("  POST   /posts")
    fmt.Println("  PUT    /posts/{id}")
    fmt.Println("  DELETE /posts/{id}")

    http.ListenAndServe(":8080", mux)
}
```

### Probar la API

```bash
# Health check
curl http://localhost:8080/health

# Listar posts (con paginación)
curl "http://localhost:8080/posts?page=1&limit=5"

# Obtener un post específico
curl http://localhost:8080/posts/123

# Crear un post
curl -X POST http://localhost:8080/posts

# Actualizar un post
curl -X PUT http://localhost:8080/posts/123

# Eliminar un post
curl -X DELETE http://localhost:8080/posts/123
```

## Manejo Manual de Métodos (Fallback)

Si necesitas manejar múltiples métodos en una misma función:

```go
func postsHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        // Lógica para GET
        listPosts(w, r)
    case http.MethodPost:
        // Lógica para POST
        createPost(w, r)
    default:
        // Método no permitido
        w.WriteHeader(http.StatusMethodNotAllowed)
        fmt.Fprintf(w, `{"error": "Método no permitido"}`)
    }
}

mux.HandleFunc("/posts", postsHandler)
```

⚠️ **Nota:** Esto es menos recomendado con Go 1.22+. Mejor usar rutas específicas por método.

## Resumen

✅ En Go 1.22+ puedes especificar métodos HTTP en las rutas: `"GET /users"`
✅ Los parámetros de ruta se definen con `{nombre}` y se obtienen con `r.PathValue("nombre")`
✅ Los query parameters se obtienen con `r.URL.Query().Get("nombre")`
✅ Path parameters para identificar recursos, Query parameters para filtros
✅ Go maneja automáticamente `405 Method Not Allowed`

## Ejercicios

### Ejercicio 1: API de Usuarios

Crea una API con estas rutas:

- `GET /users` - Listar usuarios
- `GET /users/{id}` - Obtener un usuario
- `POST /users` - Crear usuario
- `DELETE /users/{id}` - Eliminar usuario

### Ejercicio 2: Filtros y Búsqueda

Implementa:

- `GET /products?category=electronics&price_max=1000`
- `GET /search?q=laptop&sort=price`

### Ejercicio 3: Rutas Anidadas

Implementa:

- `GET /users/{userId}/posts` - Posts de un usuario
- `GET /users/{userId}/posts/{postId}` - Post específico de un usuario

---

**Anterior:** [Servidor Básico](04-servidor-basico.md) | **Siguiente:** [Request y Response](06-request-response.md)
