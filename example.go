package main

type X byte
type Y int64

type C struct {
	XC X     // CX Comment
	CA X     // CAA Comment
	CB int32 // CAB Comment
	CC struct {
		C1 X     // CCAA Comment
		C2 int32 // CCAB Comment
		C3 X     // CCAC Comment
		C4 int64 // CCAD Comment
	}
	CD X   // CAD Comment
	DF Y   // CDF Comment
	CE int // CAE Comment
	XX map[string]struct {
		YY int      // CXXYY Comment
		ST []string // ST Comment
	}
}
