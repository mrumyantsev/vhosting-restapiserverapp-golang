package messages

import (
	"app/pkg/download"
	"app/pkg/logger"
)

func ErrorExtensionIsNotMp4() *logger.Log {
	return &logger.Log{StatusCode: 405, ErrCode: 1000, Message: "Extension is not .mp4", ErrLevel: logger.ErrLevelError}
}

func InfoPutDownloadLink(dload *download.Download) *logger.Log {
	return &logger.Log{StatusCode: 200, Message: dload}
}
