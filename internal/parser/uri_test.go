package parser

import "testing"

func TestParseVLESS(t *testing.T) {
	raw := "vless://123e4567-e89b-12d3-a456-426614174000@vless.example.com:443?encryption=none&flow=xtls-rprx-vision&security=reality&sni=example.com&pbk=abcdefg#TestNode"

	node, err := ParseURI(raw)
	if err != nil {
		t.Fatalf("ParseURI returned error: %v", err)
	}
	if node.Type != "vless" {
		t.Fatalf("expected type vless, got %q", node.Type)
	}
	if node.Server != "vless.example.com" {
		t.Fatalf("expected server vless.example.com, got %q", node.Server)
	}
	if node.Port != 443 {
		t.Fatalf("expected port 443, got %d", node.Port)
	}
	if node.Tag != "TestNode" {
		t.Fatalf("expected tag TestNode, got %q", node.Tag)
	}
	if node.ParsedFields["uuid"] != "123e4567-e89b-12d3-a456-426614174000" {
		t.Fatalf("expected uuid parsed field")
	}
}
