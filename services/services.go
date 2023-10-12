package services

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atul-007/GreedyGameAssignment/models"
)

type DbServicesInterface interface {
	Set(key, value string, expiryTime *time.Time, condition string) error
	Get(key string) (string, error)
	QPush(key string, values []string) error
	QPop(key string) (string, error)
	BQPop(key string, timeout string) (string, error)
	ProcessRequest(key string, timeout string) string
	ParseExpiryTime(cmd []string) *time.Time
	ParseCondition(cmd []string) string
}

type DbServices struct {
	store map[string]models.KeyValue
	mu    sync.RWMutex
}

func (ds *DbServices) Set(key, value string, expiryTime *time.Time, condition string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	if condition == "NX" {
		if _, exists := ds.store[key]; exists {
			return nil
		}
	} else if condition == "XX" {
		if _, exists := ds.store[key]; !exists {
			return nil
		}
	}
	ds.store[key] = models.KeyValue{
		Value:      value,
		ExpiryTime: expiryTime,
	}

	go func() {
		for {
			time.Sleep(1 * time.Second)

			ds.mu.Lock()
			for key, kv := range ds.store {
				if kv.ExpiryTime != nil && time.Now().After(*kv.ExpiryTime) {
					delete(ds.store, key)
				}
			}
			ds.mu.Unlock()
		}
	}()

	return nil
}

func (ds *DbServices) Get(key string) (string, error) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()

	if kv, exists := ds.store[key]; exists {
		return kv.Value, nil
	}

	return "", nil
}

func (ds *DbServices) QPush(key string, values []string) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	Check.Store(key)

	queue, exists := ds.store[key]
	if !exists {
		queue = models.KeyValue{
			Value:      "",
			ExpiryTime: nil,
		}
	}

	queueValue := models.Queue{Items: strings.Fields(queue.Value)}
	queueValue.Items = append(queueValue.Items, values...)

	queue.Value = strings.Join(queueValue.Items, " ")
	ds.store[key] = queue

	return nil
}

func (ds *DbServices) QPop(key string) (string, error) {

	queue, exists := ds.store[key]
	if !exists || queue.Value == "" {
		return "", nil
	}

	queueValue := models.Queue{Items: strings.Fields(queue.Value)}
	if len(queueValue.Items) == 0 {
		return "", nil
	}

	poppedValue := queueValue.Items[len(queueValue.Items)-1]
	queueValue.Items = queueValue.Items[:len(queueValue.Items)-1]

	queue.Value = strings.Join(queueValue.Items, " ")
	ds.store[key] = queue

	return poppedValue, nil

}

var Check atomic.Value

func (ds *DbServices) BQPop(key string, timeout string) (string, error) {

	if Check.Load() != nil && Check.Load() != "b" {
		return "BUSY", nil
	}

	Check.Store("a")

	response := ds.ProcessRequest(key, timeout)
	fmt.Println(response)
	if response == "ok" {
		queue, exists := ds.store[key]
		if !exists || queue.Value == "" {
			return "", nil
		}

		queueValue := models.Queue{Items: strings.Fields(queue.Value)}
		if len(queueValue.Items) == 0 {
			return "", nil
		}

		poppedValue := queueValue.Items[len(queueValue.Items)-1]
		queueValue.Items = queueValue.Items[:len(queueValue.Items)-1]

		queue.Value = strings.Join(queueValue.Items, " ")
		ds.store[key] = queue

		return poppedValue, nil
	}

	return "", nil

}
func (ds *DbServices) ProcessRequest(key string, timeout string) string {

	seconds, _ := strconv.Atoi(timeout)
	duration := time.Duration(seconds) * time.Second
	curr_time := time.Now()
	t := time.Now().Add(duration)
	then := t.Unix()
	for then >= curr_time.Unix() {
		if Check.Load() == key {
			break
		} else {
			curr_time = time.Now()
		}
	}
	return "ok"

}

func (ds *DbServices) ParseExpiryTime(cmd []string) *time.Time {
	for i := 3; i < len(cmd); i++ {
		if cmd[i] == "EX" && i+1 < len(cmd) {
			seconds, err := strconv.Atoi(cmd[i+1])
			if err == nil {
				expiryTime := time.Now().Add(time.Duration(seconds) * time.Second)
				return &expiryTime
			}
		}
	}
	return nil
}

func (ds *DbServices) ParseCondition(cmd []string) string {
	for _, param := range cmd[3:] {
		if param == "NX" || param == "XX" {
			return param
		}
	}
	return ""
}
