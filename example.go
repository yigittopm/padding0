package main

type X byte  // X is a byte
type Y int64 // Y is an int64

type C struct {
	C1 X     // C1
	C2 X     // C2
	C3 int32 // C3
	C4 struct {
		C5 X     // C5
		C6 int32 // C6
		C7 X     // C7
		C8 int64 // C8
	}
	C9  X   // C9
	C10 Y   // C10
	C11 int // C11
	C12 map[string]struct {
		C13 int      // C13
		C14 []string // C14
	}
}
