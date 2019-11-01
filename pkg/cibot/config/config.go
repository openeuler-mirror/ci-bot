package config

type Config struct {
	Owner                 string `yaml:"owner"`
	Repo                  string `yaml:"repository"`
	GiteeToken            string `yaml:"giteeToken"`
	WebhookSecret         string `yaml:"webhookSecret"`
	DataBaseType          string `yaml:"databaseType"`
	DataBaseHost          string `yaml:"databaseHost"`
	DataBasePort          int    `yaml:"databasePort"`
	DataBaseName          string `yaml:"databaseName"`
	DataBaseUserName      string `yaml:"databaseUserName"`
	DataBasePassword      string `yaml:"databasePassword"`
	WatchProjectFileOwner string `yaml:"watchProjectFileOwner"`
	WatchprojectFileRepo  string `yaml:"watchprojectFileRepo"`
	WatchprojectFilePath  string `yaml:"watchprojectFilePath"`
	WatchProjectFileRef   string `yaml:"watchProjectFileRef"`
}
