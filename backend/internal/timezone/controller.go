package timezone

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidFrame/backend/internal/atomicfile"
)

const (
	defaultZoneTable = "/usr/share/zoneinfo/zone.tab"
	commandTimeout   = 10 * time.Second
)

var countryCodes = map[string]string{
	"Australia":     "AU",
	"Bangladesh":    "BD",
	"Barbados":      "BB",
	"Botswana":      "BW",
	"Canada":        "CA",
	"England":       "GB",
	"Eswatini":      "SZ",
	"Ethiopia":      "ET",
	"Germany":       "DE",
	"Grenada":       "GD",
	"India":         "IN",
	"Lebanon":       "LB",
	"Malawi":        "MW",
	"Mozambique":    "MZ",
	"Namibia":       "NA",
	"New Zealand":   "NZ",
	"Nigeria":       "NG",
	"Pakistan":      "PK",
	"Panama":        "PA",
	"Saudi Arabia":  "SA",
	"South Africa":  "ZA",
	"United States": "US",
	"Zambia":        "ZM",
	"Zimbabwe":      "ZW",
}

type Zone struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type persistedState struct {
	Name string `json:"name"`
}

type commandRunner interface {
	SetTimezone(context.Context, string) error
}

type timedatectlRunner struct{}

func (timedatectlRunner) SetTimezone(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, commandTimeout)
	defer cancel()

	output, err := exec.CommandContext(
		ctx,
		"timedatectl",
		"set-timezone",
		name,
	).CombinedOutput()
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return errors.New("setting the timezone timed out")
		}
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return fmt.Errorf("set timezone: %s", message)
	}
	return nil
}

type Controller struct {
	statePath string
	zoneTable string
	runner    commandRunner
}

func NewController(statePath string) *Controller {
	return &Controller{
		statePath: statePath,
		zoneTable: defaultZoneTable,
		runner:    timedatectlRunner{},
	}
}

func newController(
	statePath string,
	zoneTable string,
	runner commandRunner,
) *Controller {
	return &Controller{
		statePath: statePath,
		zoneTable: zoneTable,
		runner:    runner,
	}
}

func (c *Controller) Zones(country string) ([]Zone, error) {
	code, ok := countryCodes[strings.TrimSpace(country)]
	if !ok {
		return nil, errors.New("unsupported country")
	}

	zones, err := c.readZones()
	if err != nil {
		return nil, err
	}

	result := append([]Zone(nil), zones[code]...)
	if len(result) == 0 {
		return nil, fmt.Errorf("no timezones available for %s", country)
	}
	return result, nil
}

func (c *Controller) Current() (string, bool, error) {
	data, err := os.ReadFile(c.statePath)
	if errors.Is(err, os.ErrNotExist) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("read timezone state: %w", err)
	}

	var state persistedState
	if err := json.Unmarshal(data, &state); err != nil {
		return "", false, fmt.Errorf("decode timezone state: %w", err)
	}
	state.Name = strings.TrimSpace(state.Name)
	if state.Name == "" {
		return "", false, errors.New("persisted timezone is empty")
	}
	if err := c.validate(state.Name); err != nil {
		return "", false, fmt.Errorf("invalid persisted timezone: %w", err)
	}
	return state.Name, true, nil
}

func (c *Controller) Set(ctx context.Context, name string) error {
	name = strings.TrimSpace(name)
	if err := c.validate(name); err != nil {
		return err
	}
	if err := c.runner.SetTimezone(ctx, name); err != nil {
		return err
	}
	if err := atomicfile.WriteJSON(
		c.statePath,
		persistedState{Name: name},
		0644,
	); err != nil {
		return fmt.Errorf("persist timezone: %w", err)
	}
	return nil
}

func (c *Controller) Restore(ctx context.Context) (bool, error) {
	name, ok, err := c.Current()
	if err != nil || !ok {
		return false, err
	}
	if err := c.runner.SetTimezone(ctx, name); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Controller) validate(name string) error {
	if name == "" {
		return errors.New("timezone is required")
	}
	zones, err := c.readZones()
	if err != nil {
		return err
	}
	for _, countryZones := range zones {
		for _, zone := range countryZones {
			if zone.Name == name {
				return nil
			}
		}
	}
	return errors.New("unsupported timezone")
}

func (c *Controller) readZones() (map[string][]Zone, error) {
	file, err := os.Open(c.zoneTable)
	if err != nil {
		return nil, fmt.Errorf("open timezone database: %w", err)
	}
	defer file.Close()

	zones := make(map[string][]Zone)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) < 3 {
			continue
		}
		code := strings.TrimSpace(fields[0])
		name := strings.TrimSpace(fields[2])
		if len(code) != 2 || name == "" {
			continue
		}
		description := ""
		if len(fields) >= 4 {
			description = strings.TrimSpace(fields[3])
		}
		zones[code] = append(zones[code], Zone{
			Name:        name,
			Description: description,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read timezone database: %w", err)
	}

	for code := range zones {
		sort.Slice(zones[code], func(i, j int) bool {
			return zones[code][i].Name < zones[code][j].Name
		})
	}
	return zones, nil
}
