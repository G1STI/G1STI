package singbox

import (
	"testing"

	"g1sti/internal/model"
)

func TestBuildConfigVLESSReality(t *testing.T) {
	node := model.Node{
		Type:   "vless",
		Server: "example.com",
		Port:   443,
		ParsedFields: map[string]string{
			"uuid":     "123e4567-e89b-12d3-a456-426614174000",
			"security": "reality",
			"sni":      "example.com",
			"pbk":      "public-key",
			"sid":      "shortid",
			"fp":       "chrome",
		},
	}
	profile := model.Profile{TunEnabled: true}

	config, err := BuildConfig(profile, node, "singbox.log")
	if err != nil {
		t.Fatalf("BuildConfig returned error: %v", err)
	}
	if len(config.Outbounds) == 0 {
		t.Fatalf("expected outbounds")
	}
	proxy := config.Outbounds[0]
	if proxy.Type != "vless" {
		t.Fatalf("expected vless outbound, got %q", proxy.Type)
	}
	if proxy.TLS == nil {
		t.Fatalf("expected tls settings")
	}
	if proxy.TLS["reality"] == nil {
		t.Fatalf("expected reality settings")
	}
}
