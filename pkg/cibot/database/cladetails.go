package database

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
)

// CLADetailsTableName defines
var CLADetailsTableName = "cla_details"

// CLADetailsTableSQL matches with CLADetails Object
var CLADetailsTableSQL = fmt.Sprintf(`CREATE TABLE %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT,
	created_at timestamp NULL DEFAULT NULL,
	updated_at timestamp NULL DEFAULT NULL,
	deleted_at timestamp NULL DEFAULT NULL,
	type int(10) unsigned DEFAULT NULL,
	name varchar(255) DEFAULT NULL,
	title varchar(255) DEFAULT NULL,
	corporation varchar(255) DEFAULT NULL,
	address varchar(255) DEFAULT NULL,
	date varchar(255) DEFAULT NULL,
	email varchar(255) DEFAULT NULL,
	telephone varchar(255) DEFAULT NULL,
	fax varchar(255) DEFAULT NULL,
	additional_info text,
	PRIMARY KEY (id)
  ) ENGINE=InnoDB DEFAULT CHARSET=utf8`, CLADetailsTableName)

// CLADetails defines
type CLADetails struct {
	gorm.Model
	Type           int
	Name           string
	Title          string
	Corporation    string
	Address        string
	Date           string
	Email          string
	Telephone      string
	Fax            string
	AdditionalInfo string `sql:"type:text"`
}

// GetAdditionalInfo for CLADetails
func (cds CLADetails) GetAdditionalInfo(additionalinfo interface{}) error {
	if cds.AdditionalInfo != "" {
		err := json.Unmarshal([]byte(cds.AdditionalInfo), &additionalinfo)
		if err != nil {
			return err
		}
	}
	return nil
}
