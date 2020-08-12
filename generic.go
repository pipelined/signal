package signal

// Slice is a generic wrapper of Slice function for all signal types.
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

// AsFloating is a generic wrapper for conversion to floating for any signal type.
func AsFloating(s Signal, dst Floating) (read int) {
	switch src := s.(type) {
	case Signed:
		read = SignedAsFloating(src, dst)
	case Unsigned:
		read = UnsignedAsFloating(src, dst)
	case Floating:
		read = FloatingAsFloating(src, dst)
	}
	return
}

// AsSigned is a generic wrapper for conversion to signed for any signal type.
func AsSigned(s Signal, dst Signed) (read int) {
	switch src := s.(type) {
	case Signed:
		read = SignedAsSigned(src, dst)
	case Unsigned:
		read = UnsignedAsSigned(src, dst)
	case Floating:
		read = FloatingAsSigned(src, dst)
	}
	return
}

// AsUnsigned is a generic wrapper for conversion to unsigned for any signal type.
func AsUnsigned(s Signal, dst Unsigned) (read int) {
	switch src := s.(type) {
	case Signed:
		read = SignedAsUnsigned(src, dst)
	case Unsigned:
		read = UnsignedAsUnsigned(src, dst)
	case Floating:
		read = FloatingAsUnsigned(src, dst)
	}
	return
}
