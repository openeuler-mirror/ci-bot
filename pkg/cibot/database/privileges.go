package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// PrivilegesTableName defines
var PrivilegesTableName = "privileges"

// PrivilegesTableSQL matches with Privileges Object
var PrivilegesTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	user varchar(255) DEFAULT NULL,
	type varchar(255) DEFAULT NULL,
	repository_id int(10) unsigned DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, PrivilegesTableName)

// Privileges defines
type Privileges struct {
	gorm.Model
	User           string
	Type           string
	RepositoryID   uint
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for Privileges
func (ps Privileges) GetAdditionalInfo(additionalinfo interface{}) error {
	if ps.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(ps.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (ps Privileges) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(ps)
	if err != nil {
		return "", fmt.Errorf("marshal privileges failed. Error: %s", err)
	}
	return string(datas), nil
}
