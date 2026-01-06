package fluxo

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/mykytaserdiuk/fluxo/pool"
	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	cases := []struct {
		name string

		callback      any
		finalCount    int
		expectedError error
	}{
		{
			name: "OK",
			callback: func(count *int) {
				*count++
			},
			finalCount: 1,
		}, {
			name: "OK",
			callback: func(count *int) {
				*count++
			},
			// After that case, will 2 times counter increase
			finalCount: 3,
		},
		{
			name:          "Wrong function",
			callback:      errors.New("wrong fn"),
			expectedError: ErrNotAFunc,
		},
	}
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)
	id := "sub_once"
	count := 0

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := bus.Subscribe(id, c.callback)
			assert.Equal(t, c.expectedError, err)
			bus.Emit(id, &count)
			if err == nil {
				assert.Equal(t, c.finalCount, count)
			}
		})
	}
}

func TestSubscribeOnce(t *testing.T) {
	cases := []struct {
		name       string
		finalCount int
	}{
		{
			name:       "OK_One_Emit_accept",
			finalCount: 1,
		},
		{
			name:       "OK_Without_Emit",
			finalCount: 0,
		},
	}
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	id := "sub_once"
	count := 0
	fn := func() {
		count++
	}

	err := bus.SubscribeOnce(id, fn)
	assert.NoError(t, err)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bus.Emit(id)
			bus.Emit(id)
			assert.EqualValues(t, c.finalCount, count)
		})
		count = 0
	}
}

func TestUnsubscribe(t *testing.T) {
	cases := []struct {
		name          string
		sub           bool
		expectedError error
	}{
		{
			name: "OK",
			sub:  true,
		},
		{
			name:          "NoSub",
			sub:           false,
			expectedError: ErrNoHandlers,
		},
	}
	id := "ubsub"
	call := func(t *testing.T) { t.Error("emit was succefull") }

	for i, c := range cases {
		pool := pool.NewCorePool()
		bus := NewEventBus(pool)
		t.Run(c.name, func(t *testing.T) {
			nid := id + fmt.Sprint(i)
			if c.sub {
				err := bus.Subscribe(nid, call)
				assert.NoError(t, err)
			}
			err := bus.Unsubscribe(nid, call)
			assert.Equal(t, c.expectedError, err)
			assert.Equal(t, 0, len(bus.(*EventBus).subs[nid]))
			bus.Emit(nid, t)
		})
	}
}

func TestUnregister(t *testing.T) {
	id := "test_id"
	type behaviour func(bus Bus, id string)
	cases := []struct {
		name          string
		behaviour     behaviour
		expectedError error
	}{{
		name: "OK",
		behaviour: func(bus Bus, id string) {
			bus.Subscribe(id, func() {})
		},
	}, {
		name:          "No_Subs",
		behaviour:     func(bus Bus, id string) {},
		expectedError: ErrNoHandlers,
	}, {
		name: "Many_Subs",
		behaviour: func(bus Bus, id string) {
			bus.Subscribe(id, func() {})
			bus.Subscribe(id, func() {})
			bus.Subscribe(id, func() {})
		},
	}}
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)
	for i, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			newID := fmt.Sprintf("%s_%v", id, i)
			cs.behaviour(bus, newID)
			err := bus.Unregister(newID)
			assert.Equal(t, cs.expectedError, err)
			assert.Empty(t, len(bus.(*EventBus).subs[newID]))
		})
	}
}

func TestArgsToValues(t *testing.T) {
	cases := []struct {
		name     string
		args     []any
		expected int
	}{
		{
			name:     "NoArgs",
			args:     []any{},
			expected: 0,
		},
		{
			name:     "SingleArg",
			args:     []any{42},
			expected: 1,
		},
		{
			name:     "MultipleArgs",
			args:     []any{42, "hello", 3.14, true},
			expected: 4,
		},
		{
			name:     "MixedTypes",
			args:     []any{1, "test", nil},
			expected: 3,
		},
	}
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := bus.(*EventBus).argsToValues(c.args...)
			assert.EqualValues(t, c.expected, len(result))
			for i, v := range result {
				assert.Equal(t, reflect.ValueOf(c.args[i]), v)
			}
		})
	}
}

func TestPrepareArguments(t *testing.T) {
	name1 := "john"
	name2 := "smith"
	cases := []struct {
		name           string
		callback       any
		args           []any
		expectedResult []reflect.Value
	}{{
		name:           "OK",
		callback:       func(name string) {},
		args:           []any{name1},
		expectedResult: []reflect.Value{reflect.ValueOf(name1)},
	}, {
		name:           "OK_1_arg_more",
		callback:       func(name, second string) {},
		args:           []any{name1},
		expectedResult: []reflect.Value{reflect.ValueOf(name1), reflect.Zero(reflect.TypeOf(""))},
	}, {
		name:           "OK_1_arg_less",
		callback:       func(name string) {},
		args:           []any{name1, name2},
		expectedResult: []reflect.Value{reflect.ValueOf(name1)},
	}}

	pool := pool.NewCorePool()
	bus := NewEventBus(pool)
	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			res := bus.(*EventBus).prepareArguments(handler{
				callback: reflect.ValueOf(cs.callback),
				once:     false,
			}, cs.args...)
			assert.Equal(t, cs.expectedResult, res)
		})
	}
}

func BenchmarkEmit_WithArguments(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.Subscribe("event", func(count *int) {
		*count++
	})
	c := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event", &c)
	}
}

func BenchmarkEmit_NoArguments(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.Subscribe("event", func() {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event")
	}
}

func BenchmarkSubscribe(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	fn := func(count *int) { *count++ }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Subscribe("event", fn)
	}
}

func BenchmarkMultipleSubscribers(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.SubscribeAsync("event", func(count *int) { *count++ })
	bus.SubscribeAsync("event", func(count *int) { *count++ })
	bus.SubscribeAsync("event", func(count *int) { *count++ })
	c := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event", &c)
	}
}

func BenchmarkEmitAsync_WithArguments(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.SubscribeAsync("event", func(count *int) {
		*count++
	})
	c := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event", &c)
	}
}

func BenchmarkEmit_NoArgumentsAsync(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.SubscribeAsync("event", func() {})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event")
	}
}

func BenchmarkSubscribeAsync(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	fn := func(count *int) { *count++ }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.SubscribeAsync("event", fn)
	}
}

func BenchmarkMultipleSubscribersAsync(b *testing.B) {
	pool := pool.NewCorePool()
	bus := NewEventBus(pool)

	bus.SubscribeAsync("event", func(count *int) { *count++ })
	bus.SubscribeAsync("event", func(count *int) { *count++ })
	bus.SubscribeAsync("event", func(count *int) { *count++ })
	c := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bus.Emit("event", &c)
	}
}

// Benchmark result
// NO POOL
//
// goos: windows
// goarch: amd64
// pkg: github.com/mykytaserdiuk/fluxo
// cpu: 11th Gen Intel(R) Core(TM) i5-11400H @ 2.70GHz
// BenchmarkEmit_WithArguments-12           6618330               178.0 ns/op            24 B/op          1 allocs/op
// BenchmarkEmit_NoArguments-12             9132231               128.3 ns/op             0 B/op          0 allocs/op
// BenchmarkSubscribe-12                   13519893                87.46 ns/op          186 B/op          0 allocs/op
// BenchmarkMultipleSubscribers-12          2244873               523.8 ns/op           168 B/op          4 allocs/op
// PASS
// ok      github.com/mykytaserdiuk/fluxo  7.151s

// WITH POOL
// pkg: github.com/mykytaserdiuk/fluxo
// cpu: 11th Gen Intel(R) Core(TM) i5-11400H @ 2.70GHz

// BenchmarkEmit_WithArguments-12          	 5742362	       212.5 ns/op	      40 B/op	       2 allocs/op
// BenchmarkEmit_NoArguments-12            	 8999967	       134.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkSubscribe-12                   	14192054	        71.14 ns/op	     177 B/op	       0 allocs/op
// BenchmarkMultipleSubscribers-12         	  794632	      1425 ns/op	     424 B/op	       8 allocs/op

// BenchmarkEmitAsync_WithArguments-12     	 2318260	       509.1 ns/op	     120 B/op	       3 allocs/op
// BenchmarkEmit_NoArgumentsAsync-12       	 2567932	       492.4 ns/op	      80 B/op	       1 allocs/op
// BenchmarkSubscribeAsync-12              	15988488	        66.65 ns/op	     197 B/op	       0 allocs/op
// BenchmarkMultipleSubscribersAsync-12    	  746212	      1442 ns/op	     424 B/op	       8 allocs/op
