package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// SigRepositoriesTableName defines
var SigRepositoriesTableName = "sig_repositories"

// SigRepositoriesTableSQL matches with sig_repositories Object
var SigRepositoriesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	name varchar(255) DEFAULT NULL,
	repo_name varchar(255) DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, SigRepositoriesTableName)

// SigRepositories defines
type SigRepositories struct {
	gorm.Model
	Name           string
	RepoName       string
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for SigRepositories
func (srs SigRepositories) GetAdditionalInfo(additionalinfo interface{}) error {
	if srs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(srs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (srs SigRepositories) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(srs)
	if err != nil {
		return "", fmt.Errorf("marshal sig repositories failed. Error: %s", err)
	}
	return string(datas), nil
}
