package worker

import "time"

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusDone       JobStatus = "done"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           string
	OriginalKey  string // path to raw uploaded file in MinIO
	OriginalName string
	SubmittedAt  time.Time
	Status       JobStatus
	Error        string
	Result       *JobResult
}

type JobResult struct {
	Tracks   []TranscodedTrack
	Waveform []float64
	Metadata AudioMetadata
}

type TranscodedTrack struct {
	Format  string // "mp3", "aac"
	Bitrate int    // kbps
	Key     string // MinIO object key
}

type AudioMetadata struct {
	Title    string
	Artist   string
	Album    string
	Duration float64 // seconds
	Format   string
	Bitrate  int
}
