package command

type Interface interface {
	Name() string
	Description() string
	Process(string) string
}
