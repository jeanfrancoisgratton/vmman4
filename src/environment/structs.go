// vmman4
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename : structs.go
// Original timestamp : 2024/07/29 08:01

package environment

type EnvironmentConfig struct {
	VaultEnabled bool   `json:"VaultEnabled"`
	Host         string `json:"Host"`
	Port         uint8  `json:"Port,omitempty"`
	User         string `json:"User"`
	Passwd       string `json:"Passwd"`
	Database     string `json:"Database,omitempty"`
}

type VaultConfig struct {
	URL             string `json:"URL"`
	User            string `json:"User"`
	Passwd          string `json:"Passwd"`
	SecretEngine    string `json:"SecretEngine"`
	Secret          string `json:"Secret"`
	DbHostKeyName   string `json:"DbHostKeyName"`
	DBPortKeyName   uint8  `json:"DBPortKeyName"`
	DBUserKeyName   string `json:"DBUserKeyName"`
	DBPasswdKeyName string `json:"DBPasswdKeyName"`
	DBNameKeyName   string `json:"DBNameKeyName"`
}
