package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// WebView2 normally uses the executable name as its profile directory.
// Reuse the previous KO profile or the original OpenSave profile so a
// renamed portable executable can read its existing UI preferences.
func webviewUserDataPath() string {
	if runtime.GOOS != "windows" {
		return ""
	}
	root, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	for _, name := range []string{"OpenSave-KO.exe", "OpenSave.exe"} {
		path := filepath.Join(root, name)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path
		}
	}
	return filepath.Join(root, "OpenSave-KO.exe")
}
