package install

import (
	"fmt"
	"net/url"
	"strings"

	"xorm.io/xorm/schemas"
)

// CheckConfigFileResp check config file if exist or not response
type CheckConfigFileResp struct {
	ConfigFileExist     bool `json:"config_file_exist"`
	DBConnectionSuccess bool `json:"db_connection_success"`
	DbTableExist        bool `json:"db_table_exist"`
}

// CheckDatabaseReq check database
type CheckDatabaseReq struct {
	DbType     string `validate:"required,oneof=postgres sqlite3 mysql" json:"db_type"`
	DbUsername string `json:"db_username"`
	DbPassword string `json:"db_password"`
	DbHost     string `json:"db_host"`
	DbName     string `json:"db_name"`
	DbFile     string `json:"db_file"`
}

// GetConnection get connection string
func (r *CheckDatabaseReq) GetConnection() string {
	if r.DbType == string(schemas.SQLITE) {
		return r.DbFile
	}
	if r.DbType == string(schemas.MYSQL) {
		return fmt.Sprintf("%s:%s@tcp(%s)/%s",
			r.DbUsername, r.DbPassword, r.DbHost, r.DbName)
	}
	if r.DbType == string(schemas.POSTGRES) {
		host, port := parsePgSQLHostPort(r.DbHost)
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, r.DbUsername, r.DbPassword, r.DbName)
	}
	return ""
}

func parsePgSQLHostPort(dbHost string) (host string, port string) {
	if strings.Contains(dbHost, ":") {
		idx := strings.LastIndex(dbHost, ":")
		host, port = dbHost[:idx], dbHost[idx+1:]
	} else if len(dbHost) > 0 {
		host = dbHost
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "5432"
	}
	return host, port
}

// CheckDatabaseResp check database response
type CheckDatabaseResp struct {
	ConnectionSuccess bool `json:"connection_success"`
}

// InitEnvironmentResp init environment response
type InitEnvironmentResp struct {
	Success            bool   `json:"success"`
	CreateConfigFailed bool   `json:"create_config_failed"`
	DefaultConfig      string `json:"default_config"`
	ErrType            string `json:"err_type"`
}

// InitBaseInfoReq init base info request
type InitBaseInfoReq struct {
	Language      string `validate:"required,gt=0,lte=30" json:"lang"`
	SiteName      string `validate:"required,gt=0,lte=30" json:"site_name"`
	SiteURL       string `validate:"required,gt=0,lte=512,url" json:"site_url"`
	ContactEmail  string `validate:"required,email,gt=0,lte=500" json:"contact_email"`
	AdminName     string `validate:"required,gt=4,lte=30" json:"name"`
	AdminPassword string `validate:"required,gte=8,lte=32" json:"password"`
	AdminEmail    string `validate:"required,email,gt=0,lte=500" json:"email"`
}

func (r *InitBaseInfoReq) FormatSiteUrl() {
	parsedUrl, err := url.Parse(r.SiteURL)
	if err != nil {
		return
	}
	r.SiteURL = fmt.Sprintf("%s://%s", parsedUrl.Scheme, parsedUrl.Host)
}
