package main

import (
	"fmt"

	"github.com/pallandos/benchmarker/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("configs/.env")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded configuration: %+v\n", *cfg)
}
