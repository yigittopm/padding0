package main

import (
	analysis "github.com/yigittopm/padding0/internal"
)

func main() {
	println("Padding0 starting")
	if err := analysis.Start(); err != nil {
		panic(err)
	}
	println("Padding0 finished")
}
