package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// SigRecordsTableName defines
var SigRecordsTableName = "sig_records"

// SigRecordsTableSQL matches with sig_records Object
var SigRecordsTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	name varchar(255) DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, SigRecordsTableName)

// SigRecords defines
type SigRecords struct {
	gorm.Model
	Name           string
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for SigRecords
func (srs SigRecords) GetAdditionalInfo(additionalinfo interface{}) error {
	if srs.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(srs.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ToString for convert
func (srs SigRecords) ToString() (string, error) {
	// Marshal datas
	datas, err := json.Marshal(srs)
	if err != nil {
		return "", fmt.Errorf("marshal sig records failed. Error: %s", err)
	}
	return string(datas), nil
}
