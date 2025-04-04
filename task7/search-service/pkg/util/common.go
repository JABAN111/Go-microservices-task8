package util

import (
	"io"

	log "github.com/sirupsen/logrus"
)

func SafeClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Error("Failed to close connection", "error", err)
	}
}
