package model

import "fmt"

type Config struct {
	MySQL *MySQLOptions `yaml:"MYSQL"`
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
