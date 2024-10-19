package main

type A struct {
	A byte  // AA Comment
	B int32 // AB Comment
	C byte  // AC Comment
	D int64 // AD Comment
	E byte  // AE Comment
}

type X byte
type Y int64

type B struct {
	BI X     // BIA Comment
	BA Y     // BAA Comment
	BB int32 // BAB Comment
	BC byte  // BAC Comment
	BD X     // BAD Comment
	BE Y     // BAE Comment
	BF int64 // BAF Comment
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
