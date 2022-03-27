package datasource

type sourceloader interface {
	load() []map[string]string
}
