package utils

import (
	"context"
	"net"
	"testing"
)

func TestValidateTargetURLRejectsUnsafeLocalTargets(t *testing.T) {
	tests := []string{
		"http://localhost",
		"http://localhost:8080",
		"http://127.0.0.1",
		"http://10.1.2.3",
		"http://172.16.10.20",
		"http://192.168.1.10",
		"http://169.254.1.1",
		"http://[::1]",
		"http://[fc00::1]",
	}

	for _, rawURL := range tests {
		if _, err := ValidateTargetURL(context.Background(), rawURL); err == nil {
			t.Fatalf("expected %s to be rejected", rawURL)
		}
	}
}

func TestIsUnsafeIPAllowsPublicAddress(t *testing.T) {
	if IsUnsafeIP(net.ParseIP("8.8.8.8")) {
		t.Fatal("expected public IPv4 address to be allowed")
	}
	if IsUnsafeIP(net.ParseIP("2001:4860:4860::8888")) {
		t.Fatal("expected public IPv6 address to be allowed")
	}
}

func TestScoreToGrade(t *testing.T) {
	tests := map[int]string{
		100: "A+",
		93:  "A",
		88:  "B+",
		75:  "C",
		61:  "D-",
		12:  "F",
	}

	for score, want := range tests {
		if got := ScoreToGrade(score); got != want {
			t.Fatalf("ScoreToGrade(%d) = %s, want %s", score, got, want)
		}
	}
}
