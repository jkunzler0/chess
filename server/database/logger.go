package database

type EventType byte

const (
	_                     = iota // iota == 0; ignore this value
	EventPut    EventType = iota // iota == 1
	EventDelete                  // iota == 2
	EventIncr                    // iota == 3
)

type Event struct {
	Sequence  uint64
	EventType EventType
	User1     string
	User2     string
	Value     Score
}

type TransactionLogger interface {
	WritePut(user string, value Score)
	WriteDelete(user string)
	WriteIncr(winner string, losser string)

	Err() <-chan error

	LastSequence() uint64

	Run()
	Wait()
	Close() error

	ReadEvents() (<-chan Event, <-chan error)

	StartTransactionLog() error
}
