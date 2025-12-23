package middleware

import (
	"fmt"

	"github.com/fatih/color"
)

func OK(msg string) {
	fmt.Println(color.GreenString("✔ %s", msg))
}

func Warn(msg string) {
	fmt.Println(color.YellowString("⚠ %s", msg))
}

func Fail(msg string) {
	fmt.Println(color.RedString("✖ %s", msg))
}

func Info(msg string) {
	fmt.Println(color.CyanString("ℹ %s", msg))
}
