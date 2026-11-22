package main

import (
	"fmt"
	"os"

	"github.com/gchumillas/prism/internal/app"
	"github.com/gchumillas/prism/internal/platform"
)

func main() {
	env, err := platform.LoadRuntime()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	a := app.NewApp(env.Touch, env.Screen, env.Canvas, env.FontSize)
	if err := a.Refresh(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
