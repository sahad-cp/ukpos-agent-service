package printer

import (
	"log"
	"time"
)

type PrintJob struct {
	IP        string
	Port      int
	Timeout   time.Duration
	Data      []byte
	JobType   string
	OnSuccess func()
	OnError   func(error)
}
var (
	queue chan PrintJob
)

func InitQueue(buffer int) {
	queue = make(chan PrintJob, buffer)
	go worker()
	log.Println("🖨️ Printer queue initialized")
}

func Enqueue(job PrintJob) {
	queue <- job
}

func worker() {
	for job := range queue {
		log.Printf("🖨️ Printing job (%s) → %s:%d\n", job.JobType, job.IP, job.Port)

		// err := SendTCP(job.IP, job.Port, job.Timeout, job.Data)
		// if err != nil {
		// 	log.Println("❌ Printer error:", err)
		// }

		// ✅ CRITICAL: cooldown
		time.Sleep(2 * time.Second)
	}
}