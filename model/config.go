package model

type Config struct {
	App struct {
		GeminiKey  []string `yaml:"GeminiKey"`
		GeminiUrl  string   `yaml:"GeminiUrl"`
		UserAgents []string `yaml:"UserAgents"`
	} `yaml:"App"`
	Server struct {
		Host        string `yaml:"Host"`
		Port        string `yaml:"Port"`
		MaxFileSize int    `yaml:"MaxFileSize"`
	} `yaml:"Server"`
	Proxy struct {
		Protocol string `yaml:"Protocol"`
	} `yaml:"Proxy"`
}
