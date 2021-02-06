package main

import "testing"

var samples = []struct {
	amount   string
	upper    bool
	expected string
}{
  {"8*rub", false, ""},
	{"1", true, "Один рубль 00 копеек"},
	{"100.21", false, "сто рублей 21 копейка"},
}

func Test_RuMoney(t *testing.T) {
	for _, tt := range samples {
		res := moneyRu(tt.amount, tt.upper)
		if res != tt.expected {
			t.Errorf("RuMoney(%s): expected '%s', got '%s'", tt.amount, tt.expected, res)
		}
	}
}
