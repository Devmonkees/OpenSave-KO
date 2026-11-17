//go:build windows

package main

import (
	"net/url"
	"os"
	"path/filepath"

	toast "git.sr.ht/~jackmordaunt/go-toast/v2"
	"github.com/google/uuid"
)

// Windows uses the app ID as the toast header. Give the KO build and test
// builds their own activation classes to keep clicks routed to the right EXE.
func toastIdentity() (appID, guid string) {
	appID = displayName
	if instanceID != defaultInstanceID {
		appID += " (" + instanceID + ")"
	}
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte("opensave-toast:"+appID))
	return appID, "{" + id.String() + "}"
}

func (a *App) initNotifications() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	// The header's icon: the app's own, as a file Windows can read.
	icon := filepath.Join(os.TempDir(), "OpenSave", "opensave-notification-icon.png")
	if os.MkdirAll(filepath.Dir(icon), 0o755) != nil || os.WriteFile(icon, appIcon, 0o644) != nil {
		icon = ""
	}
	appID, guid := toastIdentity()
	if err := toast.SetAppData(toast.AppData{AppID: appID, GUID: guid, ActivationExe: exe, IconPath: icon}); err != nil {
		if a.daemon != nil {
			a.daemon.Log.Log("warn", "desktop notifications are unavailable: "+err.Error())
		}
		return
	}
	toast.SetActivationCallback(func(args string, _ []toast.UserData) {
		a.openFromNotification(args)
	})
}

func (a *App) showNote(n DesktopNote) error {
	appID, _ := toastIdentity()
	t := toast.Notification{
		AppID:               appID,
		Title:               plainText(n.Title),
		Body:                plainText(n.Body),
		ActivationType:      toast.Foreground,
		ActivationArguments: encodeOpen(n.Open),
		// Silent: what needs you has OpenSave's own chime already, and what
		// merely happened should not make a sound over a game.
		Audio: toast.Silent,
	}
	if img := a.noteImage(n.Image); img != "" {
		t.HeroIcon = fileURI(img)
	}
	return t.Push()
}

// fileURI is how a toast is given a picture on disk: a file:/// URI, as
// Microsoft's own desktop examples build it, with spaces and the like escaped.
func fileURI(path string) string {
	return (&url.URL{Scheme: "file", Path: "/" + filepath.ToSlash(path)}).String()
}
