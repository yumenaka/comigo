//go:build !linux

package autostart

import "errors"

// State 非 systemd 平台明确报告不支持，保持阅读与桌面功能可用。
type State struct {
	Enabled   bool   `json:"enabled"`
	Supported bool   `json:"supported"`
	Scope     string `json:"scope"`
}

func Status() (State, error) { return State{Scope: "user"}, nil }
func Set(bool) error         { return errors.New("automatic startup currently requires systemd") }
