package main
import (
	"fmt"
	"go.uber.org/zap/zapcore"
)
func main() {
	l, err := zapcore.ParseLevel("")
	fmt.Printf("Level: %v, Err: %v\n", l, err)
}
