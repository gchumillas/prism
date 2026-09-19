package main

import (
	"fmt"
	"os"

	"github.com/gchumillas/prism/internal/platform"
	"github.com/gchumillas/prism/internal/storybook"
)

func main() {
	var env *platform.Env
	var err error
	if env, err = platform.LoadRuntimeEnv(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	sb := &storybook.Storybook{
		Touch:    env.Touch,
		Screen:   env.Screen,
		Canvas:   env.Canvas,
		FontSize: env.FontSize,
		DPI:      env.DPI,
	}
	if err := sb.Run(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	env.Touch.Listen()
}
