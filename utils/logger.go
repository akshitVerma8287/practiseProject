package utils

import "sync"

// Global Variables
var (
	ApiLogs  []string
	LogMutex sync.Mutex
)

func AddLog(log string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	ApiLogs = append(ApiLogs, log)
}

func GetLogs() []string {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	return append([]string{}, ApiLogs...)
}

func ClearLogs() {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	ApiLogs = []string{}
}