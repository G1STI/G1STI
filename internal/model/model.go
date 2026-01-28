package model

import "time"

const SchemaVersion = 1

type AppData struct {
	Version       int            `json:"version"`
	Subscriptions []Subscription `json:"subscriptions"`
	Nodes         []Node         `json:"nodes"`
	Profiles      []Profile      `json:"profiles"`
	Settings      Settings       `json:"settings"`
	ActiveProfile string         `json:"active_profile"`
}

type Subscription struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	AutoUpdate   string    `json:"auto_update"`
	LastUpdated  time.Time `json:"last_updated"`
	LastError    string    `json:"last_error"`
	NodeCount    int       `json:"node_count"`
	LastChecksum string    `json:"last_checksum"`
}

type Node struct {
	ID             string            `json:"id"`
	SubscriptionID string            `json:"subscription_id"`
	Tag            string            `json:"tag"`
	Type           string            `json:"type"`
	Server         string            `json:"server"`
	Port           int               `json:"port"`
	Raw            string            `json:"raw"`
	ParsedFields   map[string]string `json:"parsed_fields"`
}

type Profile struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	SubscriptionID string    `json:"subscription_id"`
	NodeID         string    `json:"node_id"`
	TunEnabled     bool      `json:"tun_enabled"`
	DNSMode        string    `json:"dns_mode"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Settings struct {
	LogLevel      string `json:"log_level"`
	SingBoxPath   string `json:"sing_box_path"`
	WintunDLLPath string `json:"wintun_dll_path"`
}

func DefaultAppData() AppData {
	return AppData{
		Version: SchemaVersion,
		Settings: Settings{
			LogLevel: "info",
		},
	}
}
