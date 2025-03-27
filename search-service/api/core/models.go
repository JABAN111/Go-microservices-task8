package core

type UpdateStatus string

const (
	StatusUpdateUnknown UpdateStatus = "unknown"
	StatusUpdateIdle    UpdateStatus = "idle"
	StatusUpdateRunning UpdateStatus = "running"
)

type StatsReply struct {
	WordsTotal    int64
	WordsUnique   int64
	ComicsTotal   int64
	ComicsFetched int64
}
type Comics struct {
	ID     string
	ImgUrl string
}
