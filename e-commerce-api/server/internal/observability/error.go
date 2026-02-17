package observability

import (
	"runtime"
)

func LogError(logger Logger, err error, context string) {
	if err == nil {
		return
	}

	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)

	logger.Error(
		"%s: %v\nstack=%s",
		context,
		err,
		string(buf[:n]),
	)
}
