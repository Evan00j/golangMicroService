package main

import (
	"context"
	"log"
	"os"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

type Env struct {
	db *pop.Connection
}

func main() {

	ctx := context.Background()
	errChan := make(chan error)
	var svc Service
	svc = URLService{}

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	db, err := pop.Connect("development")
	logger.Log(err)
	if err != nil {
		logger.Log(err)
	}

	env := &Env{db: db}

	endpoint := Endpoints{
	}

	r := MakeHttpHandler(ctx, endpoint, logger)

	// HTTP transport
	go func() {
		fmt.Println("Server started at port 8080")
		handler := r
		errChan <- http.ListenAndServe(":8080", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	fmt.Println(<-errChan)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}