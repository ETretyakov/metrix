package interfaces

type StorageHandler interface {
	Key(...string) string
	Get(string) (*uint64, error)
	Set(string, uint64) (*uint64, error)
}
