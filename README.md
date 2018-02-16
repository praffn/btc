# $ btc
## super simple cli tool to quickly check current BTC rates

### Installation
```
go get github.com/praffn/btc
```

### Usage
Simply run `btc` to get the current price of BTC. The currency will default to USD or, if set, the value of the environment variable `BTC_CURRENCY`.

#### Specify currency
```
$ btc -c <CODE>
```
Where `CODE` is any valid ISO 4217 currency and is supported by Coindesk. [Link to Coindesk's supported currency codes](https://api.coindesk.com/v1/bpi/supported-currencies.json)

-------------

btc is Powered by [CoinDesk](https://www.coindesk.com/price/)

<small>btc and the author is not affiliated, associated, authorized, endorsed by, or in any way officially connected with Coindesk</small>