package config

import "github.com/BurntSushi/toml"

type Global struct {
	DSPD  Dspd  `toml:"dspd"`
	Textd Textd `toml:"textd"`
	Timed Timed `toml:"timed"`
	Bot   Bot   `toml:"bot"`
}

type Dspd struct {
	GRPCServerAddress string `toml:"grpc_address"`
}

type Textd struct {
	GRPCServerAddress string `toml:"grpc_address"`
	FontPath          string `toml:"font_path"`
}

type Timed struct {
	GRPCServerAddress string `toml:"grpc_address"`
}

type Bot struct {
	SlackToken string `toml:"slack_token"`
}

func Load(path string) (*Global, error) {
	c := new(Global)
	_, err := toml.DecodeFile(path, &c)
	return c, err
}
