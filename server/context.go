package server

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	RWriter http.ResponseWriter
	Request *http.Request
	userID  uint
}

func (c *Context) Send(text string) {
	c.RWriter.Write([]byte(text))
}

func (c *Context) Status(code int) {
	c.RWriter.WriteHeader(code)
}

// JSON envía una respuesta en formato JSON
func (c *Context) JSON(code int, data interface{}) error {
	c.RWriter.Header().Set("Content-Type", "application/json")
	c.RWriter.WriteHeader(code)
	return json.NewEncoder(c.RWriter).Encode(data)
}

// BindJSON decodifica el cuerpo de la petición JSON
func (c *Context) BindJSON(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}

// SetUserID establece el ID del usuario en el contexto
func (c *Context) SetUserID(id uint) {
	c.userID = id
}

// GetUserID obtiene el ID del usuario del contexto
func (c *Context) GetUserID() uint {
	return c.userID
}