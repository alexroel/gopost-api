# API REST - GoPost

## Endpoints disponibles

### Autenticación

**Registro de usuario**

```bash
POST /api/auth/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "123456"
}
```

**Login**

```bash
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "123456"
}
```

**Obtener usuario autenticado**

```bash
GET /api/auth/me
Authorization: Bearer <token>
```

### Posts

**Obtener todos los posts (público)**

```bash
GET /api/posts
```

**Obtener un post por ID (público)**

```bash
GET /api/posts/{id}
```

**Crear un post (requiere autenticación)**

```bash
POST /api/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Mi primer post",
  "content": "Este es el contenido del post"
}
```

**Actualizar un post (requiere autenticación)**

```bash
PUT /api/posts/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Título actualizado",
  "content": "Contenido actualizado"
}
```

**Eliminar un post (requiere autenticación)**

```bash
DELETE /api/posts/{id}
Authorization: Bearer <token>
```

**Obtener posts del usuario autenticado**

```bash
GET /api/posts/me
Authorization: Bearer <token>
```

## Configuración

1. Crear base de datos y tablas:

```bash
mysql -u root -p < database/schema.sql
```

2. Configurar variables de entorno en `.env`:

```
PORT=:8080
JWT_SECRET=your_jwt_secret_key
DATABASE_URL=user:password@tcp(localhost:3306)/gopost_db?parseTime=true
```

3. Ejecutar el servidor:

```bash
go run main.go
```
