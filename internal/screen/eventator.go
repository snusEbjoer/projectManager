package screen

import (
	"log"

	"github.com/eiannone/keyboard"
)

type Event struct {
	Key  keyboard.Key
	Char string
}

type Eventator struct {
	eventsChan chan Event
}

func (e *Eventator) GetChan() <-chan Event {
	return e.eventsChan
}

func NewEventator() *Eventator {
	if err := keyboard.Open(); err != nil {
		log.Fatal(err.Error())
	}
	e := &Eventator{
		eventsChan: make(chan Event),
	}
	e.run()
	return e
}

func (e *Eventator) run() {
	go func() {
		for {
			r, key, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err.Error())
			}
			e.eventsChan <- Event{
				Key:  key,
				Char: string(r),
			}
		}
	}()
}
