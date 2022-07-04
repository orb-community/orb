package migrate

type Service interface {
	Up() error
	Down() error
	Drop() error
	SetSchemaVersion(int64) error
	CurrentSchemaVersion() (int64, error)
	LatestSchemaVersion() int64
}
