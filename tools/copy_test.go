package tools

import (
	"fmt"
	"testing"
)

func TestPropertiesCopy(t *testing.T) {
	type other struct {
		ID   int64
		Name string
	}
	type target struct {
		ID     int64
		Name   string
		Age    int
		Amount float32
		Male   bool
		Cards  []string
		Map    map[string]string
		Other  *other
	}
	type source struct {
		ID     int
		Name   string
		Age    int
		Amount float32
		Male   bool
		Cards  []string
		Map    map[string]string
		Other  *other
	}

	tt := &target{}
	err := PropertiesCopy(tt, source{
		ID:     23,
		Name:   "Bob",
		Age:    32,
		Amount: 1000000,
		Male:   true,
		Cards:  []string{"steam", "psn", "xgp"},
		Map:    map[string]string{"bmw": "宝马", "tesla": "特斯拉", "byd": "比亚迪"},
		Other: &other{
			ID:   10,
			Name: "Lucy",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v, %+v", tt, tt.Other)
}
