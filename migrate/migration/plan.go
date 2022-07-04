package migration

type Plan interface {
	Up() error
	Down() error
}
