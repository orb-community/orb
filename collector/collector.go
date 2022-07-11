package collector

type collector interface {
	Start() error
	Stop() error
}
