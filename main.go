/*
Copyright © 2024 Mert Yiğittop <yigittopm@hotmail.com>
*/
package main

import (
	analysis "github.com/yigittopm/padding0/internal"
)

func main() {
	if err := analysis.Start(); err != nil {
		panic(err)
	}
}
