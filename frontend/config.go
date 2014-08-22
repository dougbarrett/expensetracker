package main

type tomlConfig struct {
	Mysql     tomlMysqlConfig
	WebServer tomlWebConfig
	Session   tomlSession
}

type tomlMysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

type tomlWebConfig struct {
	IsFastCGI  bool
	Port       string
	PublicPort string
}

type tomlSession struct {
	Name string
	Key  string
}
