package docktest

import (
	"log"
)

func header() {
	log.Println("\033[1;36m", "»» DOCKTEST ««", "\033[0m")
}

func Error(msg ...string) {
	header()
	log.Print("\033[31m", "🚨", msg, "\033[0m", "\n")
}

func Warn(msg ...string) {
	header()
	log.Print("\033[33m", "⚠️", msg, "\033[0m", "\n")
}

func Info(msg ...string) {
	header()
	log.Print("\033[34m", "ℹ️", msg, "\033[0m", "\n")
}

func Success(msg ...string) {
	header()
	log.Print("\033[32m", "✅", msg, "\033[0m", "\n")
}
