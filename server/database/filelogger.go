package database

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

// #######################################################################
// (Section 1) Initialization ############################################
// #######################################################################

type FileTransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error
	lastSequence uint64   // The last used event sequence number
	file         *os.File // The location of the transaction log
	wg           *sync.WaitGroup
}

func NewFileTransactionLogger(filename string) (TransactionLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		return nil, fmt.Errorf("cannot open transaction log file: %w", err)
	}

	return &FileTransactionLogger{file: file, wg: &sync.WaitGroup{}}, nil
}

// #######################################################################
// (Section 2) Startup ###################################################
// #######################################################################

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// Start retrieving events from the events channel and writing them
	// to the transaction log
	go func() {
		for e := range events {
			l.lastSequence++

			_, err := fmt.Fprintf(
				l.file,
				"%d\t%d\t%s\t%s\t%d\t%d\n",
				l.lastSequence, e.EventType,
				e.User1, e.User2, e.Value.Win, e.Value.Loss)

			if err != nil {
				errors <- fmt.Errorf("cannot write to log file: %w", err)
			}

			l.wg.Done()
		}
	}()
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			fmt.Sscanf(
				line, "%d\t%d\t%s\t%s\t%d\t%d",
				&e.Sequence, &e.EventType,
				&e.User1, &e.User2, &e.Value.Win, &e.Value.Loss)

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func (l *FileTransactionLogger) StartTransactionLog() error {
	var err error

	events, errors := l.ReadEvents()
	count, ok, e := 0, true, Event{}

	for ok && err == nil {
		select {
		case err, ok = <-errors:

		case e, ok = <-events:
			switch e.EventType {
			case EventDelete: // Delete event
				err = Delete(e.User1)
				count++
			case EventPut: // Put event
				err = Put(e.User1, e.Value)
				count++
			case EventIncr: // Incr event
				err = IncrWinLoss(e.User1, e.User2)
				count++
			}
		}
	}

	log.Printf("%d events replayed\n", count)

	l.Run()

	go func() {
		for err := range l.Err() {
			log.Print(err)
		}
	}()

	return err

}

// #######################################################################
// (Section 3) Event Writes ##############################################
// #######################################################################

func (l *FileTransactionLogger) WritePut(user string, value Score) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventPut, User1: user, User2: "na", Value: Score{value.Win, value.Loss}}
}

func (l *FileTransactionLogger) WriteDelete(user string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventDelete, User1: user}
}

func (l *FileTransactionLogger) WriteIncr(winner string, losser string) {
	l.wg.Add(1)
	l.events <- Event{EventType: EventIncr, User1: winner, User2: losser}
}

// #######################################################################
// (Section 4) Helper/Misc Functions  ####################################
// #######################################################################

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) LastSequence() uint64 {
	return 0
}

func (l *FileTransactionLogger) Wait() {
	l.wg.Wait()
}

func (l *FileTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.file.Close()
}
