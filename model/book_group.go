package model

import "github.com/yumenaka/comigo/tools/logger"

// GenerateBookGroup 重新生成书组
func GenerateBookGroup() {
	if err := IStore.GenerateBookGroup(); err != nil {
		logger.Infof("%s", err)
	}
}
