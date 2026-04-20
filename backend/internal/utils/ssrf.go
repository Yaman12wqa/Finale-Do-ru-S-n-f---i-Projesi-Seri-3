package utils

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

type ValidatedTarget struct {
	RawURL       string
	URL          *url.URL
	CanonicalURL string
	Scheme       string
	Host         string
	Port         string
	IPs          []net.IP
}

var blockedCIDRs = mustParseCIDRs([]string{
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.168.0.0/16",
	"::/128",
	"::1/128",
	"fc00::/7",
	"fe80::/10",
})

func ValidateTargetURL(ctx context.Context, raw string) (ValidatedTarget, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ValidatedTarget{}, errors.New("url is required")
	}
	if len(trimmed) > 2048 {
		return ValidatedTarget{}, errors.New("url is too long")
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return ValidatedTarget{}, fmt.Errorf("invalid url: %w", err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ValidatedTarget{}, errors.New("only http and https URLs are allowed")
	}
	if parsed.User != nil {
		return ValidatedTarget{}, errors.New("URLs with embedded credentials are not allowed")
	}

	host := parsed.Hostname()
	if host == "" {
		return ValidatedTarget{}, errors.New("url host is required")
	}

	ips, normalizedHost, err := ResolvePublicHost(ctx, host)
	if err != nil {
		return ValidatedTarget{}, err
	}

	port := parsed.Port()
	if port == "" {
		if scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	canonical := *parsed
	canonical.Scheme = scheme
	canonical.Host = normalizedHost
	if parsed.Port() != "" {
		canonical.Host = net.JoinHostPort(normalizedHost, port)
	}
	canonical.Fragment = ""

	return ValidatedTarget{
		RawURL:       trimmed,
		URL:          &canonical,
		CanonicalURL: canonical.String(),
		Scheme:       scheme,
		Host:         normalizedHost,
		Port:         port,
		IPs:          ips,
	}, nil
}

// ResolvePublicHost is the central SSRF guard. Every outbound request path
// resolves DNS and rejects loopback, private, link-local, multicast,
// unspecified, localhost, and other internal-only address ranges before a
// connection is attempted.
func ResolvePublicHost(ctx context.Context, host string) ([]net.IP, string, error) {
	normalizedHost := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(host)), ".")
	if normalizedHost == "" {
		return nil, "", errors.New("host is required")
	}
	if normalizedHost == "localhost" || strings.HasSuffix(normalizedHost, ".localhost") {
		return nil, "", errors.New("localhost targets are blocked")
	}
	if strings.ContainsAny(normalizedHost, " \t\r\n") {
		return nil, "", errors.New("host contains invalid whitespace")
	}

	if ip := net.ParseIP(normalizedHost); ip != nil {
		if IsUnsafeIP(ip) {
			return nil, "", fmt.Errorf("target IP %s is not allowed", ip.String())
		}
		return []net.IP{ip}, normalizedHost, nil
	}

	lookupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	addrs, err := net.DefaultResolver.LookupIPAddr(lookupCtx, normalizedHost)
	if err != nil {
		return nil, "", fmt.Errorf("dns lookup failed: %w", err)
	}
	if len(addrs) == 0 {
		return nil, "", errors.New("dns lookup returned no addresses")
	}

	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		if IsUnsafeIP(addr.IP) {
			return nil, "", fmt.Errorf("resolved IP %s is not allowed", addr.IP.String())
		}
		ips = append(ips, addr.IP)
	}

	return ips, normalizedHost, nil
}

func IsUnsafeIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() || ip.IsInterfaceLocalMulticast() {
		return true
	}
	for _, block := range blockedCIDRs {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func mustParseCIDRs(values []string) []*net.IPNet {
	ranges := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		_, parsed, err := net.ParseCIDR(value)
		if err != nil {
			panic(err)
		}
		ranges = append(ranges, parsed)
	}
	return ranges
}
