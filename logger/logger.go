package logger

import (
	"os"
	"github.com/gofiber/fiber/v3/log"
)

func Logger()*os.File {
	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("Failed to open log file: ", err)
	}
	return f
}