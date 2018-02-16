package main

import (
	"fmt"
	"strconv"
	"strings"
)

type currencyFlag struct {
	set   bool
	value string
}

func (cf *currencyFlag) Set(v string) error {
	cf.value = strings.ToUpper(v)
	cf.set = true
	return nil
}

func (cf *currencyFlag) String() string {
	return cf.value
}

type exchangeFlag struct {
	set   bool
	value float64
}

func (ef *exchangeFlag) Set(v string) error {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("value must be a float")
	}
	ef.value = f
	ef.set = true
	return nil
}

func (ef *exchangeFlag) String() string {
	return strconv.FormatFloat(ef.value, 'f', -1, 64)
}
