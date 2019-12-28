package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// SigFilesTableName defines
var SigFilesTableName = "sig_files"

// SigFilesTableSQL matches with SigFiles Object
var SigFilesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
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
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, SigFilesTableName)

// SigFiles defines
type SigFiles struct {
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

// GetAdditionalInfo for SigFiles
func (sfs SigFiles) GetAdditionalInfo(additionalinfo interface{}) error {
	if sfs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(sfs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (sfs SigFiles) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(sfs)
	if err != nil {
		return "", fmt.Errorf("marshal sig files failed. Error: %s", err)
	}
	return string(datas), nil
}
