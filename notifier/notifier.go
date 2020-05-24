package notifier

type Notifier interface {
	Send(string) error
}