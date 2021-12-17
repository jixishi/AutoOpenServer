package tools

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	MindustryVersion string
	MindustryTagUrl  string
	WayZerVersion    string
	WayZerTagUrl     string
}

func GetConfig() Config {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			InitConfig()
		} else {
			fmt.Println(err)
			return Config{}
		}
	}
	viper.Unmarshal(&config)
	return config
}
func InitConfig() {
	viper.Set("MindustryVersion", "1")
	viper.Set("MindustryTagUrl", "https://api.github.com/repos/Anuken/Mindustry/releases")
	viper.Set("WayZerVersion", "1")
	viper.Set("WayZerTagUrl", "https://api.github.com/repos/way-zer/ScriptAgent4MindustryExt/releases")
	if err := viper.SafeWriteConfig(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("初始化完成请重启!!!")
}
func SaveConfig(config Config) {
	viper.Set("MindustryVersion", config.MindustryVersion)
	viper.Set("MindustryTagUrl", "https://api.github.com/repos/Anuken/Mindustry/releases")
	viper.Set("WayZerVersion", config.WayZerVersion)
	viper.Set("WayZerTagUrl", "https://api.github.com/repos/way-zer/ScriptAgent4MindustryExt/releases")
	if err := viper.WriteConfig(); err != nil {
		fmt.Println(err)
	}
}
