package main

import (
	"github.com/pak-app/gosuper/cmd"
	"github.com/pak-app/gosuper/internal/logging"
)

func main() {
	logging.Init("log/gosuper.log")

	defer logging.Close()

	cmd.Execute()
}
