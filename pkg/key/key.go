package key

import (
	"runtime"
	"strings"

	"github.com/samber/lo"
)

type Key struct {
	// The platform specific code for the key. (VK_* on Windows, KEY_* on Linux)
	Code uint32
	// The name of the key.
	Name string
	// The platform to which the key code belongs.
	Platform string
}

// Equals returns true if the key codes are equal.
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

// RegisterModifierKey registers a new modifier key.
func AddModifierKey(key Key) {
	modifierKeys = append(modifierKeys, key)
}

// RemoveModifierKey removes a modifier key.
func RemoveModifierKey(key Key) {
	modifierKeys = lo.Filter(modifierKeys, func(k Key, _ int) bool {
		return k.Code != key.Code
	})
}

// ListModifierKeys returns a list of all modifier keys.
func ListModifierKeys() []Key {
	return modifierKeys
}

// IsModifierKey returns true if the key is a modifier key.
func IsModifierKey(key Key) bool {
	return lo.Contains(lo.Map(modifierKeys, func(k Key, _ int) uint32 { return k.Code }), key.Code)
}

// FindKeyCode returns a Key based on the platform specific code. If the key is
// not found, it returns a Key with the provided code and no name.
func FindKeyCode(val uint32) Key {
	for _, code := range allKeyCodes {
		if code.Code == val {
			return code
		}
	}

	return Key{val, "", runtime.GOOS}
}

// FindKey returns a Key based on the name and a bool indicating if the key was found.
func FindKey(val string) (Key, bool) {
	for _, code := range allKeyCodes {
		if strings.ToLower(code.Name) == strings.ToLower(val) {
			return code, true
		}
	}

	return Key{}, false
}
