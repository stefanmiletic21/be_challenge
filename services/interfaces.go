package services

type WorkerDataProvider interface {
MakeAggregatedTable() error
UpdateHour(hour int64) (err error)
GetFurthestUpdatedHour() (hour int64, err error)
}