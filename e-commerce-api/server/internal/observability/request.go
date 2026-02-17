package observability

import "time"

func LogRequest(
	logger Logger,
	method string,
	path string,
	status int,
	duration time.Duration,
	clientIP string,
) {
	logger.Info(
		"%s %s %d - %v - %s",
		method,
		path,
		status,
		duration,
		clientIP,
	)
}
