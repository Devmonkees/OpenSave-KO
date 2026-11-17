package main

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/opensave/opensave/internal/logging"
	"github.com/opensave/opensave/internal/syncpause"
)

// What the tray menu says, worked out apart from the tray itself so it can be
// tested: the tray library runs its own message loop and cannot be driven
// from a test.

// trayFacts is what the tray's first line is made from.
type trayFacts struct {
	Games     int
	Conflicts int
	Pause     syncpause.Status
	// Syncing names the games moving between devices right now.
	Syncing []string
}

// trayStatus is the tray's first line: the one thing most worth knowing
// without opening the window. A decision waiting on you comes first, then a
// pause (which you might have forgotten), then what is moving.
func trayStatus(f trayFacts, catalog trayCatalog) string {
	switch {
	case f.Conflicts == 1:
		return catalog.label("conflictOne")
	case f.Conflicts > 1:
		return catalog.label("conflictMany", "n", f.Conflicts)
	case f.Pause.Paused && f.Pause.UntilRestart:
		return catalog.label("paused")
	case f.Pause.Paused:
		return catalog.label("pausedUntil", "time", minutesLeft(f.Pause.RemainingSeconds, catalog))
	case len(f.Syncing) == 1:
		return catalog.label("syncOne", "game", f.Syncing[0])
	case len(f.Syncing) > 1:
		return catalog.label("syncMany", "n", len(f.Syncing))
	case f.Games == 0:
		return catalog.label("empty")
	case f.Games == 1:
		return catalog.label("watchOne")
	default:
		return catalog.label("watchMany", "n", f.Games)
	}
}

func minutesLeft(seconds int64, catalog trayCatalog) string {
	m := int64((time.Duration(seconds)*time.Second + time.Minute - 1) / time.Minute)
	switch {
	case m >= 60 && m%60 == 0:
		return catalog.label("hours", "h", m/60)
	case m >= 60:
		return catalog.label("hoursMinutes", "h", m/60, "m", m%60)
	case m <= 1:
		return catalog.label("underMinute")
	default:
		return catalog.label("minutes", "m", m)
	}
}

// trayRecentMax is how long a line may be in the tray's recent activity; a
// menu does not wrap, and a path runs off the edge of the screen.
const trayRecentMax = 64

// trayRecent picks the last few things worth a glance from the activity log,
// newest first: what went through and what went wrong, before the routine.
// Falls back to the routine when nothing else has happened.
func trayRecent(entries []logging.Entry, n int) []string {
	var picked, routine []string
	for i := len(entries) - 1; i >= 0 && len(picked) < n; i-- {
		e := entries[i]
		line := shorten(e.Message)
		switch e.Level {
		case "success", "warn", "error":
			if e.Level != "success" {
				line = "⚠ " + line
			}
			picked = append(picked, line)
		default:
			if len(routine) < n {
				routine = append(routine, line)
			}
		}
	}
	if len(picked) == 0 {
		return routine
	}
	return picked
}

func shorten(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if utf8.RuneCountInString(s) <= trayRecentMax {
		return s
	}
	r := []rune(s)
	return string(r[:trayRecentMax-1]) + "…"
}
