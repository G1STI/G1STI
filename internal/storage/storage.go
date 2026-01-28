package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	file, err := os.Open(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return model.DefaultAppData(), nil
	}
	if err != nil {
		return model.AppData{}, fmt.Errorf("open storage file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var data model.AppData
	if err := decoder.Decode(&data); err != nil {
		return model.AppData{}, fmt.Errorf("decode storage file: %w", err)
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
