package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Conn struct {
		Host     string `yaml:"host"`
		Endpoint string `yaml:"endpoint"`
		Origin   string `yaml:"origin"`
	} `yaml:"connection"`
	Auth struct {
		ApiKey    string `yaml:"key"`
		ApiSecret string `yaml:"secret"`
		Password  string `yaml:"password"`
	} `yaml:"authentication"`
}

type message struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

func main() {

	path, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		log.Fatal(errors.New("Unable to find config file path. Make sure CONFIG_FILE is set"))
	}

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	u := url.URL{
		Scheme: "wss",
		Host:   cfg.Conn.Host,
		Path:   cfg.Conn.Endpoint,
	}

	connCfg, err := websocket.NewConfig(u.String(), cfg.Conn.Origin)
	if err != nil {
		log.Fatal(err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	payload := timestamp + "GET" + cfg.Conn.Endpoint
	hashMac := hmac.New(sha256.New, []byte(cfg.Auth.ApiSecret))
	hashMac.Write([]byte(payload))
	expectedMac := hashMac.Sum(nil)
	signature := hex.EncodeToString(expectedMac)

	connCfg.Dialer = &net.Dialer{
		Timeout: 5 * time.Second,
	}

	connCfg.Header.Add("X-QUME-API-KEY", cfg.Auth.ApiKey)
	connCfg.Header.Add("X-QUME-SIGNATURE", signature)
	connCfg.Header.Add("X-QUME-TIMESTAMP", timestamp)
	connCfg.Header.Add("X-QUME-PASSPHRASE", cfg.Auth.Password)

	conn, err := websocket.DialConfig(connCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		var msg message
		if err := websocket.JSON.Receive(conn, &msg); err != nil {
			log.Fatal(err)
			break
		}
		log.Println(msg)
	}
}
