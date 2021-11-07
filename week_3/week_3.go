package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	// HTTP handler
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	ctx := context.Background()
	// Use cancel() to cancel downstream context
	ctx, cancel := context.WithCancel(ctx)
	// Use errgroup to cancel goroutines
	group, errCtx := errgroup.WithContext(ctx)

	// Start a server
	group.Go(func() error {
		return server.ListenAndServe()
	})

	group.Go(func() error {
		// block until cancel() closes Done
		<-errCtx.Done()
		fmt.Println("Shutting down server...")
		return server.Shutdown(errCtx)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c)

	group.Go(func() error {
		for {
			select {
			case <-errCtx.Done():
				return errCtx.Err()
			// block until os signals are captured
			case <-c:
				cancel()
			}
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Println("group error: ", err)
	}
	fmt.Println("all group done.")
}
