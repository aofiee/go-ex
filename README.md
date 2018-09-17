# go-ex

## Connect DB

**db info**

- url: 192.168.99.100
- db: exchange_db
- user: root
- pass: 3nYzRaLtpM4

**phpmyadmin**

- url: http://192.168.99.100:8080
- user: root
- pass: 3nYzRaLtpM4

## Start

```bash
docker-compose up -d
```

## On Linux server

1. install make command  
   `sudo apt update && sudo apt install make -y`

2. build dockerfild  
   `make build`

3. Start DB & PHPMyadmin  
   `make up-db`

4. Start GO app  
   `make up-go`

5. Checking docker  
   `docker-compose ps`

6. Log go app  
   `docker-compose logs go-ex`
6. Log go app force real time   
   `docker-compose logs -f go-ex`

## Binance Support

```bash
API DOC
```

- url: https://github.com/binance-exchange/binance-official-api-docs

```golang
Get Depth Symbol
http.HandleFunc("/v1/binance/get/depth/", binance.GetDepthEnpoint)
```

```bash
limit between 5, 10, 20, 50, 100, 500, 1000
url: http://localhost:1234/v1/binance/get/depth/:Symbol/:limit
```

- url: http://localhost:1234/v1/binance/get/depth/ETHBTC/5

## Bittrex Support

```bash
API DOC
```

- url: https://support.bittrex.com/hc/en-us/articles/115003723911

```bash
Used to get the current tick values for a market.
- url: http://localhost:1234/v1/bittrex/get/ticker/:market
```

```bash
a string literal for the market (ex: BTC-LTC)
```

- url: http://localhost:1234/v1/bittrex/get/ticker/BTC-LTC

```golang
http.HandleFunc("/v1/bittrex/get/ticker/", bittrex.GetTickerEnpoint)
```

## Kraken Support

```bash
API DOC
```

- url: https://www.kraken.com/help/api

```bash
Used to get the current tick values for a market.
- url: http://localhost:1234/v1/kraken/get/depth/:symbol/:limit
```

```bash
a string literal for the market (ex: XBTEUR)
```

- url: http://localhost:1234/v1/kraken/get/depth/XBTEUR/19

```golang
http.HandleFunc("/v1/kraken/get/depth/", kraken.GetDepthEnpoint)
```
