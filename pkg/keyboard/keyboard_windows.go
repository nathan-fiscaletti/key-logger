package keyboard

import (
	"context"
	"syscall"
	"time"
	"unsafe"

	"github.com/nathan-fiscaletti/key-logger/pkg/key"
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

	hook         hHOOK
	keyboardHook int = 13 // WH_KEYBOARD_LL is 13
)

type (
	hHOOK     uintptr
	hWND      uintptr
	hHANDLE   uintptr
	hINSTANCE hHANDLE
	hHOOKPROC uintptr
	msg       struct {
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

const (
	flagKeyUp = 0x8000 >> 8
)

func init() {
	keyboard = keyboardWindows{}
}

type keyboardWindows struct {
	eventChan chan Event
}

func (k keyboardWindows) Events(ctx context.Context) (chan Event, error) {
	if k.eventChan == nil {
		k.eventChan = make(chan Event)
		go func() {
			hookCallback := syscall.NewCallback(func(nCode int, wParam, lParam uintptr) uintptr {
				select {
				case <-ctx.Done():
					return 0

				default:
					if nCode >= 0 {
						kbStruct := (*kbDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

						var eventType EventType = KeyboardEventTypeDown
						if kbStruct.Flags == flagKeyUp {
							eventType = KeyboardEventTypeUp
						}

						event := Event{
							Key:       key.FindKeyCode(kbStruct.VKCode),
							EventType: eventType,
							Timestamp: parseTime(kbStruct.Time),
						}

						k.eventChan <- event
					}

					r1, _, _ := procCallNextHookEx.Call(0, uintptr(nCode), wParam, lParam)
					return r1
				}
			})

			h, _, _ := procSetWindowsHookEx.Call(uintptr(keyboardHook), hookCallback, 0, 0)
			hook = hHOOK(h)
			defer procUnhookWindowsHookEx.Call(uintptr(hook))

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
		}()
	}

	return k.eventChan, nil
}

func parseTime(t uint32) time.Time {
	sysTime, _, _ := procGetTickCount.Call()
	now := time.Now()
	bootTime := now.Add(-time.Duration(int64(sysTime)) * time.Millisecond)
	return bootTime.Add(time.Duration(t) * time.Millisecond)
}
