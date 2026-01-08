package main

import (
	"log"
	"os"
	"path/filepath"

	"ukpos-agent/internal/agent"
	"ukpos-agent/internal/config"
	
)

func main() {
	// ------------------------------------------------
	// Resolve executable directory (CRITICAL)
	// ------------------------------------------------
	exePath, err := os.Executable()
	if err != nil {
		panic("cannot resolve executable path: " + err.Error())
	}
	baseDir := filepath.Dir(exePath)

	// ------------------------------------------------
	// Setup logging (MANDATORY FOR AGENTS)
	// ------------------------------------------------
	logDir := filepath.Join(baseDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic("cannot create log directory: " + err.Error())
	}

	logFilePath := filepath.Join(logDir, "agent.log")
	logFile, err := os.OpenFile(
		logFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic("cannot open log file: " + err.Error())
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("🚀 UKPOS Agent starting...")

	// ------------------------------------------------
	// Load config (relative to binary)
	// ------------------------------------------------
	configPath := filepath.Join(baseDir, "config.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("❌ Config error:", err)
	}

	log.Printf(
		"✅ Config loaded | outlet=%s | gateway=%s",
		cfg.OutletID,
		cfg.GatewayURL,
	)

	// ------------------------------------------------
	// Start agent (NON-BLOCKING EXPECTED)
	// ------------------------------------------------
	go func() {
		log.Println("🟢 Agent run loop starting")
		agent.Run(cfg)
		log.Println("🔴 Agent.Run exited unexpectedly")
	}()

	// ------------------------------------------------
	// Block forever (REQUIRED)
	// ------------------------------------------------
	select {}
}


// package main

// import (
// 	"log"
// 	"os"
// 	"path/filepath"

// 	"ukpos-agent/internal/agent"
// 	"ukpos-agent/internal/config"
// 	"ukpos-agent/internal/printer"
// )

// func main() {
// 	// ------------------------------------------------
// 	// Resolve executable directory
// 	// ------------------------------------------------
// 	exePath, err := os.Executable()
// 	if err != nil {
// 		panic("cannot resolve executable path: " + err.Error())
// 	}
// 	baseDir := filepath.Dir(exePath)

// 	// ------------------------------------------------
// 	// Setup logging
// 	// ------------------------------------------------
// 	logDir := filepath.Join(baseDir, "logs")
// 	if err := os.MkdirAll(logDir, 0755); err != nil {
// 		panic("cannot create log directory: " + err.Error())
// 	}

// 	logFilePath := filepath.Join(logDir, "agent.log")
// 	logFile, err := os.OpenFile(
// 		logFilePath,
// 		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
// 		0644,
// 	)
// 	if err != nil {
// 		panic("cannot open log file: " + err.Error())
// 	}

// 	log.SetOutput(logFile)
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

// 	log.Println("🚀 UKPOS Agent starting...")

// 	// ------------------------------------------------
// 	// Load config
// 	// ------------------------------------------------
// 	configPath := filepath.Join(baseDir, "config.json")
// 	cfg, err := config.Load(configPath)
// 	if err != nil {
// 		log.Fatal("❌ Config error:", err)
// 	}

// 	log.Printf(
// 		"✅ Config loaded | outlet=%s | gateway=%s",
// 		cfg.OutletID,
// 		cfg.GatewayURL,
// 	)

// 	// ------------------------------------------------
// 	// ✅ INIT PRINTER QUEUE (CRITICAL)
// 	// ------------------------------------------------
// 	printer.InitQueue(100) // buffer size
// 	log.Println("🖨️ Printer queue ready")

// 	// ------------------------------------------------
// 	// Start agent
// 	// ------------------------------------------------
// 	go func() {
// 		log.Println("🟢 Agent run loop starting")
// 		agent.Run(cfg)
// 		log.Println("🔴 Agent.Run exited unexpectedly")
// 	}()

// 	// ------------------------------------------------
// 	// Block forever
// 	// ------------------------------------------------
// 	select {}
// }

