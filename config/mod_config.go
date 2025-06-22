package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	sqlcGrpcYamlConfig = "sqlc-grpc.yaml"
	sqlcGrpcYmlConfig  = "sqlc-grpc.yml"
	sqlcGrpcJsonConfig = "sqlc-grpc.json"
)

type ModField struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
}

type ModConfig struct {
	IgnoreServices []string `json:"ignore_services" yaml:"ignore_services"`
	RemoveServices []string `json:"remove_services" yaml:"remove_services"`
	AddServices    []struct {
		Name      string      `json:"name" yaml:"name"`
		Method    string      `json:"method" yaml:"method"`
		Path      string      `json:"path" yaml:"path"`
		ReqFields []*ModField `json:"req_fields" yaml:"req_fields"`
		ResFields []*ModField `json:"res_fields" yaml:"res_fields"`
	} `json:"add_services" yaml:"add_services"`
	RemoveFields []string `json:"remove_fields" yaml:"remove_fields"`
	AddFields    []struct {
		Msg    string      `json:"msg" yaml:"msg"`
		Fields []*ModField `json:"fields" yaml:"fields"`
	} `json:"add_fields" yaml:"add_fields"`
	RolesFilter map[string]ModConfig `json:"roles_filter" yaml:"roles_filter"`
	Packages    map[string]ModConfig `json:"packages" yaml:"packages"`
}

func modConfigFile() (string, error) {
	if f, err := os.Stat(sqlcGrpcYmlConfig); err == nil && !f.IsDir() {
		return sqlcGrpcYmlConfig, nil
	}

	if f, err := os.Stat(sqlcGrpcYamlConfig); err == nil && !f.IsDir() {
		return sqlcGrpcYamlConfig, nil
	}

	if f, err := os.Stat(sqlcGrpcJsonConfig); err == nil && !f.IsDir() {
		return sqlcGrpcJsonConfig, nil
	}
	return "", errors.New("no sqlc-grpc config files (sqlc-grpc.json, sqlc-grpc.yaml or sqlc-grpc.yml)")
}

func LoadModConfig() (ModConfig, error) {
	var cfg ModConfig
	name, err := modConfigFile()
	if err != nil {
		return cfg, nil
	}

	f, err := os.Open(name)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	switch name {
	case sqlcGrpcJsonConfig:
		err = json.NewDecoder(f).Decode(&cfg)
	default:
		err = yaml.NewDecoder(f).Decode(&cfg)
	}
	if errors.Is(err, io.EOF) {
		return cfg, nil
	}
	return cfg, err
}
