package main

import (
	"log"
	"github.com/pelletier/go-toml"
	"strings"
	"os"
	"fmt"
	"bufio"
	"io/ioutil"
	"os/user"
	"errors"
)

func getConfDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.New("can't get home path")
	}
	confDir := strings.Join([]string{usr.HomeDir, ".config", "megaCli"}, "/")
	_ = os.Mkdir(confDir, os.ModePerm)
	return confDir, nil
}

func setConf() Config {
	fmt.Println("input your account info!")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, " \n")
	fmt.Print("Password: ")
	password, _ := reader.ReadString('\n')
	password = strings.Trim(password, " \n")
	return Config{username, password}
}

func parseConf() (Config, error) {
	confDir, err := getConfDir()
	if err != nil {
		return Config{}, err
	}
	confPath := strings.Join([]string{confDir, "mega.toml"}, "/")
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return Config{}, err
	}

	f, err := os.Open(confPath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()
	bytesConfig, err := ioutil.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = toml.Unmarshal(bytesConfig, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil

}

func initConf(config Config) error {
	confDir, err := getConfDir()
	if err != nil {
		return err
	}
	confPath := strings.Join([]string{confDir, "mega.toml"}, "/")
	f, err := os.OpenFile(confPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	bytesConfig, err := toml.Marshal(config)
	if err != nil {
		panic(err)
	}

	_, err = f.Write(bytesConfig)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	return nil
}