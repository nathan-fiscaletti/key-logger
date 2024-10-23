package key

import (
	"github.com/nathan-fiscaletti/key-logger/pkg/platform"
	"github.com/samber/lo"
)

type Key struct {
	Code     uint32
	Name     string
	Platform platform.Platform
}

func (k Key) Equals(other Key) bool {
	return k.Code == other.Code
}

var allKeyCodes = []Key{
	Backspace,
	Tab,
	Enter,
	Escape,
	Space,
	PageUp,
	PageDown,
	End,
	Home,
	Left,
	Up,
	Right,
	Down,
	Insert,
	Delete,
	Number0,
	Number1,
	Number2,
	Number3,
	Number4,
	Number5,
	Number6,
	Number7,
	Number8,
	Number9,
	A,
	B,
	C,
	D,
	E,
	F,
	G,
	H,
	I,
	J,
	K,
	L,
	M,
	N,
	O,
	P,
	Q,
	R,
	S,
	T,
	U,
	V,
	W,
	X,
	Y,
	Z,
	LeftWin,
	RightWin,
	NumPad0,
	NumPad1,
	NumPad2,
	NumPad3,
	NumPad4,
	NumPad5,
	NumPad6,
	NumPad7,
	NumPad8,
	NumPad9,
	NumPadMultiply,
	NumPadAdd,
	NumPadSubtract,
	NumPadDecimal,
	NumPadDivide,
	F1,
	F2,
	F3,
	F4,
	F5,
	F6,
	F7,
	F8,
	F9,
	F10,
	F11,
	F12,
	F13,
	F14,
	F15,
	F16,
	F17,
	F18,
	F19,
	F20,
	F21,
	F22,
	F23,
	F24,
	NumLock,
	ScrollLock,
	LeftShift,
	RightShift,
	LeftControl,
	RightControl,
	LeftAlt,
	RightAlt,
	Backtick,
	LeftBracket,
	RightBracket,
	Backslash,
	SemiColon,
	CapsLock,
	Plus,
	Comma,
	Minus,
	Period,
	Slash,
	Quote,
}

var modifierKeys = []Key{
	LeftShift,
	RightShift,
	LeftControl,
	RightControl,
	LeftAlt,
	RightAlt,
	LeftWin,
	RightWin,
}

func IsModifierKey(key Key) bool {
	return lo.Contains(lo.Map(modifierKeys, func(k Key, _ int) uint32 { return k.Code }), key.Code)
}

func FindKeyCode(val uint32) Key {
	for _, code := range allKeyCodes {
		if code.Code == val {
			return code
		}
	}

	return Key{val, "", platform.GetPlatform()}
}