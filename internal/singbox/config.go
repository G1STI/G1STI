package singbox

import (
	"encoding/json"
	"fmt"
	"os"

	"g1sti/internal/model"
)

type Config struct {
	Log       LogConfig       `json:"log"`
	DNS       DNSConfig       `json:"dns"`
	Inbounds  []InboundConfig `json:"inbounds"`
	Outbounds []Outbound      `json:"outbounds"`
	Route     RouteConfig     `json:"route"`
}

type LogConfig struct {
	Level string `json:"level"`
	File  string `json:"output"`
}

type DNSConfig struct {
	Servers []DNSServer `json:"servers"`
}

type DNSServer struct {
	Tag     string `json:"tag"`
	Address string `json:"address"`
}

type InboundConfig struct {
	Type                   string                 `json:"type"`
	Tag                    string                 `json:"tag"`
	Interface              string                 `json:"interface,omitempty"`
	MTU                    int                    `json:"mtu,omitempty"`
	AutoRoute              bool                   `json:"auto_route,omitempty"`
	StrictRoute            bool                   `json:"strict_route,omitempty"`
	Stack                  string                 `json:"stack,omitempty"`
	Sniff                  bool                   `json:"sniff,omitempty"`
	EndpointIndependentNat bool                   `json:"endpoint_independent_nat,omitempty"`
	Extra                  map[string]interface{} `json:"-"`
}

type Outbound struct {
	Type   string                 `json:"type"`
	Tag    string                 `json:"tag"`
	Server string                 `json:"server,omitempty"`
	Port   int                    `json:"server_port,omitempty"`
	TLS    map[string]interface{} `json:"tls,omitempty"`
	Extra  map[string]interface{} `json:"-"`
}

type RouteConfig struct {
	Rules               []RouteRule `json:"rules"`
	Final               string      `json:"final"`
	AutoDetectInterface bool        `json:"auto_detect_interface"`
}

type RouteRule struct {
	Outbound string   `json:"outbound"`
	IPCIDR   []string `json:"ip_cidr,omitempty"`
	Domain   []string `json:"domain,omitempty"`
}

func BuildConfig(profile model.Profile, node model.Node, logPath string) (Config, error) {
	if node.Type == "" {
		return Config{}, fmt.Errorf("node type is empty")
	}

	config := Config{
		Log: LogConfig{
			Level: "info",
			File:  logPath,
		},
		DNS: DNSConfig{
			Servers: []DNSServer{{
				Tag:     "system",
				Address: "system",
			}},
		},
		Inbounds:  []InboundConfig{},
		Outbounds: []Outbound{},
		Route: RouteConfig{
			Rules:               []RouteRule{},
			Final:               "proxy",
			AutoDetectInterface: true,
		},
	}

	if profile.TunEnabled {
		config.Inbounds = append(config.Inbounds, InboundConfig{
			Type:                   "tun",
			Tag:                    "tun-in",
			Interface:              "auto",
			MTU:                    1500,
			AutoRoute:              true,
			StrictRoute:            true,
			Stack:                  "system",
			Sniff:                  true,
			EndpointIndependentNat: true,
		})
	}

	proxyOutbound := Outbound{
		Type:   node.Type,
		Tag:    "proxy",
		Server: node.Server,
		Port:   node.Port,
	}
	config.Outbounds = append(config.Outbounds,
		proxyOutbound,
		Outbound{Type: "direct", Tag: "direct"},
		Outbound{Type: "block", Tag: "block"},
	)

	return config, nil
}

func WriteConfig(path string, config Config) error {
	payload, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal sing-box config: %w", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		return fmt.Errorf("write sing-box config: %w", err)
	}
	return nil
}
