package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// BranchesTableName defines
var BranchesTableName = "branches"

// BranchesTableSQL matches with Branches Object
var BranchesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	owner varchar(255) DEFAULT NULL,
	repo varchar(255) DEFAULT NULL,
	name varchar(255) DEFAULT NULL,
	type varchar(255) DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, BranchesTableName)

// Branches defines
type Branches struct {
	gorm.Model
	Owner string
	Repo  string
	Name  string
	// "protected" or "readonly", only protected is supported yet
	Type           string
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for Branches
func (bs Branches) GetAdditionalInfo(additionalinfo interface{}) error {
	if bs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(bs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (bs Branches) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(bs)
	if err != nil {
		return "", fmt.Errorf("marshal branches failed. Error: %s", err)
	}
	return string(datas), nil
}
