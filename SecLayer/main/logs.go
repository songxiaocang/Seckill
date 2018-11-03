package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func InitLogger() (err error) {
	config := make(map[string]interface{})
	config["filename"] = secLayerConf.LogPath
	config["level"] = convertLogLevel(secLayerConf.LogLevel)

	configJson, err := json.Marshal(config)
	if err != nil {
		err = fmt.Errorf("config json marshal error: %s", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configJson))
	logs.SetLogFuncCall(true)
	return
}

func convertLogLevel(logLevel string) int {
	switch logLevel {
	case "debug":
		return logs.LevelDebug
	case "info":
		return logs.LevelInfo
	case "warn":
		return logs.LevelWarn
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}
