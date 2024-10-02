package models

type PPSQueueSize struct {
	PPSID     string `json:"pps_id"`
	QueueSize string `json:"queue_size"`
}
