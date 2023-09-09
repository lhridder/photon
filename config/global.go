package config

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lhridder/photon/protocol"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type DefaultStatus struct {
	Version         string
	Protocolversion int
	Description     string
	IconPath        string
}

type Globalconfig struct {
	Debug                 bool          `yaml:"debug"`
	Status                DefaultStatus `yaml:"status"`
	DefaultStatusResponse protocol.StatusResponse
}

func loadFavicon(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	imgFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer imgFile.Close()

	fileInfo, err := imgFile.Stat()
	if err != nil {
		return "", err
	}

	buffer := make([]byte, fileInfo.Size())
	fileReader := bufio.NewReader(imgFile)
	_, err = fileReader.Read(buffer)
	if err != nil {
		return "", nil
	}

	return base64.StdEncoding.EncodeToString(buffer), nil
}

func loadDefaultStatus(cfg *Globalconfig) error {
	desc := fmt.Sprintf(fmt.Sprintf("{\"text\":\"%s\"}", cfg.Status.Description))

	favicon := ""
	if cfg.Status.IconPath != "" {
		img64, err := loadFavicon(cfg.Status.IconPath)
		if err != nil {
			return fmt.Errorf("failed to load favicon: %s", err)
		}
		favicon = fmt.Sprintf("data:image/png;base64,%s", img64)
	}

	cfg.DefaultStatusResponse = protocol.StatusResponse{
		Version: protocol.VersionJSON{
			Name:     cfg.Status.Version,
			Protocol: cfg.Status.Protocolversion,
		},
		Players: protocol.PlayersJSON{
			Max:    0,
			Online: 0,
		},
		Description: json.RawMessage(desc),
		Favicon:     favicon,
	}

	return nil
}

func LoadGlobal() (*Globalconfig, error) {
	var cfg *Globalconfig
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config.yml: %s", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yml: %s", err)
	}

	err = loadDefaultStatus(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load default status: %s", err)
	}

	return cfg, nil
}
