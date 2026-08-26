package livestatus

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Status represents the last availability state received for a LiveMasjid mount.
type Status struct {
	Mount     string
	Available bool
	Payload   string
	UpdatedAt time.Time
}

// Client subscribes to LiveMasjid's mount status feed. It does not poll the
// LiveMasjid HTTP site. The MQTT connection is long-lived and subscriptions are
// restored automatically after reconnects.
type Client struct {
	broker string
	port   int
	log    *slog.Logger

	mu     sync.RWMutex
	mounts map[string]Status
	wake   chan string
	client mqtt.Client
}

func New(broker string, port int, log *slog.Logger) *Client {
	return &Client{
		broker: broker,
		port:   port,
		log:    log,
		mounts: make(map[string]Status),
		wake:   make(chan string, 16),
	}
}

func (c *Client) Start(ctx context.Context) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.broker, c.port))
	opts.SetClientID(fmt.Sprintf("masjidpi-%d", time.Now().UnixNano()))
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(10 * time.Second)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		token := client.Subscribe("mounts/#", 0, c.onMessage)
		if token.Wait() && token.Error() != nil {
			c.log.Error("LiveMasjid MQTT subscription failed", "error", token.Error())
			return
		}
		c.log.Info("Connected to LiveMasjid status feed")
	})
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		c.log.Warn("LiveMasjid status feed disconnected", "error", err)
	})

	c.client = mqtt.NewClient(opts)
	go func() {
		if token := c.client.Connect(); token.Wait() && token.Error() != nil {
			c.log.Warn("LiveMasjid MQTT connection unavailable; will retry", "error", token.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		c.Close()
	}()
}

func (c *Client) onMessage(_ mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) < 2 || parts[1] == "" {
		return
	}

	mount := parts[1]
	payload := string(msg.Payload())

	// Individual MQTT events can be frequent during normal operation. Keep
	// them at DEBUG so routine status traffic does not fill the persistent
	// system journal. Connection/subscription failures remain WARN/ERROR.
	c.log.Debug(
		"LiveMasjid MQTT event received",
		"topic", msg.Topic(),
		"mount", mount,
		"payload", payload,
	)

	lower := strings.ToLower(payload)

	available := strings.Contains(lower, "started")
	if !available && !strings.Contains(lower, "stopped") {
		return
	}

	status := Status{
		Mount:     mount,
		Available: available,
		Payload:   payload,
		UpdatedAt: time.Now(),
	}

	c.mu.Lock()
	c.mounts[mount] = status
	c.mu.Unlock()

	select {
	case c.wake <- mount:
	default:
	}
}

func (c *Client) Status(mount string) (bool, bool, time.Time) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status, ok := c.mounts[mount]
	if !ok {
		return false, false, time.Time{}
	}

	return status.Available, true, status.UpdatedAt
}

func (c *Client) Events() <-chan string { return c.wake }

func (c *Client) Close() {
	if c.client != nil && c.client.IsConnected() {
		c.client.Disconnect(250)
	}
}
