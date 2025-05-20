package main

import (
	"fmt"

	"github.com/reactivex/rxgo/v2"
)

type Event struct {
	Type string
	Data interface{}
}

type EventBus struct {
	consumer  map[string][]chan rxgo.Item
	producer  rxgo.Observable
	inputChan chan rxgo.Item
}

func NewEventBus() *EventBus {
	inputChan := make(chan rxgo.Item, 100)
	bus := &EventBus{
		consumer:  make(map[string][]chan rxgo.Item),
		inputChan: inputChan,
	}

	bus.producer = rxgo.FromChannel(inputChan)

	go func() {
		for item := range bus.producer.Observe() {
			fmt.Println(item.V)
		}
	}()

	return bus
}

func (eb *EventBus) Publish(event Event) {
	select {
	case eb.inputChan <- rxgo.Of(event):
	}
}

func (eb *EventBus) Subscribe(eventType string, handler func(event Event)) {

}

func main() {
	bus := NewEventBus()

	bus.Publish(Event{
		Type: "test",
	})

	// observable := rxgo.Defer([]rxgo.Producer{func(_ context.Context, ch chan<- rxgo.Item) {
	// 	for i := 0; i < 3; i++ {
	// 		ch <- rxgo.Of(i)
	// 	}
	// }})
	//
	// // First Observer
	// for item := range observable.Observe() {
	// 	fmt.Println(item.V)
	// }
	//
	// // Second Observer
	// for item := range observable.Observe() {
	// 	fmt.Println(item.V)
	// }
}
