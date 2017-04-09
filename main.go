package main

import (
	"fmt"
	"os"

	"github.com/davepgreene/turnstile/cmd"
)

func main() {
	if err := cmd.TurnstileCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
