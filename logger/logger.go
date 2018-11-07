package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
)

func InitLogger() (logger *zap.Logger, err error) {
	file, err := os.Open("logger/loggerConfig.json")
	if err != nil {
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	configJSON := make([]byte, stat.Size())
	_, err = file.Read(configJSON)
	if err != nil {
		return
	}

	var cfg zap.Config
	if err = json.Unmarshal(configJSON, &cfg); err != nil {
		return
	}
	logger, err = cfg.Build()
	if err != nil {
		return
	}

	fmt.Println(time.Now().UTC())
	//logger = logger.With(
	//	zap.Time("time", time.Now()),
	//)
	zap.ReplaceGlobals(logger)
	return
}
