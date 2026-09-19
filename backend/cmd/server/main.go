package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"polling-tool/backend/internal/app"
)

func main() {
	cfg := app.LoadConfig()
	a, err := app.New(context.Background(), cfg)
	if err != nil { log.Fatal(err) }
	go func() { if err := a.Serve(":" + cfg.Port); err != nil && err != http.ErrServerClosed { log.Printf("server stopped: %v", err) } }()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second); defer cancel()
	if err := a.Close(ctx); err != nil { log.Printf("shutdown: %v", err) }
}
