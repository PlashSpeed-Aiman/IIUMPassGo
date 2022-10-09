package dialog

import "github.com/sqweek/dialog"

func XPlatMessageBox(title string, info string) {
	dialog.Message(info).Title(title).Info()
}
