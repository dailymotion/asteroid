package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const FILENAME = ".asteroid.yaml"

type Config struct {
	WG 			Wireguard 	`yaml:"wireguard"`
	ClientConfig ClientConfig `yaml:"client_config_file"`
}

type Wireguard struct {
	SSHKeyName  	string `yaml:"ssh_key_name"`
	WireguardIP 	string `yaml:"ip"`
	SSHPort     	string `yaml:"ssh_port"`
	Username    	string `yaml:"username"`
	HostKey     	bool   `yaml:"verification_host_key"`
	WGPort			string `yaml:"wg_port"`
}

type ClientConfig struct {
	Name     	string `yaml:"name"`
	DNS     	string `yaml:"dns"`
	AllowedIPs  string `yaml:"allowed_ips"`
}

func isStructNil(config Config) ([]string, bool) {
	e := reflect.ValueOf(&config).Elem()
	num := e.NumField()
	var listMissing []string
	var isNil bool

	for i := 0; i < num; i++ {
		fieldTagTmp := string(e.Type().Field(i).Tag)
		fieldTag := strings.Split(fieldTagTmp, "\"")
		fieldValue := e.Field(i).Interface()

		if fieldValue == "" {
			listMissing = append(listMissing, fieldTag[1])
			isNil = true
		}
	}
	if isNil {
		return listMissing, true
	}
	return listMissing, false
}

func ReadConfigFile() (Config, error) {
	configWG := Config{}

	usr, err := user.Current()
	if err != nil {
		return configWG, errors.Wrap(err, "failed to retrieve the current user")
	}
	path := usr.HomeDir + "/" + FILENAME

	// check if file exists
	_, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		return configWG, errors.Wrap(err, "couldn't create file")
	}
	if err != nil {
		return configWG, errors.Wrap(err, "failed to get file info")
	}

	// reading file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return configWG, errors.Wrap(err, "failed to read the file")
	}

	err = yaml.Unmarshal(data, &configWG)
	if err != nil {
		return configWG, errors.Wrap(err, "failed to unmarshall the data from the file")
	}

	listMissing, isNil := isStructNil(configWG)
	if isNil {
		switch len(listMissing) {
		case 0:
			return configWG, errors.New("There is an issue with your config file")
		case 1:
			return configWG, errors.Errorf("The field %s in your config file is missing", listMissing[0])
		default:
			return configWG, errors.Errorf("The fields %v in your config file are missing\n", listMissing)
		}
	}

	return configWG, nil
}
