package index

// Index is the Unique Identification of calling `vm.Run`
type Index struct {
	Worker int64 `mapstructure:"worker"`
	VM     int64 `mapstructure:"vm"`
	Engine int64 `mapstructure:"engine"`
	Tx     int64 `mapstructure:"tx"`
}
