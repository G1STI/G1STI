package app

import (
	"encoding/json"
	"fmt"
	"time"

	"g1sti/internal/model"
)

type ProfileExport struct {
	Profile model.Profile `json:"profile"`
	Node    model.Node    `json:"node"`
}

func (a *App) ExportProfile(profileID string) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, profile := range a.data.Profiles {
		if profile.ID != profileID {
			continue
		}
		for _, node := range a.data.Nodes {
			if node.ID == profile.NodeID {
				payload, err := json.MarshalIndent(ProfileExport{Profile: profile, Node: node}, "", "  ")
				if err != nil {
					return "", fmt.Errorf("marshal profile export: %w", err)
				}
				return string(payload), nil
			}
		}
		return "", ErrNodeNotFound
	}
	return "", ErrProfileNotFound
}

func (a *App) ImportProfile(payload string) (model.Profile, error) {
	var export ProfileExport
	if err := json.Unmarshal([]byte(payload), &export); err != nil {
		return model.Profile{}, fmt.Errorf("parse profile export: %w", err)
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	profile := export.Profile
	profile.ID = newID()
	profile.CreatedAt = time.Now().UnixMilli()
	profile.UpdatedAt = profile.CreatedAt
	node := export.Node
	node.ID = newID()
	profile.NodeID = node.ID
	a.data.Nodes = append(a.data.Nodes, node)
	a.data.Profiles = append(a.data.Profiles, profile)
	a.lastSeen = time.Now().UnixMilli()
	if err := a.store.Save(a.data); err != nil {
		return model.Profile{}, err
	}
	return profile, nil
}
