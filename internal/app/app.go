package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"g1sti/internal/model"
	"g1sti/internal/parser"
	"g1sti/internal/paths"
	"g1sti/internal/singbox"
	"g1sti/internal/storage"
	"g1sti/internal/subscription"
)

var ErrProfileNotFound = errors.New("profile not found")
var ErrNodeNotFound = errors.New("node not found")
var ErrInvalidURI = errors.New("invalid uri")

const manualSubscriptionName = "Manual"

type Status struct {
	Connected     bool   `json:"connected"`
	ActiveProfile string `json:"active_profile"`
	UpdatedAtMs   int64  `json:"updated_at_ms"`
	LastError     string `json:"last_error"`
	PublicIP      string `json:"public_ip"`
}

type App struct {
	mu       sync.Mutex
	ctx      context.Context
	store    *storage.Store
	service  *subscription.Service
	manager  *singbox.Manager
	data     model.AppData
	lastErr  string
	lastSeen int64
	publicIP string
}

func New() *App {
	return &App{
		service: subscription.NewService(nil),
		data:    model.DefaultAppData(),
	}
}

func (a *App) Startup(ctx context.Context) {
	if err := a.startup(ctx); err != nil {
		a.mu.Lock()
		a.lastErr = err.Error()
		a.mu.Unlock()
	}
}

func (a *App) startup(ctx context.Context) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ctx = ctx
	rootDir, err := paths.RootDir()
	if err != nil {
		return err
	}

	a.store = storage.NewStore(rootDir)
	if err := a.store.Ensure(); err != nil {
		return err
	}

	data, err := a.store.Load()
	if err != nil {
		return err
	}
	a.data = data

	configPath, err := paths.SingBoxConfigPath()
	if err != nil {
		return err
	}
	logPath, err := paths.SingBoxLogPath()
	if err != nil {
		return err
	}
	a.manager = singbox.NewManager(a.data.Settings.SingBoxPath, configPath, logPath)
	return nil
}

func (a *App) Context() context.Context {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.ctx
}

func (a *App) GetStatus() Status {
	a.mu.Lock()
	defer a.mu.Unlock()
	return Status{
		Connected:     a.manager != nil && a.manager.IsRunning(),
		ActiveProfile: a.data.ActiveProfile,
		UpdatedAtMs:   a.lastSeen,
		LastError:     a.lastErr,
		PublicIP:      a.publicIP,
	}
}

func (a *App) ListSubscriptions() []model.Subscription {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]model.Subscription, len(a.data.Subscriptions))
	copy(result, a.data.Subscriptions)
	return result
}

func (a *App) ListProfiles() []model.Profile {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]model.Profile, len(a.data.Profiles))
	copy(result, a.data.Profiles)
	return result
}

func (a *App) ListNodes(subscriptionID string) []model.Node {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make([]model.Node, 0)
	for _, node := range a.data.Nodes {
		if subscriptionID == "" || node.SubscriptionID == subscriptionID {
			result = append(result, node)
		}
	}
	return result
}

func (a *App) AddSubscription(name string, url string, autoUpdate string) (model.Subscription, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	subscription := model.Subscription{
		ID:         newID(),
		Name:       name,
		URL:        url,
		AutoUpdate: autoUpdate,
	}
	a.data.Subscriptions = append(a.data.Subscriptions, subscription)
	a.lastSeen = time.Now().UnixMilli()
	return subscription, a.store.Save(a.data)
}

func (a *App) AddVlessURI(tag string, uri string) (model.Profile, error) {
	parsed, err := parser.ParseURI(uri)
	if err != nil {
		return model.Profile{}, fmt.Errorf("parse uri: %w", err)
	}
	if parsed.Type != "vless" {
		return model.Profile{}, fmt.Errorf("%w: only vless is supported", ErrInvalidURI)
	}
	if tag != "" {
		parsed.Tag = tag
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	subscription := a.ensureManualSubscription()
	node := model.Node{
		ID:             newID(),
		SubscriptionID: subscription.ID,
		Tag:            parsed.Tag,
		Type:           parsed.Type,
		Server:         parsed.Server,
		Port:           parsed.Port,
		Raw:            parsed.Raw,
		ParsedFields:   parsed.ParsedFields,
	}
	a.data.Nodes = append(a.data.Nodes, node)

	now := time.Now().UnixMilli()
	for i, item := range a.data.Subscriptions {
		if item.ID == subscription.ID {
			a.data.Subscriptions[i].NodeCount++
			a.data.Subscriptions[i].LastUpdated = now
			break
		}
	}

	profile := model.Profile{
		ID:             newID(),
		Name:           parsed.Tag,
		SubscriptionID: subscription.ID,
		NodeID:         node.ID,
		TunEnabled:     true,
		DNSMode:        "system",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	a.data.Profiles = append(a.data.Profiles, profile)
	a.data.ActiveProfile = profile.ID
	a.lastSeen = now

	if err := a.store.Save(a.data); err != nil {
		return model.Profile{}, err
	}
	return profile, nil
}

func (a *App) RemoveSubscription(subscriptionID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	filtered := a.data.Subscriptions[:0]
	for _, subscription := range a.data.Subscriptions {
		if subscription.ID != subscriptionID {
			filtered = append(filtered, subscription)
		}
	}
	a.data.Subscriptions = filtered

	nodes := a.data.Nodes[:0]
	for _, node := range a.data.Nodes {
		if node.SubscriptionID != subscriptionID {
			nodes = append(nodes, node)
		}
	}
	a.data.Nodes = nodes
	a.lastSeen = time.Now().UnixMilli()
	return a.store.Save(a.data)
}

func (a *App) CreateProfile(profile model.Profile) (model.Profile, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	profile.ID = newID()
	profile.CreatedAt = time.Now().UnixMilli()
	profile.UpdatedAt = profile.CreatedAt
	a.data.Profiles = append(a.data.Profiles, profile)
	a.lastSeen = time.Now().UnixMilli()
	return profile, a.store.Save(a.data)
}

func (a *App) UpdateProfile(profile model.Profile) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, item := range a.data.Profiles {
		if item.ID == profile.ID {
			profile.CreatedAt = item.CreatedAt
			profile.UpdatedAt = time.Now().UnixMilli()
			a.data.Profiles[i] = profile
			a.lastSeen = time.Now().UnixMilli()
			return a.store.Save(a.data)
		}
	}
	return ErrProfileNotFound
}

func (a *App) DeleteProfile(profileID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	filtered := a.data.Profiles[:0]
	for _, profile := range a.data.Profiles {
		if profile.ID != profileID {
			filtered = append(filtered, profile)
		}
	}
	a.data.Profiles = filtered
	if a.data.ActiveProfile == profileID {
		a.data.ActiveProfile = ""
	}
	a.lastSeen = time.Now().UnixMilli()
	return a.store.Save(a.data)
}

func (a *App) SetActiveProfile(profileID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, profile := range a.data.Profiles {
		if profile.ID == profileID {
			a.data.ActiveProfile = profileID
			a.lastSeen = time.Now().UnixMilli()
			return a.store.Save(a.data)
		}
	}
	return ErrProfileNotFound
}

func (a *App) UpdateSubscription(subscriptionID string) error {
	a.mu.Lock()
	var target *model.Subscription
	for i := range a.data.Subscriptions {
		if a.data.Subscriptions[i].ID == subscriptionID {
			target = &a.data.Subscriptions[i]
			break
		}
	}
	if target == nil {
		a.mu.Unlock()
		return fmt.Errorf("subscription not found")
	}
	url := target.URL
	a.mu.Unlock()

	nodes, checksum, err := a.service.Fetch(a.ctx, url, subscriptionID)
	a.mu.Lock()
	defer a.mu.Unlock()
	if err != nil {
		a.lastErr = err.Error()
		target.LastError = err.Error()
		_ = a.store.Save(a.data)
		return err
	}

	filtered := a.data.Nodes[:0]
	for _, node := range a.data.Nodes {
		if node.SubscriptionID != subscriptionID {
			filtered = append(filtered, node)
		}
	}
	a.data.Nodes = append(filtered, nodes...)
	target.LastUpdated = time.Now().UnixMilli()
	target.LastError = ""
	target.LastChecksum = checksum
	target.NodeCount = len(nodes)
	a.lastErr = ""
	a.lastSeen = time.Now().UnixMilli()
	return a.store.Save(a.data)
}

func (a *App) RefreshPublicIP() (string, error) {
	ip, err := fetchPublicIP(a.ctx)
	a.mu.Lock()
	defer a.mu.Unlock()
	if err != nil {
		a.lastErr = err.Error()
		return "", err
	}
	a.publicIP = ip
	a.lastErr = ""
	a.lastSeen = time.Now().UnixMilli()
	return ip, nil
}

func (a *App) Connect() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	profile, node, err := a.activeProfileNode()
	if err != nil {
		a.lastErr = err.Error()
		return err
	}

	if a.manager == nil {
		return fmt.Errorf("sing-box manager is not initialized")
	}
	if a.data.Settings.SingBoxPath == "" {
		return fmt.Errorf("sing-box path is not configured")
	}
	configPath, err := paths.SingBoxConfigPath()
	if err != nil {
		return err
	}
	logPath, err := paths.SingBoxLogPath()
	if err != nil {
		return err
	}

	config, err := singbox.BuildConfig(profile, node, logPath)
	if err != nil {
		a.lastErr = err.Error()
		return err
	}
	if err := singbox.WriteConfig(configPath, config); err != nil {
		a.lastErr = err.Error()
		return err
	}

	if err := a.manager.Restart(a.ctx); err != nil {
		a.lastErr = err.Error()
		return err
	}

	a.lastErr = ""
	a.lastSeen = time.Now().UnixMilli()
	return nil
}

func (a *App) Disconnect() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.manager == nil {
		return fmt.Errorf("sing-box manager is not initialized")
	}
	if err := a.manager.Stop(); err != nil {
		a.lastErr = err.Error()
		return err
	}
	a.lastErr = ""
	a.lastSeen = time.Now().UnixMilli()
	return nil
}

func (a *App) Shutdown(ctx context.Context) {
	a.mu.Lock()
	manager := a.manager
	a.mu.Unlock()
	if manager != nil {
		_ = manager.Stop()
	}
}

func (a *App) activeProfileNode() (model.Profile, model.Node, error) {
	if a.data.ActiveProfile == "" {
		return model.Profile{}, model.Node{}, ErrProfileNotFound
	}
	var profile model.Profile
	found := false
	for _, item := range a.data.Profiles {
		if item.ID == a.data.ActiveProfile {
			profile = item
			found = true
			break
		}
	}
	if !found {
		return model.Profile{}, model.Node{}, ErrProfileNotFound
	}
	for _, node := range a.data.Nodes {
		if node.ID == profile.NodeID {
			return profile, node, nil
		}
	}
	return profile, model.Node{}, ErrNodeNotFound
}

func (a *App) ensureManualSubscription() model.Subscription {
	for _, subscription := range a.data.Subscriptions {
		if subscription.Name == manualSubscriptionName {
			return subscription
		}
	}
	subscription := model.Subscription{
		ID:         newID(),
		Name:       manualSubscriptionName,
		AutoUpdate: "off",
	}
	a.data.Subscriptions = append(a.data.Subscriptions, subscription)
	return subscription
}
