# HTTP y REST

## ¿Qué es HTTP?

**HTTP (Hypertext Transfer Protocol)** es el protocolo de comunicación que permite la transferencia de información en la web. Es un protocolo de petición-respuesta entre un cliente y un servidor.

### Características Principales

- **Sin estado (Stateless)**: Cada petición es independiente
- **Basado en texto**: Los mensajes son legibles por humanos
- **Flexible**: Soporta diferentes tipos de contenido
- **Universal**: Funciona en cualquier plataforma

### Anatomía de una Petición HTTP

```
GET /api/posts HTTP/1.1
Host: ejemplo.com
Content-Type: application/json
Authorization: Bearer token123

{
  "dato": "valor"
}
```

**Componentes:**

1. **Método HTTP** (GET, POST, PUT, DELETE, etc.)
2. **Ruta o Path** (/api/posts)
3. **Versión de HTTP** (HTTP/1.1)
4. **Cabeceras (Headers)** - Metadatos de la petición
5. **Cuerpo (Body)** - Datos de la petición (opcional)

### Anatomía de una Respuesta HTTP

```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 85

{
  "id": 1,
  "titulo": "Mi primer post",
  "contenido": "Contenido del post"
}
```

**Componentes:**

1. **Código de estado** (200, 404, 500, etc.)
2. **Cabeceras de respuesta** - Información sobre la respuesta
3. **Cuerpo de respuesta** - Los datos solicitados

## ¿Qué es REST?

**REST (Representational State Transfer)** es un estilo arquitectónico para diseñar APIs. No es un protocolo ni un estándar, sino un conjunto de principios y restricciones.

### Principios de REST

1. **Cliente-Servidor**: Separación de responsabilidades
2. **Sin estado**: Cada petición contiene toda la información necesaria
3. **Cacheable**: Las respuestas pueden ser cacheadas
4. **Interfaz uniforme**: Uso consistente de recursos y métodos
5. **Sistema en capas**: Arquitectura modular

### Recursos

En REST, todo es un **recurso**. Un recurso es cualquier cosa que puedas nombrar:

- Un usuario
- Un post
- Una imagen
- Una colección de posts

**Ejemplos de URIs de recursos:**

```
/users          → Colección de usuarios
/users/123      → Usuario específico con ID 123
/posts          → Colección de posts
/posts/45       → Post específico con ID 45
/posts/45/comments → Comentarios del post 45
```

## Métodos HTTP y CRUD

Los métodos HTTP se mapean directamente a las operaciones CRUD:

| Método HTTP | Operación CRUD | Descripción                 | Ejemplo                       |
| ----------- | -------------- | --------------------------- | ----------------------------- |
| **POST**    | Create         | Crear un nuevo recurso      | `POST /posts`                 |
| **GET**     | Read           | Obtener recursos            | `GET /posts` o `GET /posts/1` |
| **PUT**     | Update         | Actualizar recurso completo | `PUT /posts/1`                |
| **PATCH**   | Update         | Actualizar parcialmente     | `PATCH /posts/1`              |
| **DELETE**  | Delete         | Eliminar un recurso         | `DELETE /posts/1`             |

### Ejemplos Prácticos

#### 1. Crear un Post (POST)

```http
POST /posts
Content-Type: application/json

{
  "titulo": "Aprendiendo Go",
  "contenido": "Go es un lenguaje increíble"
}
```

**Respuesta:**

```http
201 Created
Location: /posts/1

{
  "id": 1,
  "titulo": "Aprendiendo Go",
  "contenido": "Go es un lenguaje increíble",
  "createdAt": "2025-12-28T10:00:00Z"
}
```

#### 2. Obtener Todos los Posts (GET)

```http
GET /posts
```

**Respuesta:**

```http
200 OK

[
  {
    "id": 1,
    "titulo": "Aprendiendo Go",
    "contenido": "Go es un lenguaje increíble"
  },
  {
    "id": 2,
    "titulo": "REST APIs",
    "contenido": "Construyendo APIs REST"
  }
]
```

#### 3. Obtener un Post Específico (GET)

```http
GET /posts/1
```

**Respuesta:**

```http
200 OK

{
  "id": 1,
  "titulo": "Aprendiendo Go",
  "contenido": "Go es un lenguaje increíble"
}
```

#### 4. Actualizar un Post (PUT)

```http
PUT /posts/1
Content-Type: application/json

{
  "titulo": "Aprendiendo Go - Actualizado",
  "contenido": "Go es el mejor lenguaje para APIs"
}
```

**Respuesta:**

```http
200 OK

{
  "id": 1,
  "titulo": "Aprendiendo Go - Actualizado",
  "contenido": "Go es el mejor lenguaje para APIs"
}
```

#### 5. Eliminar un Post (DELETE)

```http
DELETE /posts/1
```

**Respuesta:**

```http
204 No Content
```

## Códigos de Estado HTTP (Status Codes)

Los códigos de estado indican el resultado de la petición. Se agrupan en 5 categorías:

### 1xx: Informativos

Raramente usados en APIs REST.

### 2xx: Éxito

| Código  | Significado | Uso Común                             |
| ------- | ----------- | ------------------------------------- |
| **200** | OK          | Petición exitosa con respuesta        |
| **201** | Created     | Recurso creado exitosamente           |
| **204** | No Content  | Éxito pero sin contenido de respuesta |

### 3xx: Redirección

| Código  | Significado       | Uso Común                      |
| ------- | ----------------- | ------------------------------ |
| **301** | Moved Permanently | Recurso movido permanentemente |
| **304** | Not Modified      | Recurso no modificado (cache)  |

### 4xx: Errores del Cliente

| Código  | Significado          | Uso Común                                  |
| ------- | -------------------- | ------------------------------------------ |
| **400** | Bad Request          | Datos inválidos en la petición             |
| **401** | Unauthorized         | No autenticado                             |
| **403** | Forbidden            | No autorizado (sin permisos)               |
| **404** | Not Found            | Recurso no encontrado                      |
| **409** | Conflict             | Conflicto (ej: email duplicado)            |
| **422** | Unprocessable Entity | Datos válidos pero lógicamente incorrectos |

### 5xx: Errores del Servidor

| Código  | Significado           | Uso Común                            |
| ------- | --------------------- | ------------------------------------ |
| **500** | Internal Server Error | Error genérico del servidor          |
| **502** | Bad Gateway           | Error en servidor upstream           |
| **503** | Service Unavailable   | Servicio temporalmente no disponible |

## Mejores Prácticas REST

### 1. Usa Sustantivos, No Verbos

❌ **Incorrecto:**

```
POST /createPost
GET /getUser
DELETE /deletePost
```

✅ **Correcto:**

```
POST /posts
GET /users
DELETE /posts/1
```

### 2. Usa Plurales para Colecciones

✅ **Recomendado:**

```
GET /posts          → Lista de posts
GET /posts/1        → Un post específico
GET /users          → Lista de usuarios
GET /users/5        → Un usuario específico
```

### 3. Usa Jerarquías para Relaciones

```
GET /posts/1/comments        → Comentarios del post 1
GET /users/5/posts           → Posts del usuario 5
POST /posts/1/comments       → Crear comentario en post 1
```

### 4. Devuelve el Código de Estado Apropiado

```go
// Crear: 201 Created
// Obtener: 200 OK
// Actualizar: 200 OK
// Eliminar: 204 No Content
// Error: 400, 404, 500, etc.
```

### 5. Usa Content-Type Correcto

```http
Content-Type: application/json
Content-Type: application/xml
Content-Type: text/html
```

## Resumen

- **HTTP** es el protocolo de comunicación de la web
- **REST** es un estilo arquitectónico para diseñar APIs
- Los **métodos HTTP** mapean a operaciones CRUD
- Los **códigos de estado** comunican el resultado de la petición
- Las **mejores prácticas** hacen tu API más intuitiva y mantenible

## Ejercicio Práctico

Diseña los endpoints REST para un sistema de biblioteca con los siguientes recursos:

- Libros (books)
- Autores (authors)
- Préstamos (loans)

Define:

1. Las URIs para cada recurso
2. Los métodos HTTP necesarios
3. Los códigos de estado esperados
4. Ejemplos de petición y respuesta

---

**Anterior:** [Introducción](01-intro.md) | **Siguiente:** [El Paquete net/http](03-que-net-http.md)
