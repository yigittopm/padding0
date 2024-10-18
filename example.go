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

type C struct {
	CA X     // CAA Comment
	CB int32 // CAB Comment
	CC struct {
		C1 X     // CCAA Comment
		C2 int32 // CCAB Comment
		C3 X     // CCAC Comment
		C4 int64 // CCAD Comment
	}
	CD X   // CAD Comment
	CE int // CAE Comment
}
