# Servidor Básico con Go

## Creando el Proyecto gopost-api

Vamos a crear nuestro proyecto desde cero. Sigue estos pasos:

### Paso 1: Crear el Directorio del Proyecto

```bash
# Crear directorio
mkdir gopost-api
cd gopost-api
```

### Paso 2: Inicializar el Módulo Go

```bash
go mod init github.com/tuusuario/gopost-api
```

Esto creará un archivo `go.mod` que gestiona las dependencias del proyecto:

```go
module github.com/tuusuario/gopost-api

go 1.21
```

### Paso 3: Crear el Archivo main.go

Crea el archivo `main.go` en la raíz del proyecto.

## Nuestro Primer Servidor HTTP

Vamos a crear un servidor HTTP básico paso a paso:

### Versión 1: El Servidor Más Simple

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    fmt.Println("Servidor iniciando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**Ejecutar:**

```bash
go run main.go
```

**¿Qué hace este código?**

1. `package main` - Define el paquete principal ejecutable
2. `import` - Importa los paquetes necesarios
3. `http.ListenAndServe(":8080", nil)` - Inicia el servidor en el puerto 8080

⚠️ **Problema:** Este servidor no hace nada útil. Responde con 404 a todas las peticiones porque no tiene rutas definidas (handler es `nil`).

**Probar:**

```bash
curl http://localhost:8080
# Resultado: 404 page not found
```

### Versión 2: Servidor con un Handler Simple

```go
package main

import (
    "fmt"
    "net/http"
)

// Handler que responde a todas las peticiones
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "¡Hola desde Go!")
}

func main() {
    // Registrar el handler para todas las rutas
    http.HandleFunc("/", helloHandler)

    fmt.Println("Servidor iniciando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**Ejecutar y probar:**

```bash
go run main.go

# En otra terminal:
curl http://localhost:8080
# Resultado: ¡Hola desde Go!
```

**¿Qué hace este código?**

1. `helloHandler` - Función que procesa peticiones HTTP
   - `w http.ResponseWriter` - Para escribir la respuesta
   - `r *http.Request` - Información de la petición
2. `http.HandleFunc("/", helloHandler)` - Registra el handler para la ruta "/"
3. `fmt.Fprintf(w, ...)` - Escribe la respuesta al cliente

### Versión 3: Servidor con Información de la Petición

```go
package main

import (
    "fmt"
    "net/http"
)

func requestInfoHandler(w http.ResponseWriter, r *http.Request) {
    // Establecer tipo de contenido
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

    // Información de la petición
    fmt.Fprintf(w, "Método: %s\n", r.Method)
    fmt.Fprintf(w, "URL: %s\n", r.URL.Path)
    fmt.Fprintf(w, "Host: %s\n", r.Host)
    fmt.Fprintf(w, "User-Agent: %s\n", r.Header.Get("User-Agent"))
}

func main() {
    http.HandleFunc("/", requestInfoHandler)

    fmt.Println("Servidor iniciando en http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

**Probar:**

```bash
curl http://localhost:8080
```

**Salida:**

```
Método: GET
URL: /
Host: localhost:8080
User-Agent: curl/7.81.0
```

**¿Qué aprendemos?**

- `r.Method` - El método HTTP usado (GET, POST, etc.)
- `r.URL.Path` - La ruta de la URL
- `r.Host` - El host de la petición
- `r.Header.Get()` - Obtener cabeceras HTTP
- `w.Header().Set()` - Establecer cabeceras de respuesta

### Versión 4: Servidor con Manejo de Errores

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprintf(w, "Servidor funcionando correctamente\n")
}

func main() {
    // Crear un logger personalizado
    logger := log.New(os.Stdout, "SERVER: ", log.LstdFlags)

    // Registrar handler
    http.HandleFunc("/", mainHandler)

    // Información de inicio
    port := ":8080"
    logger.Printf("Iniciando servidor en http://localhost%s", port)

    // Iniciar servidor con manejo de errores
    err := http.ListenAndServe(port, nil)
    if err != nil {
        logger.Fatalf("Error al iniciar servidor: %v", err)
    }
}
```

**¿Qué mejoras tiene?**

1. **Logger personalizado** - Mensajes con formato
2. **Manejo de errores** - Captura errores al iniciar el servidor
3. **Puerto como variable** - Más fácil de modificar

### Versión 5: Servidor con Configuración Personalizada

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`, time.Now().Format(time.RFC3339))
}

func main() {
    // Configurar el multiplexor
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthHandler)

    // Configurar el servidor con timeouts
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,  // Timeout para leer la petición
        WriteTimeout: 10 * time.Second,  // Timeout para escribir la respuesta
        IdleTimeout:  120 * time.Second, // Timeout para conexiones idle
    }

    log.Println("🚀 Servidor iniciando en http://localhost:8080")
    log.Println("📍 Endpoint disponible: GET /health")

    // Iniciar servidor
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("❌ Error al iniciar servidor: %v", err)
    }
}
```

**Explicación detallada:**

#### 1. El Handler healthHandler

```go
func healthHandler(w http.ResponseWriter, r *http.Request) {
    // Establecer el Content-Type a JSON
    w.Header().Set("Content-Type", "application/json")

    // Establecer el código de estado HTTP 200 OK
    w.WriteHeader(http.StatusOK)

    // Enviar respuesta JSON
    fmt.Fprintf(w, `{"status": "ok", "timestamp": "%s"}`,
                time.Now().Format(time.RFC3339))
}
```

- **Línea 2:** Establece que la respuesta será JSON
- **Línea 5:** Envía el código de estado 200 (OK)
- **Línea 8-9:** Escribe un JSON con el estado y timestamp

#### 2. El ServeMux

```go
mux := http.NewServeMux()
mux.HandleFunc("/health", healthHandler)
```

- `NewServeMux()` crea un nuevo enrutador
- `HandleFunc()` registra la ruta `/health` con su handler

#### 3. Configuración del Servidor

```go
server := &http.Server{
    Addr:         ":8080",              // Puerto
    Handler:      mux,                   // Enrutador
    ReadTimeout:  10 * time.Second,      // Tiempo máximo para leer request
    WriteTimeout: 10 * time.Second,      // Tiempo máximo para escribir response
    IdleTimeout:  120 * time.Second,     // Tiempo máximo de conexión idle
}
```

**¿Por qué son importantes los timeouts?**

- **ReadTimeout:** Previene que clientes lentos bloqueen el servidor
- **WriteTimeout:** Previene que respuestas lentas bloqueen recursos
- **IdleTimeout:** Libera conexiones inactivas

#### 4. Iniciar el Servidor

```go
if err := server.ListenAndServe(); err != nil {
    log.Fatalf("❌ Error al iniciar servidor: %v", err)
}
```

- `ListenAndServe()` bloquea hasta que haya un error
- Si hay error (ej: puerto ocupado), se registra y termina el programa

### Probar el Servidor

**Iniciar:**

```bash
go run main.go
```

**Salida:**

```
🚀 Servidor iniciando en http://localhost:8080
📍 Endpoint disponible: GET /health
```

**Probar con curl:**

```bash
curl http://localhost:8080/health
```

**Respuesta:**

```json
{ "status": "ok", "timestamp": "2025-12-28T10:30:00Z" }
```

**Probar con navegador:**

```
http://localhost:8080/health
```

## Estructura del Código Explicada

```
main.go
│
├── package main              → Paquete ejecutable
│
├── import (...)             → Paquetes necesarios
│
├── func healthHandler(...)  → Handler que procesa /health
│   ├── Set headers
│   ├── Set status code
│   └── Write response
│
└── func main()              → Punto de entrada
    ├── Crear ServeMux       → Enrutador
    ├── Registrar rutas      → Asociar rutas a handlers
    ├── Configurar Server    → Timeouts y configuración
    └── Iniciar servidor     → ListenAndServe
```

## Flujo de una Petición

```
1. Cliente hace petición: GET http://localhost:8080/health
                 ↓
2. Servidor recibe en el puerto 8080
                 ↓
3. ServeMux busca handler para /health
                 ↓
4. Ejecuta healthHandler
                 ↓
5. Handler construye respuesta JSON
                 ↓
6. Servidor envía respuesta al cliente
                 ↓
7. Cliente recibe: {"status": "ok", ...}
```

## Errores Comunes y Soluciones

### Error: "address already in use"

**Problema:** El puerto 8080 ya está en uso

**Solución:**

```bash
# Encontrar el proceso
lsof -i :8080

# Matar el proceso
kill -9 <PID>

# O cambiar el puerto en el código
Addr: ":8081",
```

### Error: "connection refused"

**Problema:** El servidor no está corriendo

**Solución:** Asegúrate de iniciar el servidor con `go run main.go`

### Error: "404 page not found"

**Problema:** La ruta no está registrada

**Solución:** Verifica que la ruta esté correctamente registrada con `HandleFunc`

## Ejercicios Prácticos

### Ejercicio 1: Endpoint de Información

Crea un endpoint `/info` que devuelva:

```json
{
  "app": "gopost-api",
  "version": "1.0.0",
  "description": "API de posts con Go"
}
```

### Ejercicio 2: Endpoint con Timestamp

Crea un endpoint `/time` que devuelva la hora actual en diferentes formatos.

### Ejercicio 3: Múltiples Endpoints

Crea 3 endpoints diferentes:

- `GET /` - Mensaje de bienvenida
- `GET /health` - Estado del servidor
- `GET /version` - Versión de la API

## Resumen

✅ Aprendiste a crear un proyecto Go desde cero
✅ Creaste tu primer servidor HTTP
✅ Entiendes los componentes: Handler, ServeMux, Server
✅ Sabes configurar timeouts para producción
✅ Puedes crear y probar endpoints básicos

En la próxima lección exploraremos el manejo de rutas más avanzado, incluyendo parámetros y diferentes métodos HTTP.

---

**Anterior:** [El Paquete net/http](03-que-net-http.md) | **Siguiente:** [Rutas](05-rutas.md)
