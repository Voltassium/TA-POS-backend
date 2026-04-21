package config

type Application struct {
	Port        string                 `mapstructure:"port"`
	Environment ApplicationEnvironment `mapstructure:"environment"`
}
