package server

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	mux          *http.ServeMux
	handlerCount int
}

func New() *App {
	return &App{
		mux:          http.NewServeMux(),
		handlerCount: 0,
	}
}

func (a *App) RunServer(addr string) error {
	// Mostrar el banner antes de iniciar el servidor
	a.printBanner(addr)

	// Configurar servidor con timeouts
	srv := &http.Server{
		Addr:    addr,
		Handler: a.mux,
	}

	return srv.ListenAndServe()
}

func (a *App) printBanner(addr string) {
	urlBase := fmt.Sprintf("http://localhost%s", addr)

	fmt.Println("┌───────────────────────────────────────────────────┐")
	fmt.Printf("│%s│\n", centerText("MyServer v1.0.0", 51))
	fmt.Printf("│%s│\n", centerText(urlBase, 51))
	fmt.Printf("│%s│\n", strings.Repeat(" ", 51))
	fmt.Printf("│ Handlers .........: %d|\n", a.handlerCount)
	fmt.Println("└───────────────────────────────────────────────────┘")
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}
