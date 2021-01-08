package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryDeviceDAO_GetPaging_NegativeLimitPassed(t *testing.T) {
	underTest := inMemoryDeviceDAO{}

	devices, err := underTest.GetAll(context.Background(), -1, 0)

	assert.Nil(t, devices)
	assert.Error(t, err)
}

func TestInMemoryDeviceDAO_GetPaging(t *testing.T) {
	devices := []Device{
		{Id: 0, Name: "0"},
		{Id: 1, Name: "1"},
		{Id: 2, Name: "2"},
		{Id: 3, Name: "3"},
		{Id: 4, Name: "4"},
		{Id: 5, Name: "5"},
		{Id: 6, Name: "6"},
		{Id: 7, Name: "7"},
		{Id: 8, Name: "8"},
		{Id: 9, Name: "9"},
	}

	testCases := []struct {
		limit    int
		page     int
		expected []Device
	}{
		{limit: 0, page: 0, expected: devices[:]},
		{limit: 0, page: 1, expected: devices[:]},

		{limit: 1, page: 0, expected: []Device{devices[0]}},
		{limit: 1, page: 1, expected: []Device{devices[1]}},
		{limit: 1, page: 9, expected: []Device{devices[9]}},
		{limit: 1, page: 10, expected: []Device{}},

		{limit: 2, page: 0, expected: devices[0:2]},
		{limit: 2, page: 1, expected: devices[2:4]},
		{limit: 2, page: 4, expected: devices[8:10]},
		{limit: 2, page: 5, expected: []Device{}},

		{limit: 4, page: 0, expected: devices[0:4]},
		{limit: 4, page: 1, expected: devices[4:8]},
		{limit: 4, page: 2, expected: devices[8:10]},
		{limit: 4, page: 3, expected: []Device{}},

		{limit: 10, page: 0, expected: devices[:]},
		{limit: 10, page: 1, expected: []Device{}},

		{limit: 20, page: 0, expected: devices[:]},
		{limit: 20, page: 1, expected: []Device{}},
	}

	underTest := inMemoryDeviceDAO{devices: append([]Device(nil), devices...)}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("limit:%v page:%v", testCase.limit, testCase.page), func(t *testing.T) {
			result, _ := underTest.GetAll(context.Background(), testCase.limit, testCase.page)
			assert.ElementsMatch(t, testCase.expected, result)
		})
	}
}
