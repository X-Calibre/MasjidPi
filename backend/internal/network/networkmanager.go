package network

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const commandTimeout = 25 * time.Second

type commandRunner interface {
	Run(context.Context, []string, string) ([]byte, error)
}

type execRunner struct{}

func (execRunner) Run(ctx context.Context, args []string, stdin string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "nmcli", args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	return cmd.CombinedOutput()
}

type NetworkManager struct {
	runner    commandRunner
	available func() bool
}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{
		runner: execRunner{},
		available: func() bool {
			_, err := exec.LookPath("nmcli")
			return err == nil
		},
	}
}

func newNetworkManager(runner commandRunner) *NetworkManager {
	return &NetworkManager{runner: runner, available: func() bool { return true }}
}

func (m *NetworkManager) Status(ctx context.Context) (WiFiStatus, error) {
	if !m.available() {
		return WiFiStatus{}, nil
	}

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	profiles, err := m.runner.Run(ctx, []string{"-t", "--escape", "yes", "-f", "TYPE", "connection", "show"}, "")
	if err != nil {
		return WiFiStatus{Supported: true}, fmt.Errorf("list NetworkManager profiles: %w", err)
	}
	configured := false
	for _, line := range lines(profiles) {
		if line == "802-11-wireless" || line == "wifi" {
			configured = true
			break
		}
	}

	active, err := m.runner.Run(ctx, []string{"-t", "--escape", "yes", "-f", "TYPE,NAME", "connection", "show", "--active"}, "")
	if err != nil {
		return WiFiStatus{Supported: true, Configured: configured}, fmt.Errorf("list active NetworkManager profiles: %w", err)
	}
	status := WiFiStatus{Supported: true, Configured: configured}
	for _, line := range lines(active) {
		fields := splitEscaped(line, ':')
		if len(fields) == 2 && (fields[0] == "802-11-wireless" || fields[0] == "wifi") {
			status.Connected = true
			status.SSID = fields[1]
			break
		}
	}
	return status, nil
}

func (m *NetworkManager) Scan(ctx context.Context) ([]WiFiNetwork, error) {
	if !m.available() {
		return nil, errors.New("NetworkManager is unavailable")
	}
	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	out, err := m.runner.Run(ctx, []string{"-t", "--escape", "yes", "-f", "IN-USE,SSID,SIGNAL,SECURITY", "device", "wifi", "list", "--rescan", "yes"}, "")
	if err != nil {
		return nil, fmt.Errorf("scan Wi-Fi networks: %w", err)
	}

	bySSID := make(map[string]WiFiNetwork)
	for _, line := range lines(out) {
		fields := splitEscaped(line, ':')
		if len(fields) != 4 || strings.TrimSpace(fields[1]) == "" {
			continue
		}
		signal, _ := strconv.Atoi(fields[2])
		candidate := WiFiNetwork{
			SSID:     fields[1],
			Signal:   signal,
			Security: normaliseSecurity(fields[3]),
			Active:   fields[0] == "*" || strings.EqualFold(fields[0], "yes"),
		}
		if current, ok := bySSID[candidate.SSID]; !ok || candidate.Signal > current.Signal || candidate.Active {
			bySSID[candidate.SSID] = candidate
		}
	}

	networks := make([]WiFiNetwork, 0, len(bySSID))
	for _, item := range bySSID {
		networks = append(networks, item)
	}
	sort.Slice(networks, func(i, j int) bool {
		if networks[i].Active != networks[j].Active {
			return networks[i].Active
		}
		if networks[i].Signal != networks[j].Signal {
			return networks[i].Signal > networks[j].Signal
		}
		return strings.ToLower(networks[i].SSID) < strings.ToLower(networks[j].SSID)
	})
	return networks, nil
}

func (m *NetworkManager) Connect(ctx context.Context, ssid, password string) error {
	ssid = strings.TrimSpace(ssid)
	if ssid == "" {
		return errors.New("Wi-Fi network name is required")
	}
	if strings.ContainsAny(ssid, "\r\n\x00") || strings.ContainsAny(password, "\r\n\x00") {
		return errors.New("Wi-Fi credentials contain unsupported characters")
	}
	if !m.available() {
		return errors.New("NetworkManager is unavailable")
	}

	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()
	stdin := password + "\n"
	if password == "" {
		stdin = ""
	}
	// --ask reads the secret from stdin, keeping it out of the process argument
	// list. Command output is deliberately discarded so a secret can never be
	// included in an API response or application log.
	_, err := m.runner.Run(ctx, []string{"--ask", "--wait", "20", "device", "wifi", "connect", ssid}, stdin)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errors.New("the Wi-Fi connection attempt timed out")
		}
		return errors.New("could not connect; check the Wi-Fi password and try again")
	}
	return nil
}

func lines(value []byte) []string {
	trimmed := bytes.TrimSpace(value)
	if len(trimmed) == 0 {
		return nil
	}
	return strings.Split(string(trimmed), "\n")
}

func splitEscaped(value string, separator rune) []string {
	var fields []string
	var current strings.Builder
	escaped := false
	for _, char := range value {
		if escaped {
			current.WriteRune(char)
			escaped = false
			continue
		}
		if char == '\\' {
			escaped = true
			continue
		}
		if char == separator {
			fields = append(fields, current.String())
			current.Reset()
			continue
		}
		current.WriteRune(char)
	}
	if escaped {
		current.WriteRune('\\')
	}
	return append(fields, current.String())
}

func normaliseSecurity(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "--" {
		return "Open"
	}
	return value
}
