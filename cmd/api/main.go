package main

import (
	"fmt"

	"github.com/Netflix/go-env"
	"github.com/fabulias/metrics-api/internal/app"
	"github.com/fabulias/metrics-api/internal/config"
)

func main() {
	var conf config.Config
	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		panic(fmt.Errorf("cannot load env vars properly: %v", err))
	}

	app := app.NewApplication(conf)

	app.Run()
}
