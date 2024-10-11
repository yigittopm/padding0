package padding0

func AFunc() {
	println("A")
}

type A struct {
	A byte  // AA Comment
	B int32 // AB Comment
	C byte  // AC Comment
	D int64 // AD Comment
	E byte  // AE Comment
}

type Custom int64

type B struct {
	A byte  // AA Comment
	B int32 // AB Comment
	C byte  // AC Comment
	D int64 // AD Comment
	E byte  // AE Comment
}

func BFunc() {
	println("B")
}
