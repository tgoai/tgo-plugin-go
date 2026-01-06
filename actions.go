package tgo

// Action represents an action instruction for the TGO host.
type Action struct {
	Type string         `json:"action"`
	Data map[string]any `json:"data,omitempty"`
	next *Action        // Private, used for chaining
}

// Then adds another action to be executed after this one.
func (a *Action) Then(next *Action) *Action {
	if a == nil || next == nil {
		return a
	}
	if a.next == nil {
		a.next = next
	} else {
		a.next.Then(next)
	}
	return a
}

func (a *Action) ToMap() map[string]any {
	if a.next == nil {
		return map[string]any{
			"action": a.Type,
			"data":   a.Data,
		}
	}

	// Chained actions, convert to batch
	actions := []map[string]any{
		{"action": a.Type, "data": a.Data},
	}
	for curr := a.next; curr != nil; curr = curr.next {
		actions = append(actions, map[string]any{
			"action": curr.Type,
			"data":   curr.Data,
		})
	}

	return map[string]any{
		"action": "batch",
		"data": map[string]any{
			"actions": actions,
		},
	}
}

// OpenURL opens a URL in the user's browser.
func OpenURL(url, target string) *Action {
	return &Action{
		Type: "open_url",
		Data: map[string]any{"url": url, "target": target},
	}
}

// InsertText inserts text into the agent's input field.
func InsertText(text string, replace bool) *Action {
	return &Action{
		Type: "insert_text",
		Data: map[string]any{"text": text, "replace": replace},
	}
}

// SendMessage sends a message to the visitor.
func SendMessage(content, contentType string) *Action {
	return &Action{
		Type: "send_message",
		Data: map[string]any{"content": content, "content_type": contentType},
	}
}

// ShowToast displays a notification toast.
func ShowToast(message, tp string) *Action {
	return &Action{
		Type: "show_toast",
		Data: map[string]any{"message": message, "type": tp, "duration": 3000},
	}
}

// CopyText copies text to the clipboard.
func CopyText(text, toast string) *Action {
	return &Action{
		Type: "copy_text",
		Data: map[string]any{"text": text, "toast": toast},
	}
}

// ShowModal shows a modal with UI template.
func ShowModal(title string, t Template) *Action {
	m := t.ToMap()
	data := map[string]any{
		"title":    title,
		"template": m["template"],
		"data":     m["data"],
	}
	return &Action{
		Type: "show_modal",
		Data: data,
	}
}

// Refresh re-renders the current plugin UI.
func Refresh() *Action {
	return &Action{Type: "refresh"}
}

// CloseModal closes the currently open modal.
func CloseModal() *Action {
	return &Action{Type: "close_modal"}
}

// Noop performs no operation.
func Noop() *Action {
	return &Action{Type: "noop"}
}

