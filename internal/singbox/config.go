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

func (o Outbound) MarshalJSON() ([]byte, error) {
	type base Outbound
	payload := map[string]interface{}{}
	raw, err := json.Marshal(base(o))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	for key, value := range o.Extra {
		payload[key] = value
	}
	return json.Marshal(payload)
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

	proxyOutbound, err := buildOutbound(node)
	if err != nil {
		return Config{}, err
	}
	proxyOutbound.Tag = "proxy"
	config.Outbounds = append(config.Outbounds,
		proxyOutbound,
		Outbound{Type: "direct", Tag: "direct"},
		Outbound{Type: "block", Tag: "block"},
	)

	return config, nil
}

func buildOutbound(node model.Node) (Outbound, error) {
	outbound := Outbound{
		Type:   node.Type,
		Server: node.Server,
		Port:   node.Port,
		Extra:  map[string]interface{}{},
	}
	switch node.Type {
	case "vless":
		uuid := node.ParsedFields["uuid"]
		if uuid == "" {
			return Outbound{}, fmt.Errorf("vless uuid is empty")
		}
		outbound.Extra["uuid"] = uuid
		if flow := node.ParsedFields["flow"]; flow != "" {
			outbound.Extra["flow"] = flow
		}
		applyTLS(&outbound, node)
		applyTransport(&outbound, node)
	case "trojan":
		password := node.ParsedFields["password"]
		if password == "" {
			return Outbound{}, fmt.Errorf("trojan password is empty")
		}
		outbound.Extra["password"] = password
		applyTLS(&outbound, node)
	case "shadowsocks":
		method := node.ParsedFields["method"]
		password := node.ParsedFields["password"]
		if method == "" || password == "" {
			return Outbound{}, fmt.Errorf("shadowsocks method/password missing")
		}
		outbound.Extra["method"] = method
		outbound.Extra["password"] = password
	case "hysteria2":
		if auth := node.ParsedFields["auth"]; auth != "" {
			outbound.Extra["auth"] = auth
		}
		applyTLS(&outbound, node)
	case "tuic":
		uuid := node.ParsedFields["uuid"]
		password := node.ParsedFields["password"]
		if uuid == "" || password == "" {
			return Outbound{}, fmt.Errorf("tuic uuid/password missing")
		}
		outbound.Extra["uuid"] = uuid
		outbound.Extra["password"] = password
		applyTLS(&outbound, node)
	}
	return outbound, nil
}

func applyTLS(outbound *Outbound, node model.Node) {
	security := node.ParsedFields["security"]
	sni := node.ParsedFields["sni"]
	fp := node.ParsedFields["fp"]
	pbk := node.ParsedFields["pbk"]
	sid := node.ParsedFields["sid"]

	if security == "" && sni == "" && fp == "" && pbk == "" {
		return
	}

	tls := map[string]interface{}{}
	if security != "" && security != "none" {
		tls["enabled"] = true
	}
	if sni != "" {
		tls["server_name"] = sni
	}
	if fp != "" {
		tls["utls"] = map[string]interface{}{
			"enabled":     true,
			"fingerprint": fp,
		}
	}
	if security == "reality" || pbk != "" {
		reality := map[string]interface{}{
			"enabled": true,
		}
		if pbk != "" {
			reality["public_key"] = pbk
		}
		if sid != "" {
			reality["short_id"] = sid
		}
		tls["reality"] = reality
	}
	outbound.TLS = tls
}

func applyTransport(outbound *Outbound, node model.Node) {
	transport := node.ParsedFields["type"]
	if transport == "" {
		return
	}
	outbound.Extra["transport"] = map[string]interface{}{
		"type": transport,
	}
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
