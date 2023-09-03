package signal

// Slice allows to slice arbitrary signal type.
func Slice(s Signal, start, end int) Signal {
	switch b := s.(type) {
	case Signed:
		return b.Slice(start, end)
	case Unsigned:
		return b.Slice(start, end)
	case Floating:
		return b.Slice(start, end)
	}
	return s
}

// AsFloating allows to convert arbitrary signal type to floating.
func AsFloating(s Signal, dst Floating) (read int) {
	// switch src := s.(type) {
	// case Signed:
	// 	read = SignedAsFloat(src, dst)
	// case Unsigned:
	// 	read = UnsignedAsFloat(src, dst)
	// case Floating:
	// 	read = FloatingAsFloating(src, dst)
	// }
	return
}

// AsSigned allows to convert arbitrary signal type to signed.
func AsSigned(s Signal, dst Signed) (read int) {
	// switch src := s.(type) {
	// case Signed:
	// 	read = SignedAsSigned(src, dst)
	// case Unsigned:
	// 	read = UnsignedAsSigned(src, dst)
	// case Floating:
	// 	read = FloatingAsSigned(src, dst)
	// }
	return
}

// AsUnsigned allows to convert arbitrary signal type to unsigned.
func AsUnsigned(s Signal, dst Unsigned) (read int) {
	// switch src := s.(type) {
	// case Signed:
	// 	read = SignedAsUnsigned(src, dst)
	// case Unsigned:
	// 	read = UnsignedAsUnsigned(src, dst)
	// case Floating:
	// 	read = FloatAsUnsigned(src, dst)
	// }
	return
}
