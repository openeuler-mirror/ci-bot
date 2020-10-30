package config

type Config struct {
	GiteeToken               string             `yaml:"giteeToken" envVariable:"GITEE_TOKEN"`
	WebhookSecret            string             `yaml:"webhookSecret" envVariable:"WEBHOOK_SECRET"`
	DataBaseType             string             `yaml:"databaseType"`
	DataBaseHost             string             `yaml:"databaseHost" envVariable:"DATABASE_HOST"`
	DataBasePort             int                `yaml:"databasePort" envVariable:"DATABASE_PORT"`
	DataBaseName             string             `yaml:"databaseName"`
	DataBaseUserName         string             `yaml:"databaseUserName" envVariable:"DATABASE_USERNAME"`
	DataBasePassword         string             `yaml:"databasePassword" envVariable:"DATABASE_PASSWORD"`
	PrUpdateLabelFlag        string             `yaml:"prUpdateLabelFlag"`
	DelLabels                []string           `yaml:"delLabels"`
	WatchProjectFiles        []WatchProjectFile `yaml:"watchProjectFiles"`
	WatchProjectFileDuration int                `yaml:"watchProjectFileDuration"`
	WatchSigFiles            []WatchSigFile     `yaml:"watchSigFiles"`
	WatchSigFileDuration     int                `yaml:"watchSigFileDuration"`
	WatchOwnerFiles          []WatchOwnerFile   `yaml:"watchOwnerFiles"`
	WatchOwnerFileDuration   int                `yaml:"watchOwnerFileDuration"`
	WatchFrozenFile          []WatchFrozenFile  `yaml:"watchFrozenFile"`
	WatchFrozenDuration      int                `yaml:"watchFrozenDuration"`
	BotName                  string             `yaml:"botName"`
	CommunityName            string             `yaml:"communityName"`
	ClaLink                  string             `yaml:"claLink"`
	CommandLink              string             `yaml:"commandLink"`
	ContactEmail             string             `yaml:"contactEmail"`
	LgtmCountsRequired       int                `yaml:"lgtmCountsRequired"`
	AccordingFile            string             `yaml:"accordingfile"`
	NewFileOwner             string             `yaml:"newfileowner"`
	NewFileRepo              string             `yaml:"newfilerepo"`
	NewFileBranch            string             `yaml:"newfilebranch"`
	ServiceFile              string             `yaml:"tmpservicefile"`
	ServicePath              string             `yaml:"tmpservicepath"`
	GuideURL                 string             `yaml:"guideurl"`
	CommitsThreshold         int                `yaml:"commitsThreshold"`
	SquashCommitLabel        string             `yaml:"squashCommitLabel"`
	RequiringLabels          []string           `yaml:"requiringLabels"`
	MissingLabels            []string           `yaml:"missingLabels"`
	AutoDetectCla            bool               `yaml:"autoDetectCla"`
}

type WatchProjectFile struct {
	WatchProjectFileOwner string `yaml:"watchProjectFileOwner"`
	WatchprojectFileRepo  string `yaml:"watchprojectFileRepo"`
	WatchprojectFilePath  string `yaml:"watchprojectFilePath"`
	WatchProjectFileRef   string `yaml:"watchProjectFileRef"`
}

type WatchSigFile struct {
	WatchSigFileOwner string `yaml:"watchSigFileOwner"`
	WatchSigFileRepo  string `yaml:"watchSigFileRepo"`
	WatchSigFilePath  string `yaml:"watchSigFilePath"`
	WatchSigFileRef   string `yaml:"watchSigFileRef"`
}

type WatchOwnerFile struct {
	WatchOwnerFileOwner string `yaml:"watchOwnerFileOwner"`
	WatchOwnerFileRepo  string `yaml:"watchOwnerFileRepo"`
	WatchOwnerFilePath  string `yaml:"watchOwnerFilePath"`
	WatchOwnerFileRef   string `yaml:"watchOwnerFileRef"`
}

type WatchFrozenFile struct {
	FrozenFileOwner string `yaml:"frozenFileOwner"`
	FrozenFileRepo  string `yaml:"frozenFileRepo"`
	FrozenFilePath  string `yaml:"frozenFilePath"`
	FrozenFileRef   string `yaml:"frozenFileRef"`
}
