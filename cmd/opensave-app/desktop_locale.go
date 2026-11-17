package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed frontend/src/locales/en.json
var desktopEnglish []byte

//go:embed frontend/src/locales/ko.json
var desktopKorean []byte

type trayCatalog map[string]json.RawMessage

func readTrayCatalog(data []byte) trayCatalog {
	var catalog trayCatalog
	if err := json.Unmarshal(data, &catalog); err != nil {
		panic(err)
	}
	return catalog
}

var englishTray = readTrayCatalog(desktopEnglish)
var koreanTray = readTrayCatalog(desktopKorean)

func (c trayCatalog) label(key string, values ...any) string {
	var text string
	if err := json.Unmarshal(c["tray."+key], &text); err != nil {
		panic(err)
	}
	replacements := make([]string, 0, len(values))
	for i := 0; i < len(values); i += 2 {
		replacements = append(replacements, "{"+values[i].(string)+"}", fmt.Sprint(values[i+1]))
	}
	return strings.NewReplacer(replacements...).Replace(text)
}

// SetDesktopLocale applies the WebView's resolved language to native UI.
func (a *App) SetDesktopLocale(locale string) {
	a.koreanDesktop.Store(locale == "ko")
	select {
	case a.trayRefresh <- struct{}{}:
	default:
	}
}

func (a *App) trayCatalog() trayCatalog {
	if a.koreanDesktop.Load() {
		return koreanTray
	}
	return englishTray
}
