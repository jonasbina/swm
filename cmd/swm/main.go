package main

import (
	"flag"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/janbina/swm/internal/communication"
	"github.com/janbina/swm/internal/config"
	"github.com/janbina/swm/internal/windowmanager"
	"log"
	"os/exec"
)

func main() {

	replace := flag.Bool("replace", false, "whether swm should replace current wm")
	customConfig := flag.String("c", "", "path to swmrc file")
	flag.Parse()

	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatalf("Cannot initialize x connection: %s", err)
	}
	defer X.Conn().Close()

	if err := windowmanager.Initialize(X, *replace); err != nil {
		log.Fatalf("Cannot initialize window manager: %s", err)
	}

	if err := windowmanager.SetupRoot(); err != nil {
		log.Fatalf("Cannot setup root window: %s", err)
	}

	go communication.Listen(X.Conn())

	config.FindAndRunSwmrc(*customConfig)

	windowmanager.ManageExistingClients()

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, ev xevent.KeyPressEvent) {
			exec.Command("xterm").Start()
		},
	).Connect(windowmanager.X, windowmanager.Root.Id, "control-Mod1-return", true)

	windowmanager.Run()
}
