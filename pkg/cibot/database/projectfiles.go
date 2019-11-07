package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// ProjectFilesTableName defines
var ProjectFilesTableName = "project_files"

// ProjectFilesTableSQL matches with ProjectFiles Object
var ProjectFilesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	owner varchar(255) DEFAULT NULL,
	repo varchar(255) DEFAULT NULL,
	path varchar(255) DEFAULT NULL,
	ref varchar(255) DEFAULT NULL,
	current_sha varchar(255) DEFAULT NULL,
	target_sha varchar(255) DEFAULT NULL,
	waiting_sha varchar(255) DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, ProjectFilesTableName)

// ProjectFiles defines
type ProjectFiles struct {
	gorm.Model
	Owner          string
	Repo           string
	Path           string
	Ref            string
	CurrentSha     string
	TargetSha      string
	WaitingSha     string
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for ProjectFiles
func (pfs ProjectFiles) GetAdditionalInfo(additionalinfo interface{}) error {
	if pfs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(pfs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (pfs ProjectFiles) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(pfs)
	if err != nil {
		return "", fmt.Errorf("marshal project files failed. Error: %s", err)
	}
	return string(datas), nil
}
