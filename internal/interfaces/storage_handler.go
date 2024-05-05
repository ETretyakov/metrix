package interfaces

type StorageHandler interface {
	Key(...string) string
	Get(string) (*float64, error)
	Set(string, float64) (*float64, error)
	Keys(string) ([]string, error)
}
