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
	SSHKeyName	string `yaml:"ssh_key_name"`
	WireguardIP string `yaml:"wireguard_ip"`
	SSHPort		string `yaml:"ssh_port"`
	Username	string `yaml:"username"`
	HostKey		bool   `yaml:"verification_host_key"`
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
		return configWG, err
	}
	path := usr.HomeDir + "/" + FILENAME

	// check if file exists
	_, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		return configWG, errors.Wrap(err, "couldn't create file")
	}

	// reading file
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return configWG, errors.Wrap(err, "couldn't read file")
	}

	err = yaml.Unmarshal(data, &configWG)
	if err != nil {
		return configWG, errors.Wrap(err, "error unmarshal data")
	}

	listMissing, isNil := isStructNil(configWG)
	if isNil {
		if len(listMissing) == 1 {
			return configWG, errors.Wrapf(err, "\nThe field %v in your config file is missing\n", listMissing)
			//fmt.Printf("\nThe field %v in your config file is missing\n", listMissing)
			//os.Exit(1)
		} else  if len(listMissing) >= 2 {
			return configWG, errors.Wrapf(err, "\nThe fields %v in your config file are missing\n", listMissing)
			//fmt.Printf("\nThe fields %v in your config file are missing\n", listMissing)
			//os.Exit(1)
		} else {
			return configWG, errors.Wrapf(err, "\nThere is an issue with your config file\n", listMissing)
			//fmt.Printf("\nThere is an issue with your config file\n")
			//os.Exit(1)
		}
	}

	return configWG, nil
}