package cronjobx

type Job interface {
	Name() string
	Run() error
}
