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
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	group, errCtx := errgroup.WithContext(ctx)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	group.Go(func() error {
		return server.ListenAndServe()
	})

	group.Go(func() error {
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
