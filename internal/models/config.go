package models

type Config struct {
	ListenPort string    `yaml:"listen_port"`
	SqlConfig  SqlConfig `yaml:"sql_config"`
	LogLevel   string    `yaml:"log_level"`
}

type SqlConfig struct {
	UsernameDB string `yaml:"username_db"`
	PasswordDB string `yaml:"password_db"`
	HostDB     string `yaml:"host_db"`
	PortDB     string `yaml:"port_db"`
	DBName     string `yaml:"db_name"`
	DelayTime  int    `yaml:"delay_time"`
}
