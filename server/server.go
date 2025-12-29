package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
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
		Addr:           addr,
		Handler:        a.mux,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
	
	return srv.ListenAndServe()
}

func (a *App) printBanner(addr string) {
	host := "0.0.0.0"
	port := strings.TrimPrefix(addr, ":")
	
	pid := os.Getpid()
	
	fmt.Println("┌───────────────────────────────────────────────────┐")
	fmt.Printf("│%s│\n", centerText("MyServer v1.0.0", 51))
	fmt.Printf("│%s│\n", centerText(fmt.Sprintf("http://localhost:%s", port), 51))
	fmt.Printf("│%s│\n", centerText(fmt.Sprintf("(bound on host %s and port %s)", host, port), 51))
	fmt.Printf("│%s│\n", strings.Repeat(" ", 51))
	fmt.Printf("│ Handlers .........%s%d  Processes .........%s%d │\n", 
		leftPad(fmt.Sprint(a.handlerCount), 3),
		a.handlerCount,
		leftPad("", 0),
		1)
	fmt.Printf("│ Prefork ....... Disabled  PID .............%s%d │\n",
		leftPad("", 0),
		pid)
	fmt.Println("└───────────────────────────────────────────────────┘")
}

func centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

func leftPad(text string, minSpaces int) string {
	if minSpaces > 0 {
		return strings.Repeat(" ", minSpaces-len(text))
	}
	return ""
}