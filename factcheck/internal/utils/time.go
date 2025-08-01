package utils

import (
	"log/slog"
	"time"
)

var now func() time.Time

func init() {
	slog.Info("testhelper.init.now") //nolint:noctx
	now = time.Now
}

// TimeNow is a stub for time.TimeNow. Use this function in code,
// so that we can have determinism in our tests with TimeFreeze/TimeUnfreeze.
func TimeNow() time.Time {
	return now()
}

func TimeSince(since time.Time) time.Duration {
	return TimeNow().Sub(since)
}

func TimeFreeze(t time.Time) {
	now = func() time.Time {
		return t
	}
	slog.Debug("testhelper.time.freeze") //nolint:noctx
}

func TimeUnfreeze() {
	now = time.Now
	slog.Debug("testhelper.time.unfreeze") //nolint:noctx
}
