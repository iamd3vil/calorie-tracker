package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type cfgTelegram struct {
	ApiKey string `koanf:"api_key"`
}

type cfgDb struct {
	SqlitePath string `koanf:"sqlite_path"`
}

type Config struct {
	Telegram cfgTelegram
	Db       cfgDb
}

var cfg Config

func initConfig() {
	var k = koanf.New(".")

	// Load TOML config.
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Unmarshal("telegram", &cfg.Telegram)
	k.Unmarshal("sqlite", &cfg.Db)
}
