package model

import "fmt"

type Config struct {
	MySQL      *MySQLOptions      `yaml:"MYSQL"`
	PostgreSQL *PostgreSQLOptions `yaml:"POSTGRESQL"`
}

type MySQLOptions struct {
	Host    string `yaml:"HOST"`
	Port    string `yaml:"PORT"`
	User    string `yaml:"USER"`
	Pass    string `yaml:"PASS"`
	Db      string `yaml:"DB"`
	Charset string `yaml:"CHARSET"`
}

func (options *MySQLOptions) String() string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s",
		options.User, options.Pass,
		options.Host, options.Port,
		options.Db, options.Charset,
	)
}

type PostgreSQLOptions struct {
	Host string `yaml:"HOST"`
	Port string `yaml:"PORT"`
	User string `yaml:"USER"`
	Pass string `yaml:"PASS"`
	Db   string `yaml:"DB"`
}

func (options *PostgreSQLOptions) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		options.Host, options.Port,
		options.User, options.Db,
		options.Pass,
	)
}
