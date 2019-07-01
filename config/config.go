package config

import (
	"fmt"

	"github.com/yumimobi/trace/util"

	"github.com/yumimobi/trace/log"

	"github.com/spf13/viper"
)

var Conf = &Config{}

type RemoteAddress struct {
	TestAddr       []string `mapstructure:"test_addr"`
	ProductionAddr []string `mapstructure:"production_addr"`
}

type HTTPConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type RPCConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type GRPCConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type TargetDirectory struct {
	Dir string `mapstructure:"dir"`
}

type WebSocketConfig struct {
	Address string `mapstructure:"address"`
	Port    string `mapstructure:"port"`
}

type Server struct {
	HTTP          HTTPConfig      `mapstructure:"http"`
	RPC           RPCConfig       `mapstructure:"rpc"`
	GRPC          GRPCConfig      `mapstructure:"grpc"`
	WebSocket     WebSocketConfig `mapstructure:"websocket"`
	Log           log.Config      `mapstructure:"log"`
	RemoteAddress RemoteAddress   `mapstructure:"remote_address"`
}

type Client struct {
	HTTP   HTTPConfig      `mapstructure:"http"`
	RPC    RPCConfig       `mapstructure:"rpc"`
	GRPC   GRPCConfig      `mapstructure:"grpc"`
	Log    log.Config      `mapstructure:"log"`
	Target TargetDirectory `mapstructure:"target_directory"`
}

type Config struct {
	Server Server `mapstructure:"server"`
	Client Client `mapstructure:"client"`
}

// New a config
func Init() error {
	err := load(Conf)
	if err != nil {
		return err
	}

	return nil
}

func load(c *Config) error {
	filename := "./conf.yaml"
	has, _ := util.PathExists(filename)
	if !has {
		filename = "../conf.yaml"
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("read config is failed, err:%v", err)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("unmarshal config is failed, err:%v", err)
	}

	fmt.Println("------", *c)
	return nil
}
