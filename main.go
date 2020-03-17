//go:generate sqlboiler --wipe --no-hooks --add-panic-variants psql

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}
