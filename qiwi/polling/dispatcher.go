package polling

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"time"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi"
	"github.com/GLEF1X/go-qiwi-sdk/types"

	"github.com/GLEF1X/go-qiwi-sdk/qiwi/filters"
)

const (
	maxHistoryTransactionLimit = 50

	defaultConcurrentMaximum = 10
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
	GetUpdatesFromDate time.Time

	// Max updates, that would be processed by handlers concurrently
	MaxConcurrent int
}

func NewDispatcher(config *Config) *Dispatcher {
	if config.MaxConcurrent == 0 {
		config.MaxConcurrent = defaultConcurrentMaximum
	}

	return &Dispatcher{
		config:  config,
		stop:    make(chan struct{}),
		updates: make(chan interface{}, config.MaxConcurrent),
	}
}

func (dp *Dispatcher) HandleTransaction(handler func(*types.Transaction), additionalFilters ...Filter) {
	dp.addHandler(handler, transactionFilter{}, additionalFilters...)
}

func (dp *Dispatcher) HandleError(handler func(error), additionalFilters ...Filter) {
	dp.addHandler(handler, errorFilter{}, additionalFilters...)
}

func (dp *Dispatcher) addHandler(handler interface{}, defaultFilter Filter, filters ...Filter) {
	eventFilters := []Filter{defaultFilter}
	eventFilters = append(eventFilters, filters...)

	dp.handlerStack = append(dp.handlerStack, Handler{
		fn:             handler,
		bunchOfFilters: eventFilters,
	})
}

func (dp *Dispatcher) LaunchPolling(apiClient *qiwi.APIClient) {
	defer apiClient.Close()
	dp.config.GetUpdatesFromDate = time.Now()

	go dp.fetchNewEvents(apiClient)

	for {
		select {
		case update := <-dp.updates:
			log.Println("Get update from channel, starting propagate it")
			dp.propagateEvent(update)
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
		log.Println("Start fetching new events...")
		currentTime := time.Now()
		log.Printf("Poll transactions from %s to %s", dp.config.GetUpdatesFromDate.Format(time.RFC3339), currentTime.Format(time.RFC3339))
		history, err := apiClient.RetrieveHistory(context.Background(), &filters.HistoryFilter{
			StartDate: &dp.config.GetUpdatesFromDate,
			EndDate:   &currentTime,
		})
		if err != nil {
			log.Printf("Catched exception %s", err.Error())
		}
		for _, update := range history.Transactions {
			if !(update.ID > dp.offset) {
				log.Printf("Skip update %d due to offset is bigger: %d > %d", update.ID, update.ID, dp.offset)
				continue
			}
			log.Printf("Iterate throw the history, current update id is %d", update.ID)
			dp.updates <- &update
		}
		sort.Slice(history.Transactions, func(i, j int) bool {
			return history.Transactions[i].ID < history.Transactions[j].ID
		})
		log.Println("Current sorted history: ", history.Transactions)
		if len(history.Transactions) != 0 {
			log.Printf("History is not empty, set new offset to %d", history.Transactions[len(history.Transactions)-1].ID)
			dp.offset = history.Transactions[len(history.Transactions)-1].ID
		} else if len(history.Transactions) == maxHistoryTransactionLimit {
			dp.config.GetUpdatesFromDate = history.Transactions[len(history.Transactions)-1].Date
		}
		log.Println("Sleeping...")
		time.Sleep(5 * time.Second)
	}
}

func (dp *Dispatcher) propagateEvent(event interface{}) {
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
