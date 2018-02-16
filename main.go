package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/fatih/color"
	"github.com/praffn/btc/lib"
)

const version = "0.1.3"
const defaultCurrency = "USD"

////////////////////////////////////////////////
// Quick type to support environment currency //
////////////////////////////////////////////////
type currencyFlag struct {
	set   bool
	value string
}

func (cf *currencyFlag) Set(v string) error {
	cf.value = v
	cf.set = true
	return nil
}

func (cf *currencyFlag) String() string {
	return cf.value
}

var currency currencyFlag
var shouldLog = flag.Bool("l", false, "show log output")
var showVersion = flag.Bool("v", false, "show version")
var simple = flag.Bool("s", false, "show simple")
var difference = flag.Bool("d", true, "show difference since yesterday")

func init() {
	// setup custom type for flag parsing
	flag.Var(&currency, "c", "currency to fetch")
}

func main() {
	// parse flags
	flag.Parse()
	if *showVersion {
		fmt.Println("btc v" + version)
		return
	}
	if !*shouldLog {
		log.SetOutput(ioutil.Discard)
	}
	if !currency.set {
		currencyEnv := os.Getenv("BTC_CURRENCY")
		if len(currencyEnv) > 0 {
			currency.Set(currencyEnv)
		} else {
			currency.Set(defaultCurrency)
		}
	}

	gray := color.New(color.FgHiBlack)

	var price lib.Price
	var histRate float64
	var err error

	red := color.New(color.FgRed, color.Bold)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Writer = os.Stderr
	s.Suffix = " Fetching price for " + currency.String()
	s.Start()

	fetcher := lib.NewCoindesk(strings.ToUpper(currency.String()))
	if *difference {
		ch := make(chan lib.FetchWithHistoryResponse)
		go fetcher.FetchWithHistory(ch)
		resp := <-ch
		price = resp.Price
		histRate = resp.HistRate
		err = resp.Err
	} else {
		ch := make(chan lib.FetchResponse)
		go fetcher.Fetch(ch)
		resp := <-ch
		price = resp.Price
		err = resp.Err
	}
	if err != nil {
		red.Println("An error occured!")
		gray.Println(err.Error())
		os.Exit(1)
	}
	s.Stop()

	floatString := strconv.FormatFloat(price.Rate, 'f', 2, 64)

	if *simple {
		fmt.Println(floatString)
	} else {
		green := color.New(color.FgGreen, color.Bold)
		gray.Printf("Updated: %s\n", price.Updated.Format("2 Jan 2006 at 15:04:05"))
		fmt.Printf("%s: ", price.Currency)
		green.Printf("%s\n", floatString)
		if *difference {
			// if histRate is lower than show up percentage
			if histRate < price.Rate {
				green.Printf("⇡ ")
				fmt.Printf("%s%% since yesterday\n", strconv.FormatFloat(price.Rate/histRate, 'f', 3, 64))
			} else {
				red.Printf("⇣ ")
				fmt.Printf("%s%% since yesterday\n", strconv.FormatFloat(histRate/price.Rate, 'f', 3, 64))
			}
		}
	}
}
