package app

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"g1sti/internal/model"
	"g1sti/internal/paths"
	"g1sti/internal/singbox"
	"g1sti/internal/storage"
	"g1sti/internal/subscription"
)

var ErrProfileNotFound = errors.New("profile not found")
var ErrNodeNotFound = errors.New("node not found")

type Status struct {
	Connected     bool      `json:"connected"`
	ActiveProfile string    `json:"active_profile"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastError     string    `json:"last_error"`
}

type App struct {
	mu       sync.Mutex
	ctx      context.Context
	store    *storage.Store
	service  *subscription.Service
	manager  *singbox.Manager
	data     model.AppData
	lastErr  string
	lastSeen time.Time
}

func New() *App {
	return &App{
		service: subscription.NewService(nil),
		data:    model.DefaultAppData(),
	}
}

func (a *App) Startup(ctx context.Context) error {
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
		UpdatedAt:     a.lastSeen,
		LastError:     a.lastErr,
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

func (a *App) SetActiveProfile(profileID string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, profile := range a.data.Profiles {
		if profile.ID == profileID {
			a.data.ActiveProfile = profileID
			a.lastSeen = time.Now()
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
	target.LastUpdated = time.Now()
	target.LastError = ""
	target.LastChecksum = checksum
	target.NodeCount = len(nodes)
	a.lastErr = ""
	a.lastSeen = time.Now()
	return a.store.Save(a.data)
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
	a.lastSeen = time.Now()
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
	a.lastSeen = time.Now()
	return nil
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
