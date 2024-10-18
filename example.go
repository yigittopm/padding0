package main

type A struct {
	A byte  // AA Comment
	B int32 // AB Comment
	C byte  // AC Comment
	D int64 // AD Comment
	E byte  // AE Comment
}

type X byte

type B struct {
	BA X     // BAA Comment
	BB int32 // BAB Comment
	BC X     // BAC Comment
	BD int64 // BAD Comment
	BE X     // BAE Comment
}
