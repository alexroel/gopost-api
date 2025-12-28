# El Paquete net/http

## Introducción

`net/http` es el paquete estándar de Go para construir aplicaciones HTTP. Está incluido en la biblioteca estándar, lo que significa que no necesitas instalar dependencias externas para crear servidores web robustos y eficientes.

## ¿Por Qué net/http?

### Ventajas

✅ **Parte de la biblioteca estándar** - Sin dependencias externas
✅ **Alto rendimiento** - Optimizado y eficiente
✅ **Concurrencia nativa** - Maneja múltiples conexiones simultáneas
✅ **Producción ready** - Usado por grandes empresas
✅ **Bien documentado** - Excelente documentación oficial
✅ **Comunidad activa** - Gran soporte de la comunidad Go

### Casos de Uso

- APIs REST
- Servicios web
- Microservicios
- Servidores HTTP personalizados
- Proxies y gateways
- Webhooks

## Componentes Principales

### 1. http.Server

El servidor HTTP que escucha y maneja peticiones.

```go
server := &http.Server{
    Addr:         ":8080",           // Puerto donde escuchar
    Handler:      nil,                // Manejador de rutas
    ReadTimeout:  10 * time.Second,   // Timeout de lectura
    WriteTimeout: 10 * time.Second,   // Timeout de escritura
    IdleTimeout:  120 * time.Second,  // Timeout de inactividad
}
```

### 2. http.Handler

La interfaz fundamental de net/http. Cualquier tipo que implemente esta interfaz puede manejar peticiones HTTP.

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

### 3. http.HandlerFunc

Una función que puede ser convertida a Handler. Es la forma más común de definir manejadores.

```go
func miManejador(w http.ResponseWriter, r *http.Request) {
    // Lógica del manejador
}
```

### 4. http.ServeMux

El enrutador (router) por defecto de Go. Multiplexor que asigna rutas a manejadores.

```go
mux := http.NewServeMux()
mux.HandleFunc("/ruta", miManejador)
```

### 5. http.Request

Representa una petición HTTP recibida.

```go
type Request struct {
    Method     string              // GET, POST, PUT, etc.
    URL        *url.URL            // URL de la petición
    Header     Header              // Cabeceras HTTP
    Body       io.ReadCloser       // Cuerpo de la petición
    Form       url.Values          // Datos del formulario parseados
    // ... más campos
}
```

**Campos importantes:**

```go
r.Method                    // "GET", "POST", etc.
r.URL.Path                  // "/posts/123"
r.Header.Get("Content-Type") // "application/json"
r.Body                      // io.ReadCloser
r.PathValue("id")           // Obtener parámetros de ruta (Go 1.22+)
```

### 6. http.ResponseWriter

Interfaz para construir la respuesta HTTP.

```go
type ResponseWriter interface {
    Header() Header                // Obtener/modificar cabeceras
    Write([]byte) (int, error)     // Escribir el cuerpo de la respuesta
    WriteHeader(statusCode int)    // Escribir código de estado
}
```

**Uso común:**

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Establecer cabeceras
    w.Header().Set("Content-Type", "application/json")

    // Establecer código de estado
    w.WriteHeader(http.StatusOK) // 200

    // Escribir respuesta
    w.Write([]byte(`{"message": "Hola mundo"}`))
}
```

## Flujo de una Petición HTTP

```
Cliente → Petición HTTP → Servidor → ServeMux → Handler → Respuesta → Cliente
```

**Paso a paso:**

1. **Cliente** envía petición HTTP
2. **Servidor** (http.Server) recibe la conexión
3. **ServeMux** analiza la URL y encuentra el handler apropiado
4. **Handler** procesa la petición y construye la respuesta
5. **Respuesta** se envía de vuelta al cliente

## Ejemplo Básico Completo

```go
package main

import (
    "fmt"
    "net/http"
)

// Handler para la ruta principal
func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "¡Bienvenido a mi API!")
}

// Handler para obtener información
func infoHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"app": "Mi API", "version": "1.0"}`))
}

func main() {
    // Crear el multiplexor (router)
    mux := http.NewServeMux()

    // Registrar rutas y handlers
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/info", infoHandler)

    // Crear y configurar el servidor
    server := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }

    // Iniciar el servidor
    fmt.Println("Servidor escuchando en http://localhost:8080")
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("Error al iniciar servidor: %v\n", err)
    }
}
```

## Métodos Útiles del Paquete

### Para el Servidor

```go
// Iniciar servidor
http.ListenAndServe(":8080", handler)

// Iniciar servidor HTTPS
http.ListenAndServeTLS(":443", "cert.pem", "key.pem", handler)

// Servidor con configuración personalizada
server := &http.Server{/* config */}
server.ListenAndServe()
```

### Para la Request

```go
// Leer body
body, err := io.ReadAll(r.Body)
defer r.Body.Close()

// Parsear JSON
var data struct{ Name string }
json.NewDecoder(r.Body).Decode(&data)

// Obtener query parameters
id := r.URL.Query().Get("id")

// Obtener parámetros de ruta (Go 1.22+)
userID := r.PathValue("id")

// Leer cabeceras
contentType := r.Header.Get("Content-Type")
```

### Para la Response

```go
// Establecer cabeceras
w.Header().Set("Content-Type", "application/json")
w.Header().Add("X-Custom-Header", "valor")

// Enviar código de estado
w.WriteHeader(http.StatusCreated) // 201

// Escribir respuesta
w.Write([]byte("Respuesta"))
fmt.Fprintf(w, "Hola %s", nombre)

// Enviar JSON
json.NewEncoder(w).Encode(data)

// Redireccionar
http.Redirect(w, r, "/nueva-ruta", http.StatusMovedPermanently)
```

## Códigos de Estado Disponibles

Go proporciona constantes para todos los códigos de estado HTTP:

```go
// Éxito
http.StatusOK                  // 200
http.StatusCreated             // 201
http.StatusNoContent           // 204

// Redirección
http.StatusMovedPermanently    // 301
http.StatusFound               // 302

// Errores del cliente
http.StatusBadRequest          // 400
http.StatusUnauthorized        // 401
http.StatusForbidden           // 403
http.StatusNotFound            // 404
http.StatusMethodNotAllowed    // 405
http.StatusConflict            // 409

// Errores del servidor
http.StatusInternalServerError // 500
http.StatusBadGateway          // 502
http.StatusServiceUnavailable  // 503
```

## Helpers de net/http

```go
// Servir archivos estáticos
http.FileServer(http.Dir("./static"))

// Enviar error HTTP
http.Error(w, "Mensaje de error", http.StatusBadRequest)

// Servir un archivo
http.ServeFile(w, r, "archivo.html")

// Redireccionar
http.Redirect(w, r, "/nueva-ruta", http.StatusFound)

// Not Found
http.NotFound(w, r)
```

## Ejemplo: API Simple con net/http

```go
package main

import (
    "encoding/json"
    "net/http"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

func getUsers(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "Juan"},
        {ID: 2, Name: "María"},
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
    var user User

    // Decodificar el JSON del body
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "JSON inválido", http.StatusBadRequest)
        return
    }

    // Simular creación
    user.ID = 3

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /users", getUsers)
    mux.HandleFunc("POST /users", createUser)

    http.ListenAndServe(":8080", mux)
}
```

## Concurrencia en net/http

Una de las características más poderosas de net/http es que **cada petición se maneja en su propia goroutine automáticamente**.

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Esta función se ejecuta en su propia goroutine
    // Go maneja la concurrencia automáticamente

    // Puedes lanzar más goroutines si necesitas
    go procesoEnSegundoPlano()

    w.Write([]byte("Respuesta"))
}
```

⚠️ **Importante:** Debido a la concurrencia, ten cuidado con:

- Variables compartidas (usa mutex o channels)
- Mapas compartidos (no son thread-safe)
- Recursos compartidos

## Resumen

- `net/http` es el paquete estándar de Go para HTTP
- Los componentes principales son: Server, Handler, ServeMux, Request, ResponseWriter
- Es eficiente, concurrente y producción-ready
- Proporciona todo lo necesario para construir APIs REST
- Cada petición se maneja en su propia goroutine automáticamente

## Para Profundizar

📚 **Documentación oficial:**

- [net/http package](https://pkg.go.dev/net/http)
- [Writing Web Applications](https://go.dev/doc/articles/wiki/)

---

**Anterior:** [HTTP y REST](02-http-rest.md) | **Siguiente:** [Servidor Básico](04-servidor-basico.md)
