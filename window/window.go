package window

import (
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type Window struct {
	win *xwindow.Window
}

func New(x *xgbutil.XUtil, w xproto.Window) *Window {
	return &Window{
		win: xwindow.New(x, w),
	}
}

func (w *Window) Map() {
	w.win.Map()
}

func (w *Window) Focus() {
	w.win.Focus()
}


