package config

import "fmt"

type Sqlite struct {
	Filename string `yaml:"filename"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     uint   `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Database struct {
	Postgres *Postgres `yaml:"postgres"`
	Mysql    *Mysql    `yaml:"mysql"`
	Sqlite   *Sqlite   `yaml:"sqlite"`
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		p.Host, p.Port, p.User, p.Password, p.Database)
}

func (m *Mysql) DSN() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?parseTime=true", m.User, m.Password, m.Host, m.Port, m.Database)
}
