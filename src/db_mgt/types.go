// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: src/db_mgt/types.go
// Original timestamp : 2026.10.02 15:18:19

package db_mgt

// This structure holds the basic DB config

type DBenvironmentStruct struct {
	Host     string `json:"dbhost"`
	Port     uint   `json:"dbport,omitempty"`
	User     string `json:"dbuser"`
	Password string `json:"dbpassword"`
}
