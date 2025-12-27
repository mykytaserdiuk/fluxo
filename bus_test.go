package fluxo

import (
	"errors"
	"reflect"
	"testing"

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
	bus := NewEventBus()
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
	bus := NewEventBus()
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

	bus := NewEventBus()

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

	bus := NewEventBus()
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
