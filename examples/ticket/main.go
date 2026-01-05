package main

import (
	"fmt"
	"log"
	"time"

	"github.com/tgoai/tgo-plugin-go"
)

// Ticket represents a simple ticket data structure for our mock database.
type Ticket struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"` // Open, In Progress, Closed
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
}

// Mock database
var mockTickets = map[string][]Ticket{
	"visitor_1": {
		{ID: "TK-1001", Title: "无法登录后台", Status: "Open", Priority: "High", CreatedAt: time.Now().Add(-24 * time.Hour)},
		{ID: "TK-1002", Title: "建议增加深色模式", Status: "Closed", Priority: "Low", CreatedAt: time.Now().Add(-72 * time.Hour)},
	},
}

type TicketPlugin struct {
	tgo.BasePlugin
}

func (p *TicketPlugin) ID() string      { return "com.tgo.ticket.go" }
func (p *TicketPlugin) Name() string    { return "ticket-management-go" }
func (p *TicketPlugin) Version() string { return "1.0.0" }

func (p *TicketPlugin) Capabilities() []tgo.Capability {
	return []tgo.Capability{
		// 1. Entry in Chat Toolbar for manual creation
		tgo.ChatToolbar("创建工单", tgo.WithIcon("plus-circle"), tgo.WithTooltip("为当前访客创建工单")),

		// 2. Section in Visitor Panel to show visitor tickets
		tgo.VisitorPanel("相关工单", tgo.WithIcon("ticket")),

		// 3. MCP Tool for AI Agent to create/query tickets
		tgo.MCPTools(
			tgo.Tool("create_ticket", "创建工单").
				Description("根据访客对话内容创建一个新的服务工单。").
				String("title", "工单标题", true).
				String("description", "详细描述", true).
				Enum("priority", "优先级", []string{"low", "medium", "high", "urgent"}, false),
			tgo.Tool("list_tickets", "列出访客工单").
				Description("获取指定访客的所有历史工单列表。"),
		),
	}
}

// --- Visitor Panel Rendering ---

func (p *TicketPlugin) OnVisitorPanelRender(ctx *tgo.RenderContext) tgo.Template {
	visitorID := ctx.VisitorID
	if visitorID == "" {
		return tgo.NewText("请选择一个访客以查看工单信息。")
	}

	tickets := mockTickets[visitorID]

	group := tgo.NewGroup()

	// Title and Action Button
	header := tgo.NewGroup().SetHorizontal().
		Add(tgo.NewText("访客工单记录").SetSize("lg").SetBold(true)).
		Add(tgo.NewButton("手动创建", "open_create_form").SetIcon("plus").SetType("primary"))
	group.Add(header)

	if len(tickets) == 0 {
		group.Add(tgo.NewText("该访客暂无工单记录。").SetColor("#999"))
	} else {
		table := tgo.NewTable("").Columns("ID", "标题", "状态", "优先级")
		for _, t := range tickets {
			statusColor := "blue"
			if t.Status == "Closed" {
				statusColor = "green"
			} else if t.Priority == "High" {
				statusColor = "orange"
			}

			table.Row(map[string]any{
				"ID":  t.ID,
				"标题":  t.Title,
				"状态":  map[string]any{"text": t.Status, "color": statusColor},
				"优先级": t.Priority,
			})
		}
		group.Add(table)
	}

	return group
}

// --- Chat Toolbar Rendering ---

func (p *TicketPlugin) OnChatToolbarRender(ctx *tgo.RenderContext) tgo.Template {
	// Directly return the create form when the toolbar entry is clicked (if supported by host)
	// or return a button that triggers the form.
	return tgo.NewForm("新建工单").
		Add(tgo.NewFormField("title", "工单标题", "text").SetRequired(true).SetPlaceholder("简述问题...")).
		Add(tgo.NewFormField("priority", "优先级", "select").
			AddOption("低", "low").
			AddOption("中", "medium").
			AddOption("高", "high").
			SetDefault("medium")).
		Add(tgo.NewFormField("description", "详细描述", "textarea").SetPlaceholder("请输入详细描述..."))
}

// --- Event Handling (Buttons/Forms) ---

func (p *TicketPlugin) OnVisitorPanelEvent(ctx *tgo.EventContext) *tgo.Action {
	return p.handleCommonEvents(ctx)
}

func (p *TicketPlugin) OnChatToolbarEvent(ctx *tgo.EventContext) *tgo.Action {
	return p.handleCommonEvents(ctx)
}

func (p *TicketPlugin) handleCommonEvents(ctx *tgo.EventContext) *tgo.Action {
	switch ctx.ActionID {
	case "open_create_form", "创建工单":
		// Show a form to create a ticket
		form := tgo.NewForm("新建工单").
			Add(tgo.NewFormField("title", "工单标题", "text").SetRequired(true).SetPlaceholder("简述问题...")).
			Add(tgo.NewFormField("priority", "优先级", "select").
				AddOption("低", "low").
				AddOption("中", "medium").
				AddOption("高", "high").
				SetDefault("medium")).
			Add(tgo.NewFormField("description", "详细描述", "textarea").SetPlaceholder("请输入详细描述..."))

		return tgo.ShowModal("新建工单", form)

	case "submit_ticket":
		// Handle form submission
		title, _ := ctx.FormData["title"].(string)
		priority, _ := ctx.FormData["priority"].(string)

		// In a real app, you'd save to DB here
		newID := fmt.Sprintf("TK-%d", 1000+len(mockTickets[ctx.VisitorID])+1)
		mockTickets[ctx.VisitorID] = append(mockTickets[ctx.VisitorID], Ticket{
			ID:        newID,
			Title:     title,
			Status:    "Open",
			Priority:  priority,
			CreatedAt: time.Now(),
		})

		return tgo.ShowToast(fmt.Sprintf("工单 %s 创建成功", newID), "success")
	}

	return tgo.Noop()
}

// --- MCP Tool Execution ---

func (p *TicketPlugin) OnToolExecute(ctx *tgo.ToolContext, toolName string, args map[string]any) (*tgo.ToolResult, error) {
	log.Printf("MCP Tool Execute: %s (Visitor: %s)", toolName, ctx.VisitorID)
	fmt.Println("MCP Tool Execute: %s (Visitor: %s)", toolName, ctx.VisitorID)

	switch toolName {
	case "create_ticket":
		title, _ := args["title"].(string)
		desc, _ := args["description"].(string)
		priority, _ := args["priority"].(string)
		if priority == "" {
			priority = "medium"
		}

		if ctx.VisitorID == "" {
			return &tgo.ToolResult{Success: false, Content: "无法识别访客，请在会话中调用。"}, nil
		}

		// Create ticket in mock DB
		newID := fmt.Sprintf("TK-%d", 1000+len(mockTickets[ctx.VisitorID])+1)
		mockTickets[ctx.VisitorID] = append(mockTickets[ctx.VisitorID], Ticket{
			ID:        newID,
			Title:     title,
			Status:    "Open",
			Priority:  priority,
			CreatedAt: time.Now(),
		})

		return &tgo.ToolResult{
			Success: true,
			Content: fmt.Sprintf("成功为访客创建工单！工单 ID: %s, 标题: %s, 优先级: %s。描述: %s", newID, title, priority, desc),
			Data: map[string]any{
				"ticket_id": newID,
				"status":    "Open",
			},
		}, nil

	case "list_tickets":
		tickets := mockTickets[ctx.VisitorID]
		if len(tickets) == 0 {
			return &tgo.ToolResult{Success: true, Content: "该访客目前没有工单。"}, nil
		}

		content := "该访客的工单列表如下：\n"
		for _, t := range tickets {
			content += fmt.Sprintf("- [%s] %s (状态: %s, 优先级: %s)\n", t.ID, t.Title, t.Status, t.Priority)
		}
		return &tgo.ToolResult{Success: true, Content: content, Data: map[string]any{"tickets": tickets}}, nil
	}

	return &tgo.ToolResult{Success: false, Content: "未知工具"}, nil
}

func main() {
	// Start the plugin, connecting to the TGO API via TCP
	// Use 8005 for local debugging with Docker-based TGO API
	if err := tgo.Run(&TicketPlugin{}, tgo.WithTCPAddr("localhost:8005"), tgo.WithDevToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwcm9qZWN0X2lkIjoiOTI2MGEwODMtYWI3YS00MGU5LWI4MDQtNzUzMjhjMjFlZDNhIiwidXNlcl9pZCI6IjQ4NTQ2YmU2LTc2NTQtNDU4NS05MmEwLTMwMjNlY2MyYjZlMyIsInR5cGUiOiJwbHVnaW5fZGV2IiwiZXhwIjoxNzY3NjkwNzM0fQ.2ULd2ztlGO517zRkz2c-EWrPNXgH332Gb4dwbQYqJUY")); err != nil {
		log.Fatalf("Ticket Plugin failed: %v", err)
	}
}
