package tlsscanner

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "tls"

type Scanner struct{}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)

	if target.Scheme != "https" {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "tls-http-target",
			Title:          "Target URL does not use HTTPS",
			Description:    "The submitted target was an HTTP URL, so the scanner did not negotiate TLS for that endpoint.",
			Severity:       "high",
			Recommendation: "Use HTTPS for production applications and redirect HTTP traffic to HTTPS.",
		})
		result.Metadata["checked"] = false
		scanner.Complete(&result, 45)
		return result
	}

	if len(target.IPs) == 0 {
		return scanner.Fail(moduleName, result.StartedAt, fmt.Errorf("no resolved IPs available"))
	}

	versions := []uint16{tls.VersionTLS10, tls.VersionTLS11, tls.VersionTLS12, tls.VersionTLS13}
	supported := make([]string, 0)
	var firstState *tls.ConnectionState
	var verifyErr error

	for _, version := range versions {
		state, err := handshake(ctx, target, target.IPs[0], version)
		if err == nil {
			versionName := tlsVersionName(version)
			supported = append(supported, versionName)
			if firstState == nil {
				copied := state
				firstState = &copied
				if len(state.PeerCertificates) > 0 {
					verifyErr = verifyCertificate(state.PeerCertificates[0], target.Host)
				}
			}
		}
	}

	result.Metadata["supported_versions"] = supported

	score := 100
	if len(supported) == 0 {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "tls-handshake-failed",
			Title:          "TLS handshake failed",
			Description:    "The scanner could not establish a TLS connection to the HTTPS endpoint.",
			Severity:       "high",
			Recommendation: "Confirm the service is reachable on HTTPS and that the certificate chain is valid.",
		})
		scanner.Complete(&result, 35)
		return result
	}

	if contains(supported, "TLS 1.0") || contains(supported, "TLS 1.1") {
		score -= 25
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "tls-legacy-protocols",
			Title:          "Legacy TLS protocol support detected",
			Description:    "The endpoint negotiated TLS 1.0 or TLS 1.1 during testing.",
			Severity:       "medium",
			Evidence:       fmt.Sprintf("Supported versions: %v", supported),
			Recommendation: "Disable TLS 1.0 and TLS 1.1 and require TLS 1.2 or newer.",
		})
	}
	if !contains(supported, "TLS 1.3") {
		score -= 5
	}

	if firstState != nil && len(firstState.PeerCertificates) > 0 {
		cert := firstState.PeerCertificates[0]
		now := time.Now()
		result.Metadata["certificate"] = map[string]any{
			"subject":       cert.Subject.String(),
			"issuer":        cert.Issuer.String(),
			"not_before":    cert.NotBefore,
			"not_after":     cert.NotAfter,
			"dns_names":     cert.DNSNames,
			"serial_number": cert.SerialNumber.String(),
		}

		if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
			score -= 35
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "tls-certificate-invalid-date",
				Title:          "Certificate is outside its validity period",
				Description:    "The leaf certificate is either expired or not valid yet.",
				Severity:       "high",
				Evidence:       fmt.Sprintf("Valid from %s to %s", cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339)),
				Recommendation: "Renew or replace the certificate with a currently valid certificate.",
			})
		} else if time.Until(cert.NotAfter) < 30*24*time.Hour {
			score -= 10
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "tls-certificate-expiring",
				Title:          "Certificate expires soon",
				Description:    "The leaf certificate expires in less than 30 days.",
				Severity:       "low",
				Evidence:       cert.NotAfter.Format(time.RFC3339),
				Recommendation: "Renew the certificate before expiration to avoid service disruption.",
			})
		}

		if verifyErr != nil {
			score -= 20
			result.Findings = append(result.Findings, domain.Finding{
				ID:             "tls-certificate-verification",
				Title:          "Certificate verification warning",
				Description:    "The certificate could not be verified against the host using the system trust store.",
				Severity:       "medium",
				Evidence:       verifyErr.Error(),
				Recommendation: "Use a certificate issued by a trusted CA for the scanned hostname.",
			})
		}
	}

	scanner.Complete(&result, score)
	return result
}

func handshake(ctx context.Context, target utils.ValidatedTarget, ip net.IP, version uint16) (tls.ConnectionState, error) {
	handshakeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dialer := net.Dialer{Timeout: 4 * time.Second}
	conn, err := dialer.DialContext(handshakeCtx, "tcp", net.JoinHostPort(ip.String(), target.Port))
	if err != nil {
		return tls.ConnectionState{}, err
	}
	defer conn.Close()

	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         target.Host,
		MinVersion:         version,
		MaxVersion:         version,
		InsecureSkipVerify: true,
	})
	if err := tlsConn.HandshakeContext(handshakeCtx); err != nil {
		return tls.ConnectionState{}, err
	}
	return tlsConn.ConnectionState(), nil
}

func verifyCertificate(cert *x509.Certificate, host string) error {
	_, err := cert.Verify(x509.VerifyOptions{
		DNSName:     host,
		CurrentTime: time.Now(),
	})
	return err
}

func tlsVersionName(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return fmt.Sprintf("0x%x", version)
	}
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
