package main

import (
	"bufio"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
)

type config struct {
	currency    string
	currencySet bool
	exchange    float64
	exchangeSet bool
}

func readConfig() (cfg config) {
	hd, err := homedir.Dir()
	if err != nil {
		return
	}
	cpath := path.Join(hd, ".btc")
	if _, err := os.Stat(cpath); os.IsNotExist(err) {
		return
	}
	file, err := os.Open(cpath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), "\t")
		if len(lines) < 2 {
			continue
		}
		if lines[0] == "currency" {
			cfg.currency = strings.ToUpper(lines[1])
			cfg.currencySet = true
		} else if lines[0] == "exchange" {
			f, err := strconv.ParseFloat(lines[1], 64)
			if err != nil {
				continue
			}
			cfg.exchange = f
			cfg.exchangeSet = true
		}
	}
	return
}
