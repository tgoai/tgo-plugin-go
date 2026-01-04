package main

import (
	"fmt"
	"log"

	"github.com/tgoai/tgo-plugin-go"
)

type CRMPlugin struct {
	tgo.BasePlugin
}

func (p *CRMPlugin) ID() string      { return "com.tgo.crm.go" }
func (p *CRMPlugin) Name() string    { return "crm-info-go" }
func (p *CRMPlugin) Version() string { return "1.0.0" }

func (p *CRMPlugin) Capabilities() []tgo.Capability {
	return []tgo.Capability{
		tgo.VisitorPanel("客户档案 (Go)", tgo.WithIcon("user")),
		tgo.ChatToolbar("CRM 操作 (Go)", tgo.WithIcon("database")),
	}
}

func (p *CRMPlugin) OnVisitorPanelRender(ctx *tgo.RenderContext) tgo.Template {
	title := "客户信息"
	if ctx.Language == "en" {
		title = "Customer Info"
	}

	if ctx.Visitor == nil {
		if ctx.Language == "en" {
			return tgo.NewKeyValue("Unknown Visitor").Add("Tip", "Unable to get visitor info")
		}
		return tgo.NewKeyValue("未知访客").Add("提示", "无法获取访客信息")
	}

	// Build UI
	group := tgo.NewGroup()

	// 1. Basic Info
	idLabel, nameLabel, levelLabel := "ID", "姓名", "等级"
	levelValue := "铂金会员"
	if ctx.Language == "en" {
		nameLabel, levelLabel = "Name", "Level"
		levelValue = "Platinum"
	}

	info := tgo.NewKeyValue(title).
		Add(idLabel, ctx.VisitorID, tgo.KeyValueCopyable(true)).
		Add(nameLabel, ctx.Visitor.Name).
		Add(levelLabel, levelValue, tgo.KeyValueIcon("crown"), tgo.KeyValueColor("#FFD700"))
	group.Add(info)

	// 2. Orders Table
	orderTitle := "最近订单"
	orderCol, amountCol, statusCol := "订单号", "金额", "状态"
	statusText := "配送中"
	statusDone := "已完成"
	if ctx.Language == "en" {
		orderTitle = "Recent Orders"
		orderCol, amountCol, statusCol = "Order ID", "Amount", "Status"
		statusText = "Shipping"
		statusDone = "Completed"
	}

	table := tgo.NewTable(orderTitle).
		Columns(orderCol, amountCol, statusCol).
		Row(map[string]any{
			orderCol:  "GO-001",
			amountCol: "¥1,299",
			statusCol: map[string]any{"text": statusText, "color": "blue"},
		}).
		Row(map[string]any{
			orderCol:  "GO-002",
			amountCol: "¥88",
			statusCol: map[string]any{"text": statusDone, "color": "green"},
		})
	group.Add(table)

	// 3. Actions
	viewBtn, sendBtn := "查看 CRM 详情", "发送优惠券"
	if ctx.Language == "en" {
		viewBtn, sendBtn = "View CRM Details", "Send Coupon"
	}

	actions := tgo.NewGroup().SetHorizontal().
		Add(tgo.NewButton(viewBtn, "view_crm").SetIcon("external-link")).
		Add(tgo.NewButton(sendBtn, "send_coupon").SetType("secondary").SetIcon("ticket"))
	group.Add(actions)

	return group
}

func (p *CRMPlugin) OnVisitorPanelEvent(ctx *tgo.EventContext) *tgo.Action {
	if ctx.ActionID == "view_crm" {
		return tgo.OpenURL(fmt.Sprintf("https://crm.example.com/visitor/%s", ctx.VisitorID), "_blank")
	}
	if ctx.ActionID == "send_coupon" {
		return tgo.ShowToast("优惠券已发送给访客", "success")
	}
	return tgo.ShowToast("收到事件: "+ctx.EventType, "info")
}

func (p *CRMPlugin) OnChatToolbarEvent(ctx *tgo.EventContext) *tgo.Action {
	fmt.Println("OnChatToolbarEvent", ctx.EventType, ctx.ActionID, ctx.VisitorID, ctx.SessionID)
	return tgo.ShowToast("收到事件: "+ctx.EventType, "info")
}

func (p *CRMPlugin) OnChatToolbarRender(ctx *tgo.RenderContext) tgo.Template {
	return tgo.NewForm("创建工单").
		AddField("title", "标题", "text", true, tgo.FormPlaceholder("请输入工单标题")).
		AddField("priority", "优先级", "select", true, tgo.FormOptions([]map[string]any{
			{"label": "高", "value": "high"},
			{"label": "低", "value": "low"},
		}))
}

func main() {
	plugin := &CRMPlugin{}
	// On macOS, Unix socket bind mounts from Docker are not accessible from host.
	// Use TCP port 8005 for local debugging.
	if err := tgo.Run(plugin, tgo.WithTCPAddr("localhost:8005")); err != nil {
		log.Fatalf("Plugin exited: %v", err)
	}
}
