package main

import (
	"os"

	"github.com/aslammmuhammed/run-secret-reloader/internal/app"
)

func main() {
	exitCode := app.Run()
	os.Exit(exitCode)
}
