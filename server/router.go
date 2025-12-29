package server

import "net/http"

type HandleFunc func(c *Context)

func (a *App) Get(path string, handler HandleFunc) {
	a.mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	a.handlerCount++
}

func (a *App) Post(path string, handler HandleFunc) {
	a.mux.HandleFunc("POST "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	a.handlerCount++
}

func (a *App) Put(path string, handler HandleFunc) {
	a.mux.HandleFunc("PUT "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	a.handlerCount++
}

func (a *App) Delete(path string, handler HandleFunc) {
	a.mux.HandleFunc("DELETE "+path, func(w http.ResponseWriter, r *http.Request) {
		handler(&Context{
			RWriter: w,
			Request: r,
			Ctx:     r.Context(),
		})
	})
	a.handlerCount++
}
