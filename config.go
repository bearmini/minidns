package main

import (
	"io/ioutil"
	"net"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v1"
)

type RawConfig struct {
	Ports   Ports                  `yaml:"ports"`
	Records map[string][]*RawEntry `yaml:"records"`
}

type Ports struct {
	UDP4 uint16 `yaml:"udp4"`
	TCP4 uint16 `yaml:"tcp4"`
}

type RawEntry struct {
	Type  string `yaml:"type"`
	Class string `yaml:"class"`
	TTL   uint32 `yaml:"ttl"`
	A     string `yaml:"A"`
}

type Config struct {
	Ports   Ports `yaml:"ports"`
	Records map[string][]*Entry
}

type Entry struct {
	Type  uint16
	Class uint16
	TTL   uint32
	A     net.IP
}

func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to open file: %s", path)
	}

	yamlData := strings.TrimSpace(string(data))

	var rc RawConfig
	err = yaml.Unmarshal([]byte(yamlData), &rc)
	if err != nil {
		return nil, err
	}

	return convertConfig(&rc)
}

func convertConfig(rc *RawConfig) (*Config, error) {
	var cfg Config
	cfg.Ports = rc.Ports
	cfg.Records = make(map[string][]*Entry)

	for name, rr := range rc.Records {
		r, err := convertRecord(rr)
		if err != nil {
			return nil, err
		}
		cfg.Records[name] = r
	}

	return &cfg, nil
}

func convertRecord(src []*RawEntry) ([]*Entry, error) {
	dst := make([]*Entry, 0, len(src))

	for _, e := range src {
		etype, err := convertRRType(e.Type)
		if err != nil {
			return nil, err
		}

		eclass, err := convertClass(e.Class)
		if err != nil {
			return nil, err
		}

		if e.TTL == 0 {
			e.TTL = 60
		}

		dst = append(dst, &Entry{
			Type:  etype,
			Class: eclass,
			TTL:   e.TTL,
			A:     net.ParseIP(e.A),
		})
	}

	return dst, nil
}

func convertRRType(src string) (uint16, error) {
	if src == "" {
		src = "A"
	}
	etype, found := g_rrTypeMap[src]
	if !found {
		return 0, errors.Errorf("unknown or unsupported type: %s\n", src)
	}
	return etype, nil
}

func convertClass(src string) (uint16, error) {
	if src == "" {
		src = "IN"
	}
	eclass, found := g_classMap[src]
	if !found {
		return 0, errors.Errorf("unknown or unsupported class: %s\n", src)
	}
	return eclass, nil
}
