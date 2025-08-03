package keyboard

import (
	"context"
	"syscall"
	"time"
	"unsafe"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
)

const (
	// This is the event we will hook into in the Windows API.
	WH_KEYBOARD_LL int = 13
	// LLKHF_UP is used to determine if a key event is an up or down event.
	LLKHF_UP = 0x00000080
)

var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessageW         = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessageW    = user32.NewProc("DispatchMessageW")

	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procGetTickCount = kernel32.NewProc("GetTickCount64")

	hook hHOOK
)

type (
	hHOOK uintptr
	hWND  uintptr
	msg   struct {
		HWND    hWND
		Message uint32
		WPARAM  uintptr
		LPARAM  uintptr
		Time    uint32
		Pt      struct {
			X, Y int32
		}
	}
	kbDLLHOOKSTRUCT struct {
		VKCode      uint32
		ScanCode    uint32
		Flags       uint32
		Time        uint32
		DwExtraInfo uintptr
	}
)

func init() {
	keyboard = &keyboardWindows{}
}

type keyboardWindows struct {
	eventChan chan Event
}

func (k *keyboardWindows) Events(ctx context.Context) (chan Event, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan Event)
		go func() {
			h, _, _ := procSetWindowsHookEx.Call(
				uintptr(WH_KEYBOARD_LL),
				syscall.NewCallback(k.hook(ctx)),
				0,
				0,
			)
			hook = hHOOK(h)
			defer procUnhookWindowsHookEx.Call(uintptr(hook))

			k.executeMessageLoop(ctx)
		}()
	}
	return k.eventChan, nil
}

type HookFunc func(int, uintptr, uintptr) uintptr

// hook is the callback that will be invoked when a WH_KEYBOARD_LL event is received
// from the Windows API.
func (k *keyboardWindows) hook(ctx context.Context) HookFunc {
	return func(nCode int, wParam, lParam uintptr) uintptr {
		select {
		case <-ctx.Done():
			return 0

		default:
			if nCode >= 0 {
				var (
					// ignore "possible misuse of unsafe.Pointer" warning, lParam is guaranteed by
					// the Windows API to be a pointer to a kbDLLHOOKSTRUCT.
					kbStruct = (*kbDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

					// The event we're going to build to send to the main channel.
					event Event
					// Assume that the key event is a down event.
					eventType = KeyboardEventTypeDown
				)

				// Check for up key event
				if kbStruct.Flags&LLKHF_UP != 0 {
					eventType = KeyboardEventTypeUp
				}

				// Build the event
				event = Event{
					Key:       key.FindKeyCode(kbStruct.VKCode),
					EventType: eventType,
					Timestamp: parseTime(kbStruct.Time),
				}

				// Send the event to the main channel.
				select {
				case k.eventChan <- event:
				case <-ctx.Done():
				}
			}

			r1, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
			return r1
		}
	}
}

// executeMessageLoop is an internal function used to listen for events from the
// Windows API. It blocks it's caller and calls GetMessageW(), TranslateMessage()
// and DispatchMessageW() in a loop until the provided context is Done().
func (k *keyboardWindows) executeMessageLoop(ctx context.Context) {
	var msg msg
	for {
		select {
		case <-ctx.Done():
			close(k.eventChan)
			k.eventChan = nil
			return

		default:
			ret, _, _ := procGetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
			if ret != 0 {
				procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
				procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
			}
		}
	}
}

func parseTime(t uint32) time.Time {
	sysTime, _, _ := procGetTickCount.Call()
	now := time.Now()
	bootTime := now.Add(-time.Duration(int64(sysTime)) * time.Millisecond)
	return bootTime.Add(time.Duration(t) * time.Millisecond)
}
