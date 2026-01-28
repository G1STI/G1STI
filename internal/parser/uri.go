package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

type ParsedNode struct {
	Type         string
	Tag          string
	Server       string
	Port         int
	Raw          string
	ParsedFields map[string]string
}

func ParseURI(raw string) (ParsedNode, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ParsedNode{}, fmt.Errorf("empty uri")
	}
	lower := strings.ToLower(trimmed)
	switch {
	case strings.HasPrefix(lower, "vless://"):
		return parseVLESS(trimmed)
	case strings.HasPrefix(lower, "vmess://"):
		return parseVMess(trimmed)
	case strings.HasPrefix(lower, "trojan://"):
		return parseTrojan(trimmed)
	case strings.HasPrefix(lower, "ss://"):
		return parseShadowsocks(trimmed)
	case strings.HasPrefix(lower, "hysteria2://") || strings.HasPrefix(lower, "hy2://"):
		return parseHysteria2(trimmed)
	case strings.HasPrefix(lower, "tuic://"):
		return parseTuic(trimmed)
	default:
		return ParsedNode{}, fmt.Errorf("unsupported uri scheme")
	}
}

func parseVLESS(raw string) (ParsedNode, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("parse vless uri: %w", err)
	}
	if parsed.User == nil {
		return ParsedNode{}, fmt.Errorf("vless uri missing userinfo")
	}
	uuid := parsed.User.Username()
	host, port, err := splitHostPort(parsed.Host)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("vless host/port: %w", err)
	}

	fields := map[string]string{
		"uuid": uuid,
	}
	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		fields[key] = values[len(values)-1]
	}

	tag := fragmentOrHost(parsed.Fragment, host, port)
	return ParsedNode{
		Type:         "vless",
		Tag:          tag,
		Server:       host,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func parseVMess(raw string) (ParsedNode, error) {
	payload := strings.TrimSpace(strings.TrimPrefix(raw, "vmess://"))
	decoded, err := decodeBase64(payload)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("decode vmess payload: %w", err)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(decoded, &data); err != nil {
		return ParsedNode{}, fmt.Errorf("parse vmess json: %w", err)
	}
	server, _ := data["add"].(string)
	portStr := fmt.Sprintf("%v", data["port"])
	port, _ := strconv.Atoi(portStr)
	tag, _ := data["ps"].(string)
	fields := mapStringValues(data)

	if tag == "" {
		tag = fragmentOrHost("", server, port)
	}
	return ParsedNode{
		Type:         "vmess",
		Tag:          tag,
		Server:       server,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func parseTrojan(raw string) (ParsedNode, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("parse trojan uri: %w", err)
	}
	if parsed.User == nil {
		return ParsedNode{}, fmt.Errorf("trojan uri missing userinfo")
	}
	password := parsed.User.Username()
	host, port, err := splitHostPort(parsed.Host)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("trojan host/port: %w", err)
	}
	fields := map[string]string{"password": password}
	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		fields[key] = values[len(values)-1]
	}
	tag := fragmentOrHost(parsed.Fragment, host, port)
	return ParsedNode{
		Type:         "trojan",
		Tag:          tag,
		Server:       host,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func parseShadowsocks(raw string) (ParsedNode, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("parse ss uri: %w", err)
	}

	fields := map[string]string{}
	host := parsed.Hostname()
	portStr := parsed.Port()
	port := 0
	if portStr != "" {
		port, _ = strconv.Atoi(portStr)
	}

	if host == "" {
		payload := strings.TrimPrefix(raw, "ss://")
		if idx := strings.Index(payload, "#"); idx >= 0 {
			payload = payload[:idx]
		}
		decoded, err := decodeBase64(payload)
		if err == nil {
			parts := strings.SplitN(string(decoded), "@", 2)
			if len(parts) == 2 {
				userinfo := strings.SplitN(parts[0], ":", 2)
				if len(userinfo) == 2 {
					fields["method"] = userinfo[0]
					fields["password"] = userinfo[1]
				}
				hostPart := parts[1]
				host, port, _ = splitHostPort(hostPart)
			}
		}
	}

	tag := fragmentOrHost(parsed.Fragment, host, port)
	return ParsedNode{
		Type:         "shadowsocks",
		Tag:          tag,
		Server:       host,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func parseHysteria2(raw string) (ParsedNode, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("parse hysteria2 uri: %w", err)
	}
	host, port, err := splitHostPort(parsed.Host)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("hysteria2 host/port: %w", err)
	}
	fields := map[string]string{}
	if parsed.User != nil {
		fields["auth"] = parsed.User.Username()
	}
	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		fields[key] = values[len(values)-1]
	}
	tag := fragmentOrHost(parsed.Fragment, host, port)
	return ParsedNode{
		Type:         "hysteria2",
		Tag:          tag,
		Server:       host,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func parseTuic(raw string) (ParsedNode, error) {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("parse tuic uri: %w", err)
	}
	host, port, err := splitHostPort(parsed.Host)
	if err != nil {
		return ParsedNode{}, fmt.Errorf("tuic host/port: %w", err)
	}
	fields := map[string]string{}
	if parsed.User != nil {
		fields["uuid"] = parsed.User.Username()
		if pwd, hasPwd := parsed.User.Password(); hasPwd {
			fields["password"] = pwd
		}
	}
	for key, values := range parsed.Query() {
		if len(values) == 0 {
			continue
		}
		fields[key] = values[len(values)-1]
	}
	tag := fragmentOrHost(parsed.Fragment, host, port)
	return ParsedNode{
		Type:         "tuic",
		Tag:          tag,
		Server:       host,
		Port:         port,
		Raw:          raw,
		ParsedFields: fields,
	}, nil
}

func splitHostPort(hostport string) (string, int, error) {
	if hostport == "" {
		return "", 0, fmt.Errorf("empty host")
	}
	if strings.Contains(hostport, ":") {
		host, portStr, err := net.SplitHostPort(hostport)
		if err != nil {
			return "", 0, err
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return "", 0, fmt.Errorf("invalid port %q", portStr)
		}
		return host, port, nil
	}
	return hostport, 0, fmt.Errorf("missing port")
}

func fragmentOrHost(fragment string, host string, port int) string {
	if fragment != "" {
		return fragment
	}
	if port > 0 {
		return fmt.Sprintf("%s:%d", host, port)
	}
	return host
}

func decodeBase64(payload string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err == nil {
		return decoded, nil
	}
	decoded, err = base64.RawStdEncoding.DecodeString(payload)
	if err == nil {
		return decoded, nil
	}
	return nil, err
}

func DecodeBase64String(payload string) (string, error) {
	decoded, err := decodeBase64(payload)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func mapStringValues(data map[string]interface{}) map[string]string {
	fields := make(map[string]string, len(data))
	for key, value := range data {
		fields[key] = fmt.Sprintf("%v", value)
	}
	return fields
}
