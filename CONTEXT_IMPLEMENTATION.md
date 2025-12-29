# Implementación de Context en gopost-api

## Resumen de Cambios

Se ha integrado `context.Context` de Go en toda la aplicación para mejorar el manejo de cancelación, timeouts y propagación de valores.

## 🎯 Beneficios Implementados

### 1. **Cancelación Automática de Operaciones**

- Si un cliente cancela una petición HTTP (cierra el navegador, cancela el request), todas las operaciones de base de datos asociadas se cancelarán automáticamente
- Las queries SQL dejan de ejecutarse cuando el context se cancela
- Libera recursos del servidor y la base de datos

### 2. **Timeouts en el Servidor HTTP**

```go
// server/server.go
srv := &http.Server{
    ReadTimeout:  15 * time.Second,  // Timeout para leer el request
    WriteTimeout: 15 * time.Second,  // Timeout para escribir la respuesta
    IdleTimeout:  60 * time.Second,  // Timeout para conexiones idle
}
```

### 3. **Context en Operaciones de Base de Datos**

Todos los métodos de repositorio ahora usan `*Context` methods:

- `db.ExecContext(ctx, ...)`
- `db.QueryContext(ctx, ...)`
- `db.QueryRowContext(ctx, ...)`

### 4. **Propagación de Valores**

El `server.Context` personalizado ahora incluye `context.Context` nativo y métodos para trabajar con valores.

## 📁 Archivos Modificados

### Core

- [server/context.go](server/context.go) - Agregado `Ctx context.Context` y métodos auxiliares
- [server/router.go](server/router.go) - Pasa `r.Context()` a cada handler
- [server/server.go](server/server.go) - Configuración de timeouts HTTP

### Repositories

- [repositories/user_repository.go](repositories/user_repository.go) - Todos los métodos aceptan `context.Context`
- [repositories/post_repository.go](repositories/post_repository.go) - Todos los métodos aceptan `context.Context`

### Services

- [services/user_service.go](services/user_service.go) - Propagación de context a repositories
- [services/post_service.go](services/post_service.go) - Propagación de context a repositories

### Handlers

- [handlers/user_handler.go](handlers/user_handler.go) - Usa `c.Context()` para pasar context
- [handlers/post_handler.go](handlers/post_handler.go) - Usa `c.Context()` para pasar context

## 💡 Ejemplos de Uso

### Ejemplo 1: Timeout en una Operación Específica

Si quieres agregar un timeout específico a una operación (por ejemplo, 5 segundos para una query lenta):

```go
// En un handler
func (h *PostHandler) GetPostsHandler(c *server.Context) {
    // Crear un context con timeout de 5 segundos
    ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
    defer cancel()

    posts, err := h.postService.GetAllPosts(ctx)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            RespondError(c.RWriter, NewAppError("La operación tardó demasiado", http.StatusRequestTimeout))
            return
        }
        RespondError(c.RWriter, NewAppError(err.Error(), http.StatusInternalServerError))
        return
    }

    RespondJSON(c.RWriter, http.StatusOK, map[string]interface{}{
        "posts": posts,
    })
}
```

### Ejemplo 2: Pasar Valores Adicionales via Context

```go
// En middleware o handler
func SomeMiddleware(next server.HandleFunc) server.HandleFunc {
    return func(c *server.Context) {
        // Agregar request ID
        requestID := uuid.New().String()
        c.WithValue("request_id", requestID)

        next(c)
    }
}

// En otro handler
func (h *UserHandler) SomeHandler(c *server.Context) {
    requestID := c.Value("request_id").(string)
    log.Printf("Request ID: %s", requestID)
}
```

### Ejemplo 3: Detectar Cancelación del Cliente

```go
func (r *PostRepository) FindAll(ctx context.Context) ([]models.Post, error) {
    query := "SELECT id, user_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC"
    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        // Si el cliente canceló la petición
        if errors.Is(err, context.Canceled) {
            return nil, fmt.Errorf("operación cancelada por el cliente")
        }
        return nil, fmt.Errorf("error al obtener posts: %w", err)
    }
    defer rows.Close()

    // ... resto del código
}
```

## 🔥 Buenas Prácticas

1. **Siempre pasar el context**: Nunca usar `context.Background()` en handlers, usa `c.Context()`

2. **No almacenar context**: Los contexts no deben almacenarse en structs, solo pasarse como argumentos

3. **Defer cancel**: Siempre hacer `defer cancel()` cuando creas un context con timeout/cancelación

4. **Verificar errores de context**: Comprueba `context.Canceled` y `context.DeadlineExceeded`

5. **Context debe ser el primer parámetro**: Por convención en Go, siempre es el primer argumento

## 🚀 Mejoras Futuras

- **Logging con context**: Agregar request IDs para trazabilidad
- **Tracing distribuido**: Integrar con OpenTelemetry
- **Metrics**: Usar context para métricas de rendimiento
- **Rate limiting**: Usar context values para limitar requests por usuario
- **Database connection pooling**: Configurar timeouts específicos por operación

## 📊 Antes vs Después

### Antes

```go
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
    err := r.db.QueryRow(query, email).Scan(...)
    // ❌ No se puede cancelar si el cliente desconecta
    // ❌ Sin timeout para queries lentas
}
```

### Después

```go
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    err := r.db.QueryRowContext(ctx, query, email).Scan(...)
    // ✅ Se cancela automáticamente si el cliente desconecta
    // ✅ Respeta timeouts del servidor HTTP (15s)
    // ✅ Permite agregar timeouts personalizados
}
```

## 🧪 Testing

Para tests, puedes usar:

```go
// Test normal
ctx := context.Background()

// Test con timeout
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

// Test con cancelación
ctx, cancel := context.WithCancel(context.Background())
// Simular cancelación
cancel()
```
