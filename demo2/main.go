package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	myHandler := http.NewServeMux()
	myHandler.HandleFunc("/myhttp", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("myhttp response"))
	})
	svr := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	// shutdown
	myShutDown := make(chan struct{})
	myHandler.HandleFunc("/myshutdown", func(w http.ResponseWriter, r *http.Request) {
		myShutDown <- struct{}{}
	})
	// signal
	g.Go(func() error {
		sig := make(chan os.Signal, 0)
		signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case <-ctx.Done():
				log.Println("ctx.Done()...")
				return ctx.Err()
			case <-sig:
				return errors.Errorf("get os signal: %v", sig)
			}
		}
	})
	// http server
	g.Go(func() error {
		return svr.ListenAndServe()
	})
	// myShutDown
	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit...")
		case <-myShutDown:
			log.Println("server will out...")
		}

		log.Println("shutting down server...")
		return svr.Shutdown(ctx)
	})
	// errgroup exit
	fmt.Println(g.Wait())
}
