package ports

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/secscan/secscan/backend/internal/domain"
	"github.com/secscan/secscan/backend/internal/scanner"
	"github.com/secscan/secscan/backend/internal/utils"
)

const moduleName = "ports"

type Scanner struct{}

type probeJob struct {
	IP   net.IP
	Port int
}

type openPort struct {
	IP      string `json:"ip"`
	Port    int    `json:"port"`
	Service string `json:"service"`
}

func New() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Name() string {
	return moduleName
}

func (s *Scanner) Scan(ctx context.Context, target utils.ValidatedTarget) domain.ModuleResult {
	result := scanner.NewResult(moduleName)

	ports := targetPorts(target)
	jobs := make(chan probeJob)
	found := make(chan openPort, len(ports)*len(target.IPs))
	workerCount := 16
	if len(ports) < workerCount {
		workerCount = len(ports)
	}

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				if service, ok := probe(ctx, job.IP, job.Port); ok {
					found <- openPort{IP: job.IP.String(), Port: job.Port, Service: service}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, ip := range target.IPs {
			for _, port := range ports {
				select {
				case <-ctx.Done():
					return
				case jobs <- probeJob{IP: ip, Port: port}:
				}
			}
		}
	}()

	wg.Wait()
	close(found)

	openPorts := make([]openPort, 0)
	for item := range found {
		openPorts = append(openPorts, item)
	}
	sort.Slice(openPorts, func(i, j int) bool {
		if openPorts[i].IP == openPorts[j].IP {
			return openPorts[i].Port < openPorts[j].Port
		}
		return openPorts[i].IP < openPorts[j].IP
	})

	result.Metadata["scanned_ports"] = ports
	result.Metadata["open_ports"] = openPorts
	result.Metadata["ip_count"] = len(target.IPs)

	score := 100
	for _, item := range openPorts {
		if isSensitivePort(item.Port) {
			score -= 20
			result.Findings = append(result.Findings, domain.Finding{
				ID:             fmt.Sprintf("ports-%d", item.Port),
				Title:          "Sensitive service appears reachable",
				Description:    "A service commonly associated with administrative or database access is reachable from the internet.",
				Severity:       "high",
				Evidence:       fmt.Sprintf("%s:%d (%s)", item.IP, item.Port, item.Service),
				Recommendation: "Restrict this service with a firewall, VPN, or allowlist unless it must be public.",
			})
			continue
		}
		if item.Port != 80 && item.Port != 443 {
			score -= 4
		}
	}

	if len(openPorts) == 0 {
		result.Findings = append(result.Findings, domain.Finding{
			ID:             "ports-none-open",
			Title:          "No common TCP ports were reachable",
			Description:    "The scanner could not connect to the small set of common ports it checked.",
			Severity:       "info",
			Recommendation: "Confirm the target URL and network path if this was unexpected.",
		})
	}

	scanner.Complete(&result, score)
	return result
}

func targetPorts(target utils.ValidatedTarget) []int {
	common := []int{80, 443, 8080, 8443, 3000, 5000, 22, 25, 110, 143, 993, 995, 3306, 5432, 6379, 27017}
	seen := make(map[int]struct{}, len(common)+1)
	ports := make([]int, 0, len(common)+1)
	if port, err := strconv.Atoi(target.Port); err == nil && port > 0 && port <= 65535 {
		ports = append(ports, port)
		seen[port] = struct{}{}
	}
	for _, port := range common {
		if _, exists := seen[port]; !exists {
			ports = append(ports, port)
		}
	}
	return ports
}

func probe(ctx context.Context, ip net.IP, port int) (string, bool) {
	probeCtx, cancel := context.WithTimeout(ctx, 900*time.Millisecond)
	defer cancel()

	dialer := net.Dialer{Timeout: 900 * time.Millisecond}
	conn, err := dialer.DialContext(probeCtx, "tcp", net.JoinHostPort(ip.String(), strconv.Itoa(port)))
	if err != nil {
		return "", false
	}
	defer conn.Close()

	return serviceName(port), true
}

func serviceName(port int) string {
	switch port {
	case 22:
		return "ssh"
	case 25:
		return "smtp"
	case 80, 8080, 3000, 5000:
		return "http"
	case 110:
		return "pop3"
	case 143:
		return "imap"
	case 443, 8443:
		return "https"
	case 993:
		return "imaps"
	case 995:
		return "pop3s"
	case 3306:
		return "mysql"
	case 5432:
		return "postgresql"
	case 6379:
		return "redis"
	case 27017:
		return "mongodb"
	default:
		return "tcp"
	}
}

func isSensitivePort(port int) bool {
	switch port {
	case 22, 3306, 5432, 6379, 27017:
		return true
	default:
		return false
	}
}
