//go:build trace
// +build trace

package functrace

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var mu sync.Mutex
var m = make(map[uint64]int)

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func printTrace(id uint64, name, typ string, indent int, cost time.Duration) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}
	fmt.Printf("g[%02d]:%s%s%s cost:%v\n", id, indents, typ, name, cost)
}

func Trace() func() {
	start := time.Now()
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("not found caller")
	}

	id := getGID()
	fn := runtime.FuncForPC(pc)
	name := fn.Name()

	mu.Lock()
	v := m[id]
	m[id] = v + 1
	mu.Unlock()
	printTrace(id, name, "->", v+1, time.Duration(0))
	return func() {
		mu.Lock()
		v := m[id]
		m[id] = v - 1
		mu.Unlock()
		stop := time.Now()
		cost := time.Now().Sub()
		printTrace(id, name, "<-", v, cost)
	}
}
