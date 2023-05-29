package settings

import (
	"bytes"
	"fmt"

	"github.com/d2jvkpn/gotk"
	"github.com/d2jvkpn/gotk/impls"
	"github.com/spf13/viper"
)

var (
	_ServiceName string
	_Project     *viper.Viper
	_Config      *viper.Viper
	Logger       *impls.Logger
	Meta         map[string]any
)

func SetProject(bts []byte) (err error) {
	_Project = viper.New()
	_Project.SetConfigType("yaml")

	// _Project.ReadConfig(strings.NewReader(str))
	if err = _Project.ReadConfig(bytes.NewReader(bts)); err != nil {
		return err
	}

	Meta = gotk.BuildInfo()
	Meta["project"] = _Project.GetString("project")
	Meta["version"] = _Project.GetString("version")

	return nil
}

func SetConfig(vp *viper.Viper) (err error) {
	_Config = vp

	if _ServiceName = _Config.GetString("service_name"); _ServiceName == "" {
		return fmt.Errorf("service_name is empty in config")
	}

	return nil
}

func ConfigSub(key string) *viper.Viper {
	return _Config.Sub(key)
}

func Project() string {
	return _Project.GetString("project")
}

func Version() string {
	return _Project.GetString("version")
}

func DemoConfig() string {
	return _Project.GetString("config")
}

func DotEnv() string {
	return _Project.GetString("dotenv")
}
