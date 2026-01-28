package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"g1sti/internal/model"
)

const dataFileName = "data.json"

var ErrNotInitialized = errors.New("storage not initialized")

type Store struct {
	rootDir string
	path    string
}

func NewStore(rootDir string) *Store {
	return &Store{
		rootDir: rootDir,
		path:    filepath.Join(rootDir, dataFileName),
	}
}

func (s *Store) Ensure() error {
	if s.rootDir == "" {
		return ErrNotInitialized
	}
	return os.MkdirAll(s.rootDir, 0o755)
}

func (s *Store) Load() (model.AppData, error) {
	if s.path == "" {
		return model.AppData{}, ErrNotInitialized
	}

	payload, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return model.DefaultAppData(), nil
	}
	if err != nil {
		return model.AppData{}, fmt.Errorf("read storage file: %w", err)
	}

	var data model.AppData
	if err := json.Unmarshal(payload, &data); err != nil {
		legacy, legacyErr := decodeLegacy(payload)
		if legacyErr != nil {
			return model.AppData{}, fmt.Errorf("decode storage file: %w", err)
		}
		data = legacy
	}

	migrated, err := migrate(data)
	if err != nil {
		return model.AppData{}, err
	}

	return migrated, nil
}

func (s *Store) Save(data model.AppData) error {
	if s.path == "" {
		return ErrNotInitialized
	}
	if err := s.Ensure(); err != nil {
		return err
	}

	data.Version = model.SchemaVersion
	payload, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal storage data: %w", err)
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, payload, 0o644); err != nil {
		return fmt.Errorf("write storage temp file: %w", err)
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		return fmt.Errorf("rename storage temp file: %w", err)
	}

	return nil
}

func migrate(data model.AppData) (model.AppData, error) {
	if data.Version == 0 {
		data.Version = model.SchemaVersion
	}
	if data.Version > model.SchemaVersion {
		return model.AppData{}, fmt.Errorf("storage schema version %d is newer than supported %d", data.Version, model.SchemaVersion)
	}
	return data, nil
}

type legacyAppData struct {
	Version       int                  `json:"version"`
	Subscriptions []legacySubscription `json:"subscriptions"`
	Nodes         []model.Node         `json:"nodes"`
	Profiles      []legacyProfile      `json:"profiles"`
	Settings      model.Settings       `json:"settings"`
	ActiveProfile string               `json:"active_profile"`
}

type legacySubscription struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	AutoUpdate   string    `json:"auto_update"`
	LastUpdated  time.Time `json:"last_updated"`
	LastError    string    `json:"last_error"`
	NodeCount    int       `json:"node_count"`
	LastChecksum string    `json:"last_checksum"`
}

type legacyProfile struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	SubscriptionID string    `json:"subscription_id"`
	NodeID         string    `json:"node_id"`
	TunEnabled     bool      `json:"tun_enabled"`
	DNSMode        string    `json:"dns_mode"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func decodeLegacy(payload []byte) (model.AppData, error) {
	var legacy legacyAppData
	if err := json.Unmarshal(payload, &legacy); err != nil {
		return model.AppData{}, err
	}
	data := model.AppData{
		Version:       model.SchemaVersion,
		Nodes:         legacy.Nodes,
		Settings:      legacy.Settings,
		ActiveProfile: legacy.ActiveProfile,
	}
	data.Subscriptions = make([]model.Subscription, 0, len(legacy.Subscriptions))
	for _, sub := range legacy.Subscriptions {
		data.Subscriptions = append(data.Subscriptions, model.Subscription{
			ID:           sub.ID,
			Name:         sub.Name,
			URL:          sub.URL,
			AutoUpdate:   sub.AutoUpdate,
			LastUpdated:  sub.LastUpdated.UnixMilli(),
			LastError:    sub.LastError,
			NodeCount:    sub.NodeCount,
			LastChecksum: sub.LastChecksum,
		})
	}
	data.Profiles = make([]model.Profile, 0, len(legacy.Profiles))
	for _, profile := range legacy.Profiles {
		data.Profiles = append(data.Profiles, model.Profile{
			ID:             profile.ID,
			Name:           profile.Name,
			SubscriptionID: profile.SubscriptionID,
			NodeID:         profile.NodeID,
			TunEnabled:     profile.TunEnabled,
			DNSMode:        profile.DNSMode,
			CreatedAt:      profile.CreatedAt.UnixMilli(),
			UpdatedAt:      profile.UpdatedAt.UnixMilli(),
		})
	}
	return data, nil
}
