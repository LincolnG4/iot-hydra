package main

import (
	"fmt"

	"github.com/LincolnG4/iot-hydra/internal/runtime"
)

func main() {
	podmanRuntime := runtime.NewPodmanRuntime()

	err := podmanRuntime.Start()
	if err != nil {
		fmt.Println(err)
	}
}
