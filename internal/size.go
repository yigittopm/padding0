package analysis

var typeSizes = map[string]int{
	"bool":       1,  // 1 byte
	"byte":       1,  // 1 byte
	"int8":       1,  // 1 byte
	"uint8":      1,  // 1 byte
	"int16":      2,  // 2 byte
	"uint16":     2,  // 2 byte
	"int32":      4,  // 4 byte
	"uint32":     4,  // 4 byte
	"int64":      8,  // 8 byte
	"uint64":     8,  // 8 byte
	"int":        8,  // 8 byte (64-bit)
	"uint":       8,  // 8 byte (64-bit)
	"float32":    4,  // 4 byte
	"float64":    8,  // 8 byte
	"complex64":  8,  // 8 byte
	"complex128": 16, // 16 byte

	"uintptr":        8,  // 8 byte (64-bit)
	"string":         16, // 16 byte (header size, actual string data is elsewhere)
	"unsafe.Pointer": 8,  // 8 byte

	"slice":     24, // 24 byte (header size, actual data is elsewhere)
	"map":       8,  // Maps are references (pointer to internal structures)
	"chan":      8,  // Channels are references (pointer to internal structures)
	"interface": 16, // 16 byte (2 words, type info and data)
}
