package config

import "github.com/BurntSushi/toml"

type Global struct {
	DSPD  Dspd  `toml:"dspd"`
	Textd Textd `toml:"textd"`
}

type Dspd struct {
	GRPCServerAddress string `toml:"grpc_address"`
}

type Textd struct {
	GRPCServerAddress string `toml:"grpc_address"`
}

func Load(path string) (*Global, error) {
	c := new(Global)
	_, err := toml.DecodeFile(path, &c)
	return c, err
}
