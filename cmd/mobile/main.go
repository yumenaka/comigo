package server

import (
	"github.com/yumenaka/comigo/config"
	"strconv"
	// _ "golang.org/x/mobile/bind"
)

func Start(path string) (string, error) {
	config.SetOpenBrowser(false)
	return strconv.Itoa(config.GetCfg().Port), nil
}
