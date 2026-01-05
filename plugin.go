package tgo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Plugin is the interface that all TGO plugins must implement.
type Plugin interface {
	ID() string
	Name() string
	Version() string
	Capabilities() []Capability
}

// Handler methods that a plugin can optionally implement.
type VisitorPanelRenderer interface {
	OnVisitorPanelRender(ctx *RenderContext) Template
}
type VisitorPanelEventHandler interface {
	OnVisitorPanelEvent(ctx *EventContext) *Action
}
type ChatToolbarRenderer interface {
	OnChatToolbarRender(ctx *RenderContext) Template
}
type ChatToolbarEventHandler interface {
	OnChatToolbarEvent(ctx *EventContext) *Action
}
type SidebarIframeConfigurator interface {
	OnSidebarIframeConfig(params map[string]any) any
}
type ChannelIntegrationManifestProvider interface {
	OnChannelIntegrationManifest(params map[string]any) any
}
type ToolHandler interface {
	OnToolExecute(ctx *ToolContext, toolName string, args map[string]any) (*ToolResult, error)
}

// Options for running a plugin.
type Options struct {
	SocketPath string
	TCPAddr    string
	DevToken   string
}

type Option func(*Options)

func WithSocketPath(path string) Option {
	return func(o *Options) { o.SocketPath = path }
}

func WithTCPAddr(addr string) Option {
	return func(o *Options) { o.TCPAddr = addr }
}

func WithDevToken(token string) Option {
	return func(o *Options) { o.DevToken = token }
}

// Run starts the plugin and handles communication with TGO.
func Run(p Plugin, opts ...Option) error {
	options := &Options{
		SocketPath: "/var/run/tgo/tgo.sock",
	}
	for _, opt := range opts {
		opt(options)
	}

	var transport *Transport
	if options.TCPAddr != "" {
		transport = NewTCPTransport(options.TCPAddr)
	} else {
		transport = NewUnixTransport(options.SocketPath)
	}

	if err := transport.Connect(); err != nil {
		return err
	}
	defer transport.Close()

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Register the plugin
	if err := register(p, transport, options.DevToken); err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	log.Printf("Plugin '%s' v%s is running", p.Name(), p.Version())

	// Main request loop
	done := make(chan error, 1)
	go func() {
		for {
			msg, err := transport.RecvMessage()
			if err != nil {
				done <- err
				return
			}

			go handleRequest(p, transport, msg)
		}
	}()

	select {
	case err := <-done:
		log.Printf("Connection lost: %v", err)
		return err
	case sig := <-sigChan:
		log.Printf("Received signal %v, shutting down...", sig)
		return nil
	}
}

func register(p Plugin, t *Transport, devToken string) error {
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "register",
		"params": map[string]any{
			"id":           p.ID(),
			"name":         p.Name(),
			"version":      p.Version(),
			"capabilities": p.Capabilities(),
			"dev_token":    devToken,
		},
	}

	if err := t.SendMessage(req); err != nil {
		return err
	}

	resp, err := t.RecvMessage()
	if err != nil {
		return err
	}

	result, ok := resp["result"].(map[string]any)
	if !ok || result["success"] != true {
		return fmt.Errorf("registration failed: %v", resp["error"])
	}

	return nil
}

func handleRequest(p Plugin, t *Transport, msg map[string]any) {
	method, _ := msg["method"].(string)
	id, _ := msg["id"]
	params, _ := msg["params"].(map[string]any)

	if method == "" {
		return
	}

	if method == "shutdown" {
		t.SendMessage(map[string]any{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  map[string]any{"success": true},
		})
		return
	}

	if method == "ping" {
		t.SendMessage(map[string]any{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  map[string]any{"pong": true},
		})
		return
	}

	var result any
	var err error

	switch method {
	case "visitor_panel/render":
		if h, ok := p.(VisitorPanelRenderer); ok {
			ctx := &RenderContext{}
			mapToStruct(params, ctx)
			result = h.OnVisitorPanelRender(ctx)
		}
	case "visitor_panel/event":
		if h, ok := p.(VisitorPanelEventHandler); ok {
			ctx := &EventContext{}
			mapToStruct(params, ctx)
			result = h.OnVisitorPanelEvent(ctx)
		}
	case "chat_toolbar/render":
		if h, ok := p.(ChatToolbarRenderer); ok {
			ctx := &RenderContext{}
			mapToStruct(params, ctx)
			result = h.OnChatToolbarRender(ctx)
		}
	case "chat_toolbar/event":
		if h, ok := p.(ChatToolbarEventHandler); ok {
			ctx := &EventContext{}
			mapToStruct(params, ctx)
			result = h.OnChatToolbarEvent(ctx)
		}
	case "sidebar_iframe/config":
		if h, ok := p.(SidebarIframeConfigurator); ok {
			result = h.OnSidebarIframeConfig(params)
		}
	case "channel_integration/manifest":
		if h, ok := p.(ChannelIntegrationManifestProvider); ok {
			result = h.OnChannelIntegrationManifest(params)
		}
	case "tool/execute":
		if h, ok := p.(ToolHandler); ok {
			ctx := &ToolContext{}
			mapToStruct(params, ctx)
			toolName, _ := params["tool_name"].(string)
			args, _ := params["arguments"].(map[string]any)
			result, err = h.OnToolExecute(ctx, toolName, args)
		}
	default:
		err = fmt.Errorf("method not found: %s", method)
	}

	if err != nil {
		t.SendMessage(map[string]any{
			"jsonrpc": "2.0",
			"id":      id,
			"error":   map[string]any{"code": -32601, "message": err.Error()},
		})
		return
	}

	// If no handler was implemented but method exists
	if result == nil {
		t.SendMessage(map[string]any{
			"jsonrpc": "2.0",
			"id":      id,
			"result":  map[string]any{"success": true},
		})
		return
	}

	// Unwrap potential builders
	if b, ok := result.(interface{ ToMap() map[string]any }); ok {
		result = b.ToMap()
	}

	t.SendMessage(map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	})
}

// Helper to convert map[string]any to struct via JSON (simple approach)
func mapToStruct(m map[string]any, s any) {
	data, _ := json.Marshal(m)
	json.Unmarshal(data, s)
}

// BasePlugin provides default implementations for Plugin interface.
type BasePlugin struct {
	PID      string
	PName    string
	PVersion string
	Caps     []Capability
}

func (b *BasePlugin) ID() string                 { return b.PID }
func (b *BasePlugin) Name() string               { return b.PName }
func (b *BasePlugin) Version() string            { return b.PVersion }
func (b *BasePlugin) Capabilities() []Capability { return b.Caps }
