package config

import (
	"path/filepath"
	"runtime"
)

// ConfigFileInfo 描述运行配置的文件来源；文件删除后，已加载配置仍保留在内存中。
type ConfigFileInfo struct {
	Path     string `json:"path"`
	Location string `json:"location"`
	Type     string `json:"type"`
	Format   string `json:"format"`
	Exists   bool   `json:"exists"`
}

type Status struct {
	In   string
	Path struct {
		WorkingDirectory string
		HomeDirectory    string
		ProgramDirectory string
	}
	Current ConfigFileInfo `json:"current"`
}

// SetConfigStatus 区分实际加载的文件与各保存位置的现存文件。
func (c *Status) SetConfigStatus() error {
	ConfigFileLock.Lock()
	defer ConfigFileLock.Unlock()
	*c = Status{In: "None", Current: ConfigFileInfo{Location: "None", Type: configProfile(), Format: "toml"}}
	if runtime.GOOS == "js" || cfg.TemporaryReaderMode {
		return nil
	}
	c.Path.WorkingDirectory = GetWorkingDirectoryConfig()
	c.Path.HomeDirectory = GetHomeDirectoryConfig()
	c.Path.ProgramDirectory = GetProgramDirectoryConfig()
	if cfg.ConfigFile == "" {
		return nil
	}
	current, err := filepath.Abs(cfg.ConfigFile)
	if err != nil {
		return err
	}
	c.Current.Path = current
	c.Current.Exists = fileExists(current)
	c.Current.Location = "Custom"
	for _, location := range configSearchLocations() {
		directory, err := filepath.Abs(location.dir)
		if err != nil {
			return err
		}
		if filepath.Dir(current) == directory {
			c.Current.Location = location.name
			break
		}
	}
	c.In = c.Current.Location
	return nil
}
