//go:build trace
// +build trace

package functrace

import (
	"bytes"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Ltime|log.Ldate)
}

var (
	mu sync.Mutex
	m  = make(map[uint64]int)
)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func printTraceEntry(id uint64, name, typ string, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}

	logger.Printf("g[%02d]:%s%s%s\n", id, indents, typ, name)
}

func printTraceExit(id uint64, name, typ string, indent int, cost time.Duration) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}
	logger.Printf("g[%02d]:%s%s%s\t\tcost: %s\n", id, indents, typ, name, cost)
}

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	id := getGID()
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	started := time.Now()

	mu.Lock()
	v := m[id]
	m[id] = v + 1
	mu.Unlock()
	printTraceEntry(id, name, "->", v+1)
	return func() {
		mu.Lock()
		v := m[id]
		m[id] = v - 1
		mu.Unlock()
		printTraceExit(id, name, "<-", v, time.Since(started))
	}
}
