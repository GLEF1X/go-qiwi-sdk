package polling

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi"
	"github.com/GLEF1X/go-qiwi-sdk/types"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi/filters"
)

const maxHistoryTransactionLimit = 50

const (
	defaultConcurrentMaximum = 1
	defaultPollingTimeout    = 5 * time.Second
)

type Handler struct {
	bunchOfFilters []Filter

	// function, that handle some event
	fn interface{}
}

func (handler Handler) ProcessEvent(event interface{}, alreadyProcessed chan struct{}) {
	select {
	case <-alreadyProcessed:
		return
	default:
		for _, filter := range handler.bunchOfFilters {
			if !filter.Check(event) {
				return
			}
		}
		reflect.ValueOf(handler.fn).Call([]reflect.Value{reflect.ValueOf(event)})
		alreadyProcessed <- struct{}{}
	}
}

type Dispatcher struct {
	config *Config

	// Contains all handlers, except of error handlers
	handlerStack []Handler
	updates      chan interface{}
	stop         chan struct{}

	errorHandlers []Handler

	// special ID of last transaction, which was being processed. It helps us to determine which transactions
	// are not processed yet.
	offset int
}

type Config struct {
	GetUpdatesFromDate *time.Time
	PollingTimeout     time.Duration

	// Max updates, that would be processed by handlers concurrently
	MaxConcurrent int
}

func NewDispatcher(config *Config) *Dispatcher {
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = defaultConcurrentMaximum
	}
	if config.PollingTimeout == 0 {
		config.PollingTimeout = defaultPollingTimeout
	}

	return &Dispatcher{
		config:  config,
		stop:    make(chan struct{}),
		updates: make(chan interface{}, config.MaxConcurrent),
	}
}

func (dp *Dispatcher) HandleTransaction(handler func(*types.Transaction), additionalFilters ...Filter) {
	dp.addHandler(
		handler,
		func(event interface{}) bool {
			_, ok := event.(*types.Transaction)
			return ok
		},
		additionalFilters...,
	)
}

func (dp *Dispatcher) HandleError(handler func(error), additionalFilters ...Filter) {
	dp.addHandler(
		handler,
		func(event interface{}) bool {
			_, ok := event.(error)
			return ok
		},
		additionalFilters...,
	)
}

func (dp *Dispatcher) addHandler(handler interface{}, defaultFilter func(interface{}) bool, filters ...Filter) {
	eventFilters := []Filter{makeDefaultFilterFromFunc(defaultFilter)}
	eventFilters = append(eventFilters, filters...)

	dp.handlerStack = append(dp.handlerStack, Handler{
		fn:             handler,
		bunchOfFilters: eventFilters,
	})
}

func (dp *Dispatcher) StartPolling(apiClient *qiwi.APIClient) {
	defer apiClient.Close()
	if dp.config.GetUpdatesFromDate == nil {
		currentTime := time.Now()
		dp.config.GetUpdatesFromDate = &currentTime
	}

	go dp.fetchNewEvents(apiClient)

	for {
		select {
		case update := <-dp.updates:
			log.Println("Get update from channel, starting propagate it")
			dp.PropagateEventToHandlers(update)
		case <-dp.stop:
			close(dp.stop)
			return
		}
	}
}

func (dp *Dispatcher) StopPolling() {
	dp.stop <- struct{}{}
}

func (dp *Dispatcher) fetchNewEvents(apiClient *qiwi.APIClient) {
	for {
		select {
		case <-dp.stop:
			return
		default:
		}

		currentTime := time.Now()
		history, err := apiClient.RetrieveHistory(context.Background(), &filters.HistoryFilter{
			StartDate: dp.config.GetUpdatesFromDate,
			EndDate:   &currentTime,
		})

		if err != nil {
			log.Printf("Catched exception %s", err.Error())
		}

		for _, update := range history.Transactions {
			if !(update.ID > dp.offset) {
				continue
			}
			dp.updates <- &update
		}

		sortHistoryByIDAscending(history)
		if len(history.Transactions) != 0 {
			dp.offset = history.Transactions[len(history.Transactions)-1].ID
		} else if len(history.Transactions) == maxHistoryTransactionLimit {
			dp.config.GetUpdatesFromDate = &history.Transactions[len(history.Transactions)-1].Date
		}

		time.Sleep(dp.config.PollingTimeout)
	}
}

func (dp *Dispatcher) PropagateEventToHandlers(event interface{}) {
	processChan := make(chan struct{}, 1)
	for _, handler := range dp.handlerStack {
		log.Printf("Get event %T. Start iterating throw handlerStack", event)
		go func(handler Handler, event interface{}) {
			defer dp.deferDebug()
			handler.ProcessEvent(event, processChan)
		}(handler, event)
	}
}

func (dp *Dispatcher) deferDebug() {
	if r := recover(); r != nil {
		if err, ok := r.(error); ok {
			dp.debug(err)
		} else if str, ok := r.(string); ok {
			dp.debug(fmt.Errorf("%s", str))
		}
	}
}
func (dp *Dispatcher) debug(err error) {
	processChan := make(chan struct{}, 1)
	for _, handler := range dp.errorHandlers {
		handler.ProcessEvent(err, processChan)
	}
	log.Println(err)
}
