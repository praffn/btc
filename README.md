# > ₿TC
**btc** is a cli tool for fetching the latest bitcoin (BTC) rate from [CoinDesk](https://www.coindesk.com/)

![GIF of btc](btc.gif)

### Installation
No binaries yet, but you can easily fetch btc if you have Go installed:
```
go get github.com/praffn/btc
```

### Flags

* `-c <CODE>`<br>
  *Where CODE is a valid ISO 4217 currency code and is in [CoinDesk's list of supported currencies](https://api.coindesk.com/v1/bpi/supported-currencies.json).*

  Set the currency to show the BTC rate. Defaults to USD. Automatically uppercases value so "dkk" is equivalent to "DKK"

* `-x <AMOUNT>`<br>
  *Where amount is a float*

  If set, the price for AMOUNT of Bitcoin will be shown (in chosen currency)

  For example `btc -c dkk -x 0.3` will show the price for 0.3 BTC in DKK

* `-d`<br>
  Shows the difference since yesterday in percentages. On by default. To disable: `-d=0` or `-d=false`

* `-s`<br>
  If toggled **ONLY** the float value of a single Bitcoin will be printed to stdout. All flags are ignored expect -c and -l. Useful for piping the result.

* `-l`<br>
  If set, will print log information. Only useful for debuggin purposes

* `-v`<br>
  If set, will print version information and terminate

* `-h`<br>
  Show help

-------------

**btc** is Powered by [CoinDesk](https://www.coindesk.com/price/)

**btc** and the author is not affiliated, associated, authorized, endorsed by, or in any way officially connected with CoinDesk