package utils

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

var clipboard ClipboardService

// Clipboard returns an object that provides access to the system clipboard.
func Clipboard() *ClipboardService {
	return &clipboard
}

// ClipboardService provides access to the system clipboard.
type ClipboardService struct {
	hwnd                     win.HWND
	contentsChangedPublisher walk.EventPublisher
}

// ContentsChanged returns an Event that you can attach to for handling
// clipboard content changes.
func (c *ClipboardService) ContentsChanged() *walk.Event {
	return c.contentsChangedPublisher.Event()
}

// Clear clears the contents of the clipboard.
func (c *ClipboardService) Clear() error {
	return c.withOpenClipboard(func() error {
		if !win.EmptyClipboard() {
			return lastError("EmptyClipboard")
		}

		return nil
	})
}

// ContainsText returns whether the clipboard currently contains text data.
func (c *ClipboardService) ContainsText() (available bool, err error) {
	err = c.withOpenClipboard(func() error {
		available = win.IsClipboardFormatAvailable(win.CF_UNICODETEXT)

		return nil
	})

	return
}

// Text returns the current text data of the clipboard.
func (c *ClipboardService) Text() (text string, err error) {
	err = c.withOpenClipboard(func() error {
		hMem := win.HGLOBAL(win.GetClipboardData(win.CF_UNICODETEXT))
		if hMem == 0 {
			return lastError("GetClipboardData")
		}

		p := win.GlobalLock(hMem)
		if p == nil {
			return lastError("GlobalLock()")
		}
		defer win.GlobalUnlock(hMem)

		text = win.UTF16PtrToString((*uint16)(p))

		return nil
	})

	return
}

// SetText sets the current text data of the clipboard.
func (c *ClipboardService) SetText(s string) error {
	return c.withOpenClipboard(func() error {
		utf16, err := syscall.UTF16FromString(s)
		if err != nil {
			return err
		}

		hMem := win.GlobalAlloc(win.GMEM_MOVEABLE, uintptr(len(utf16)*2))
		if hMem == 0 {
			return lastError("GlobalAlloc")
		}

		p := win.GlobalLock(hMem)
		if p == nil {
			return lastError("GlobalLock()")
		}

		win.MoveMemory(p, unsafe.Pointer(&utf16[0]), uintptr(len(utf16)*2))

		win.GlobalUnlock(hMem)

		if 0 == win.SetClipboardData(win.CF_UNICODETEXT, win.HANDLE(hMem)) {
			// We need to free hMem.
			defer win.GlobalFree(hMem)

			return lastError("SetClipboardData")
		}

		// The system now owns the memory referred to by hMem.

		return nil
	})
}

type DROPFILES struct {
	pFiles uintptr
	pt     uintptr
	fNC    bool
	fWide  bool
}

// SetFile sets the current file drop data of the clipboard.
func (c *ClipboardService) SetFile(s string) error {
	return c.withOpenClipboard(func() error {
		utf16, err := syscall.UTF16FromString(s)
		if err != nil {
			return err
		}

		size := unsafe.Sizeof(DROPFILES{}) + uintptr((len(utf16)+2)*2)

		hMem := win.GlobalAlloc(win.GMEM_MOVEABLE, size)
		if hMem == 0 {
			return lastError("GlobalAlloc")
		}

		p := win.GlobalLock(hMem)
		if p == nil {
			return lastError("GlobalLock()")
		}

		zeroMem := make([]byte, size)
		win.MoveMemory(p, unsafe.Pointer(&zeroMem[0]), size)

		pD := (*DROPFILES)(p)
		pD.pFiles = unsafe.Sizeof(DROPFILES{})
		pD.fWide = true
		win.MoveMemory(unsafe.Pointer(uintptr(p)+unsafe.Sizeof(DROPFILES{})), unsafe.Pointer(&utf16[0]), uintptr(len(utf16)*2))

		win.GlobalUnlock(hMem)

		if 0 == win.SetClipboardData(win.CF_HDROP, win.HANDLE(hMem)) {
			// We need to free hMem.
			defer win.GlobalFree(hMem)

			return lastError("SetClipboardData")
		}

		// The system now owns the memory referred to by hMem.

		return nil
	})
}

func (c *ClipboardService) withOpenClipboard(f func() error) error {
	if !win.OpenClipboard(c.hwnd) {
		return lastError("OpenClipboard")
	}
	defer win.CloseClipboard()

	return f()
}

func lastError(name string) error {
	return errors.New(fmt.Sprintf("%s failed", name))
}
