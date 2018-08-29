package logger

import (
	"log"
	"testing"
)

func TestGetCurrentPath(t *testing.T) {
	log.Println(GetCurrentPath())
}
