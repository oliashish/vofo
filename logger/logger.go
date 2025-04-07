package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

func Logger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Error Initializing logger : ", err)
		os.Exit(1)
	}
	return logger
}
