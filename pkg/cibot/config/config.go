package config

type Config struct {
	GiteeToken               string             `yaml:"giteeToken"`
	WebhookSecret            string             `yaml:"webhookSecret"`
	DataBaseType             string             `yaml:"databaseType"`
	DataBaseHost             string             `yaml:"databaseHost"`
	DataBasePort             int                `yaml:"databasePort"`
	DataBaseName             string             `yaml:"databaseName"`
	DataBaseUserName         string             `yaml:"databaseUserName"`
	DataBasePassword         string             `yaml:"databasePassword"`
	WatchProjectFiles        []WatchProjectFile `yaml:"watchProjectFiles"`
	WatchProjectFileDuration int                `yaml:"watchProjectFileDuration"`
}

type WatchProjectFile struct {
	WatchProjectFileOwner string `yaml:"watchProjectFileOwner"`
	WatchprojectFileRepo  string `yaml:"watchprojectFileRepo"`
	WatchprojectFilePath  string `yaml:"watchprojectFilePath"`
	WatchProjectFileRef   string `yaml:"watchProjectFileRef"`
}
