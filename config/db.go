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

type Database struct {
	Sqlite   *Sqlite   `yaml:"sqlite"`
	Postgres *Postgres `yaml:"postgres"`
}

func (p *Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		p.Host, p.Port, p.User, p.Password, p.Database)
}
