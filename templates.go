package tgo

// Template is the interface for all UI templates.
type Template interface {
	ToMap() map[string]any
}

// KeyValue template
type KeyValue struct {
	Title string           `json:"title,omitempty"`
	Items []map[string]any `json:"items"`
}

func NewKeyValue(title string) *KeyValue {
	return &KeyValue{Title: title, Items: []map[string]any{}}
}

func (kv *KeyValue) Add(label string, value any, opts ...KeyValueOption) *KeyValue {
	item := map[string]any{"label": label, "value": value}
	for _, opt := range opts {
		opt(item)
	}
	kv.Items = append(kv.Items, item)
	return kv
}

func (kv *KeyValue) ToMap() map[string]any {
	return map[string]any{
		"template": "key_value",
		"data":     kv,
	}
}

type KeyValueOption func(map[string]any)

func KeyValueIcon(icon string) KeyValueOption {
	return func(m map[string]any) { m["icon"] = icon }
}

func KeyValueColor(color string) KeyValueOption {
	return func(m map[string]any) { m["color"] = color }
}

func KeyValueCopyable(c bool) KeyValueOption {
	return func(m map[string]any) { m["copyable"] = c }
}

// Table template
type Table struct {
	Title      string           `json:"title,omitempty"`
	ColumnsArr []map[string]any `json:"columns"`
	RowsArr    []map[string]any `json:"rows"`
}

func NewTable(title string) *Table {
	return &Table{Title: title, ColumnsArr: []map[string]any{}, RowsArr: []map[string]any{}}
}

func (t *Table) Columns(cols ...any) *Table {
	for _, col := range cols {
		if s, ok := col.(string); ok {
			t.ColumnsArr = append(t.ColumnsArr, map[string]any{"key": s, "label": s})
		} else if m, ok := col.(map[string]any); ok {
			t.ColumnsArr = append(t.ColumnsArr, m)
		}
	}
	return t
}

func (t *Table) Row(row map[string]any) *Table {
	t.RowsArr = append(t.RowsArr, row)
	return t
}

func (t *Table) Rows(rows []map[string]any) *Table {
	t.RowsArr = append(t.RowsArr, rows...)
	return t
}

func (t *Table) ToMap() map[string]any {
	return map[string]any{
		"template": "table",
		"data":     t,
	}
}

// Text template
type Text struct {
	Text     string `json:"text"`
	Type     string `json:"type,omitempty"`
	Copyable bool   `json:"copyable,omitempty"`
}

func NewText(text string) *Text {
	return &Text{Text: text}
}

func (t *Text) SetType(tp string) *Text {
	t.Type = tp
	return t
}

func (t *Text) SetCopyable(c bool) *Text {
	t.Copyable = c
	return t
}

func (t *Text) ToMap() map[string]any {
	return map[string]any{
		"template": "text",
		"data":     t,
	}
}

// Group template
type Group struct {
	Layout string           `json:"layout,omitempty"` // vertical (default), horizontal
	Items  []map[string]any `json:"items"`
}

func NewGroup() *Group {
	return &Group{Items: []map[string]any{}}
}

func (g *Group) SetHorizontal() *Group {
	g.Layout = "horizontal"
	return g
}

func (g *Group) Add(t Template) *Group {
	g.Items = append(g.Items, t.ToMap())
	return g
}

func (g *Group) ToMap() map[string]any {
	return map[string]any{
		"template": "group",
		"data":     g,
	}
}

// Tabs template
type Tabs struct {
	DefaultTab string           `json:"default_tab,omitempty"`
	Items      []map[string]any `json:"items"`
}

func NewTabs(defaultTab string) *Tabs {
	return &Tabs{DefaultTab: defaultTab, Items: []map[string]any{}}
}

func (t *Tabs) AddTab(key, label string, content Template, icon string) *Tabs {
	t.Items = append(t.Items, map[string]any{
		"key":     key,
		"label":   label,
		"icon":    icon,
		"content": content.ToMap(),
	})
	return t
}

func (t *Tabs) ToMap() map[string]any {
	return map[string]any{
		"template": "tabs",
		"data":     t,
	}
}

// Form template
type Form struct {
	Title      string           `json:"title"`
	Fields     []map[string]any `json:"fields"`
	SubmitText string           `json:"submit_text,omitempty"`
	CancelText string           `json:"cancel_text,omitempty"`
}

func NewForm(title string) *Form {
	return &Form{Title: title, Fields: []map[string]any{}}
}

func (f *Form) AddField(name, label, tp string, required bool, opts ...FormFieldOption) *Form {
	field := map[string]any{
		"name":     name,
		"label":    label,
		"type":     tp,
		"required": required,
	}
	for _, opt := range opts {
		opt(field)
	}
	f.Fields = append(f.Fields, field)
	return f
}

func (f *Form) ToMap() map[string]any {
	return map[string]any{
		"template": "form",
		"data":     f,
	}
}

type FormFieldOption func(map[string]any)

func FormPlaceholder(p string) FormFieldOption {
	return func(m map[string]any) { m["placeholder"] = p }
}

func FormOptions(opts []map[string]any) FormFieldOption {
	return func(m map[string]any) { m["options"] = opts }
}

func FormDefault(d any) FormFieldOption {
	return func(m map[string]any) { m["default"] = d }
}

// Button (Action) template
type Button struct {
	Label    string `json:"label"`
	ActionID string `json:"action_id"`
	Type     string `json:"type,omitempty"` // primary, secondary, danger, link
	Icon     string `json:"icon,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

func NewButton(label, actionID string) *Button {
	return &Button{Label: label, ActionID: actionID, Type: "primary"}
}

func (b *Button) SetType(t string) *Button {
	b.Type = t
	return b
}

func (b *Button) SetIcon(icon string) *Button {
	b.Icon = icon
	return b
}

func (b *Button) SetDisabled(d bool) *Button {
	b.Disabled = d
	return b
}

func (b *Button) ToMap() map[string]any {
	return map[string]any{
		"template": "button",
		"data":     b,
	}
}
