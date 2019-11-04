package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// RepositoriesTableName defines
var RepositoriesTableName = "repositories"

// RepositoriesTableSQL matches with Repositories Object
var RepositoriesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	owner varchar(255) DEFAULT NULL,
	repo varchar(255) DEFAULT NULL,
	description varchar(255) DEFAULT NULL,
	type varchar(255) DEFAULT NULL,
	project_file_id int(10) unsigned DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, RepositoriesTableName)

// Repositories defines
type Repositories struct {
	gorm.Model
	Owner          string
	Repo           string
	Description    string
	Type           string
	ProjectFileID  uint
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for Repositories
func (rs Repositories) GetAdditionalInfo(additionalinfo interface{}) error {
	if rs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(rs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (rs Repositories) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(rs)
	if err != nil {
		return "", fmt.Errorf("marshal repositories failed. Error: %s", err)
	}
	return string(datas), nil
}
