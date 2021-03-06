package windowmanager

import (
	"log"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xinerama"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/janbina/swm/internal/groupmanager"
	"github.com/janbina/swm/internal/heads"
)

// Root window configure request
func configureRequestFun(x *xgbutil.XUtil, e xevent.ConfigureRequestEvent) {
	log.Printf("Configure request: %s", e)
	if _, ok := managedWindows[e.Window]; ok {
		return
	}

	xwindow.New(x, e.Window).Configure(
		int(e.ValueMask),
		int(e.X),
		int(e.Y),
		int(e.Width),
		int(e.Height),
		e.Sibling,
		e.StackMode,
	)
}

func mapRequestFun(_ *xgbutil.XUtil, e xevent.MapRequestEvent) {
	log.Printf("Map request: %s", e)
	manageWindow(e.Window)
}

func applyStruts() {
	rootG := xwindow.RootGeometry(X)
	wh := make(xinerama.Heads, len(heads.Heads)+1)
	wh[0] = xrect.New(rootG.Pieces())
	for i, head := range heads.Heads {
		wh[i+1] = xrect.New(head.Pieces())
	}

	for w := range managedWindows {
		strut, _ := ewmh.WmStrutPartialGet(X, w)
		if strut == nil {
			continue
		}
		xrect.ApplyStrut(wh, uint(rootG.Width()), uint(rootG.Height()),
			strut.Left, strut.Right, strut.Top, strut.Bottom,
			strut.LeftStartY, strut.LeftEndY,
			strut.RightStartY, strut.RightEndY,
			strut.TopStartX, strut.TopEndX,
			strut.BottomStartX, strut.BottomEndX,
		)
	}

	RootGeometryStruts = wh[0]
	heads.HeadsStruts = wh[1:]

	setWorkArea(groupmanager.GetNumGroups())

	for _, win := range managedWindows {
		win.RootGeometryChanged()
	}
}
