package configurator

type SqlConfig struct {
	DbDriver        string
	DbHost          string
	DbPort          string
	DbName          string
	DbUsername      string
	DbPassword      string
	DbSslMode       string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int64
}
