package wrap

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadYamlConfig(fp, name string) (vp *viper.Viper, err error) {
	vp = viper.New()
	vp.SetConfigName(name)
	vp.SetConfigFile(fp)

	if err = vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("ReadInConfig: %w", err)
	}

	return vp, nil
}
