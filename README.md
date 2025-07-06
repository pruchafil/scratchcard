# scratchcard

Scratchcards generator (html, go templates)

run GET request `localhost:8080?count={count}&maxPrice={price}`

default values are: 10000 for {count} and 1000000 for {price}

html contains {count} elements, first scratchcard wins {price} CZK

next 1023 scratchcards win folowing:

|count|  CZK  |
|-----|-------|
| 2^0 | 10000 |
| 2^1 |  5000 |
| 2^2 |  2000 |
| 2^3 |  1000 |
| 2^4 |   500 |
| 2^5 |   200 |
| 2^6 |   100 |
| 2^7 |    50 |
| 2^8 |    20 |
| 2^9 |    10 |
