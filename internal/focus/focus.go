package focus

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type FocusableWindow interface {
	Id() xproto.Window

	IsFocusable() bool // whether the window is in state it can be focused (mapped)
	IsFocused() bool

	CanFocus() bool
	ShouldSendFocusNotify() bool // WM_TAKE_FOCUS protocol

	PrepareForFocus(tmp bool)
	ApplyFocus() // set focus on the actual window, send focus notify
	Focused()
	Unfocused()
}

var x *xgbutil.XUtil
var windows []FocusableWindow

func Initialize(_x *xgbutil.XUtil) {
	x = _x
}

func Current() FocusableWindow {
	if len(windows) == 0 {
		return nil
	}

	if w := windows[len(windows)-1]; w.IsFocused() {
		return w
	}
	return nil
}

func InitialAdd(w FocusableWindow) {
	windows = append([]FocusableWindow{w}, windows...)
}

func Remove(w FocusableWindow) bool {
	for i, w2 := range windows {
		if w.Id() == w2.Id() {
			windows = append(windows[:i], windows[i+1:]...)
			return true
		}
	}
	return false
}

func Focus(w FocusableWindow) {
	focus(w, false)
}

func focus(w FocusableWindow, tmp bool) {
	if !Remove(w) {
		return
	}

	if w.CanFocus() || w.ShouldSendFocusNotify() {
		add(w)
		w.PrepareForFocus(tmp)
		w.ApplyFocus()
	}
}

func FocusLast() {
	FocusLastWithPreference(func(window xproto.Window) bool {
		return true
	})
}

func FocusLastWithPreference(isPreferred func(window xproto.Window) bool) {
	if w := lastFocused(isPreferred); w != nil {
		Focus(w)
	} else {
		xwindow.New(x, x.Dummy()).Focus()
	}
}

func SetFocus(w FocusableWindow) {
	Remove(w)
	add(w)
}

func lastFocused(isPreferred func(window xproto.Window) bool) FocusableWindow {
	var firstNonPreferred FocusableWindow
	last := len(windows) - 1
	for i := range windows {
		if w := windows[last-i]; w.IsFocusable() {
			if isPreferred(w.Id()) {
				return w
			} else if firstNonPreferred == nil {
				firstNonPreferred = w
			}
		}
	}
	return firstNonPreferred
}

func add(w FocusableWindow) {
	windows = append(windows, w)
}

func AcceptClientFocus(mode, detail byte) bool {
	return (mode == xproto.NotifyModeNormal || mode == xproto.NotifyModeWhileGrabbed) &&
		(detail == xproto.NotifyDetailVirtual || detail == xproto.NotifyDetailNonlinearVirtual)
}
