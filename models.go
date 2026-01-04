package tgo

// Capability defines a plugin's extension point.
type Capability struct {
	Type      string   `json:"type"`
	Title     string   `json:"title"`
	Icon      string   `json:"icon,omitempty"`
	Priority  int      `json:"priority,omitempty"`
	Tooltip   string   `json:"tooltip,omitempty"`
	Shortcut  string   `json:"shortcut,omitempty"`
	URL       string   `json:"url,omitempty"`
	Width     int      `json:"width,omitempty"`
	RefreshOn []string `json:"refresh_on,omitempty"`
}

// CapabilityOption is a function to configure a Capability.
type CapabilityOption func(*Capability)

func WithIcon(icon string) CapabilityOption {
	return func(c *Capability) { c.Icon = icon }
}

func WithPriority(p int) CapabilityOption {
	return func(c *Capability) { c.Priority = p }
}

func WithTooltip(t string) CapabilityOption {
	return func(c *Capability) { c.Tooltip = t }
}

func WithURL(u string) CapabilityOption {
	return func(c *Capability) { c.URL = u }
}

func WithWidth(w int) CapabilityOption {
	return func(c *Capability) { c.Width = w }
}

// VisitorPanel creates a visitor_panel capability.
func VisitorPanel(title string, opts ...CapabilityOption) Capability {
	c := Capability{Type: "visitor_panel", Title: title, Priority: 10}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// ChatToolbar creates a chat_toolbar capability.
func ChatToolbar(title string, opts ...CapabilityOption) Capability {
	c := Capability{Type: "chat_toolbar", Title: title}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// SidebarIframe creates a sidebar_iframe capability.
func SidebarIframe(title string, url string, opts ...CapabilityOption) Capability {
	c := Capability{Type: "sidebar_iframe", Title: title, URL: url, Width: 400}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// Visitor contains information about a visitor.
type Visitor struct {
	ID             string         `json:"id"`
	PlatformOpenID string         `json:"platform_open_id,omitempty"`
	Name           string         `json:"name,omitempty"`
	Email          string         `json:"email,omitempty"`
	Phone          string         `json:"phone,omitempty"`
	Avatar         string         `json:"avatar,omitempty"`
	Metadata       map[string]any `json:"metadata,omitempty"`
}

// RenderContext is provided to render handlers.
type RenderContext struct {
	VisitorID string         `json:"visitor_id"`
	SessionID string         `json:"session_id,omitempty"`
	Visitor   *Visitor       `json:"visitor,omitempty"`
	AgentID   string         `json:"agent_id,omitempty"`
	ActionID  string         `json:"action_id,omitempty"`
	Language  string         `json:"language,omitempty"`
	Context   map[string]any `json:"context"`
}

// EventContext is provided to event handlers.
type EventContext struct {
	EventType  string         `json:"event_type"`
	ActionID   string         `json:"action_id"`
	VisitorID  string         `json:"visitor_id,omitempty"`
	SessionID  string         `json:"session_id,omitempty"`
	SelectedID string         `json:"selected_id,omitempty"`
	Language   string         `json:"language,omitempty"`
	FormData   map[string]any `json:"form_data,omitempty"`
	Payload    map[string]any `json:"payload"`
}
