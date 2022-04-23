package forkjoin

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewLogger() *StandardLogger {
	baseLogger := logrus.New()
	baseLogger.SetOutput(os.Stdout)
	baseLogger.SetLevel(logrus.InfoLevel)
	var sl = &StandardLogger{baseLogger}
	sl.Formatter = &logrus.JSONFormatter{}
	return sl
}

var (
	internalError   = LogEvent{InternalError, "Request Id Message Id: Internal Error: %s %s %s"}
	requestError    = LogEvent{RequestError, "Request Id Message Id: Invalid Request: %s %s %v"}
	abortRequest    = LogEvent{RequestAborted, "Request Id Message Id: Request Aborted: %s %s %v"}
	responseError   = LogEvent{ResponseError, "Request Id Message Id: Response error: %s %s %v"}
	connectionError = LogEvent{ConnectionError, "Request Id Message Id: Connection error: %s %s %v"}
	info            = LogEvent{Info, "Request Id %s Message Id: %s: %s \n"}
)

func (l *StandardLogger) LogAbortedRequest() {

}

func (l *StandardLogger) LogInfo(requestID, messageID, message string) {
	l.Infof(info.message, requestID, messageID, message)
}
