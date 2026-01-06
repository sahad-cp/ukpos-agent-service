package agent

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ukpos-agent/internal/config"
	"ukpos-agent/internal/escpos"
	"ukpos-agent/internal/printer"
	ws "ukpos-agent/internal/websocket"
)

func Run(cfg *config.Config) {
	for {
		log.Println("🔌 Connecting to gateway...")

		client, err := ws.Connect(cfg.GatewayURL)
		if err != nil {
			log.Println("❌ Gateway offline, retrying...")
			time.Sleep(time.Duration(cfg.ReconnectInterval) * time.Millisecond)
			continue
		}

		log.Println("✅ Gateway connected")

		client.Send(map[string]any{
			"type":     "REGISTER",
			"outletId": cfg.OutletID,
			"version":  cfg.AgentVersion,
		})

		ctx, cancel := context.WithCancel(context.Background())
		go heartbeat(ctx, client, cfg)

		client.Listen(func(msg map[string]any) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("🔥 Panic recovered:", r)
				}
			}()

			msgType, _ := msg["type"].(string)

			switch msgType {

			case "PRINT_JOB":
				job := msg["job"].(map[string]any)

				jobId := job["jobId"].(string)
				ip := job["printerIP"].(string)
				port := int(job["printerPort"].(float64))
				jobType := job["jobType"].(string)

				client.Send(map[string]any{
					"type":   "JOB_STATUS",
					"jobId": jobId,
					"status": "PRINTING",
				})

				var data []byte

				switch jobType {

				case "TEST":
					data = escpos.TestPrint(cfg.OutletID, ip)

				case "KOT":
					var kot escpos.KOTData
					parse(job["data"], &kot)
					data = escpos.BuildKOT(kot)

				case "RECEIPT":
					var rec escpos.ReceiptData
					parse(job["data"], &rec)
					data = escpos.BuildReceipt(rec)

				default:
					log.Println("❌ Unknown jobType:", jobType)
					return
				}

				if err := printer.Send(ip, port, cfg.PrinterTimeout, data); err != nil {
					client.Send(map[string]any{
						"type":   "JOB_STATUS",
						"jobId": jobId,
						"status": "FAILED",
						"error":  err.Error(),
					})
					log.Println("❌ Print failed:", err)
				} else {
					client.Send(map[string]any{
						"type":   "JOB_STATUS",
						"jobId": jobId,
						"status": "COMPLETED",
					})
					log.Println("🖨️ Printed:", jobType)
				}
			}
		})

		log.Println("⚠️ Gateway disconnected")
		cancel()
		client.Close()
	}
}

func parse(src any, dst any) {
	b, _ := json.Marshal(src)
	_ = json.Unmarshal(b, dst)
}
func heartbeat(ctx context.Context, c *ws.Client, cfg *config.Config) {
	ticker := time.NewTicker(
		time.Duration(cfg.HeartbeatInterval) * time.Millisecond,
	)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("💔 Heartbeat stopped")
			return
		case <-ticker.C:
			c.Send(map[string]any{"type": "HEARTBEAT"})
		}
	}
}