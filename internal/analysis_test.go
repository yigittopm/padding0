/*
Copyright © 2024 Mert Yiğittop <yigittopm@hotmail.com>
*/
package analysis

import "testing"

func TestProcessStruct(t *testing.T) {
	beforeContent := `
	package analysis
	
	type Person struct {
		isActive bool
		Age  int	
		Name string
		Height float64
		b  byte
		X int64
	}
	`

	want := `
	package analysis

	type Person struct {
		Name string 
		Age int 
		Height float64 
		X int64 
		isActive bool 
		b byte 
	}
	`

	got := ExtractStruct(beforeContent)
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}

}
