package main

import (
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
)

var niceMoney []uint32

type Cell struct {
	Money  uint32
	Number uint32
}

type ScratcherCard struct {
	Cells         [5][5]Cell
	WinningNumber uint32
}

type ScratcherCards struct {
	Cards []ScratcherCard
}

func generateCards(count int32, maxPrice uint32) ScratcherCards {
	smallPricesUsed := make([]uint32, len(niceMoney))

	exp := uint32(1)
	for i := int32(len(smallPricesUsed) - 1); i >= 0; i-- {
		smallPricesUsed[i] = exp
		exp *= 2
	}

	cards := make([]ScratcherCard, count)

	priceUsed := false

	for i := int32(0); i < count; i++ {
		winningNumber := uint32(0)
		numsMap := map[uint32]bool{}

		for len(numsMap) < 25 {
			num := rand.Uint32() % 100
			_, ok := numsMap[num]

			if !ok {
				numsMap[num] = true
			}
		}

		nums := [5][5]uint32{}
		money := [5][5]uint32{}
		j, k := uint32(0), uint32(0)

		for key := range numsMap {
			nums[j][k] = key
			money[j][k] = niceMoney[rand.Int()%len(niceMoney)]

			k += 1
			if k == 5 {
				k = 0
				j++
			}
		}

		if !priceUsed {
			priceUsed = true
			j, k = rand.Uint32()%5, rand.Uint32()%5
			money[j][k] = maxPrice
			winningNumber = nums[j][k]
		} else {
			containsSmallWin := false
			for l := 0; l < len(smallPricesUsed); l++ {
				if smallPricesUsed[l] != 0 {
					containsSmallWin = true

					j, k = rand.Uint32()%5, rand.Uint32()%5
					money[j][k] = niceMoney[l]

					n := uint32(0)

					for {
						n = rand.Uint32() % 100
						_, ok := numsMap[n]

						if !ok {
							break
						}
					}

					nums[j][k] = n
					winningNumber = n

					smallPricesUsed[l]--
					break
				}
			}

			if !containsSmallWin {
				n := uint32(0)

				for {
					n = rand.Uint32() % 100
					_, ok := numsMap[n]

					if !ok {
						break
					}
				}

				winningNumber = n

				if rand.Int32()%2 == 0 {
					j, k = rand.Uint32()%5, rand.Uint32()%5
					money[j][k] = maxPrice
				}
			}
		}

		cells := [5][5]Cell{}
		for i := 0; i < 5; i++ {
			for j := 0; j < 5; j++ {
				cells[i][j].Number = nums[i][j]
				cells[i][j].Money = money[i][j]
			}
		}

		cards[i] = ScratcherCard{cells, winningNumber}
	}

	return ScratcherCards{cards}
}

func errorHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func queryHandler(s string, defaultValue string) string {
	if "" == s {
		return defaultValue
	}

	return s
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := template.ParseFiles("index.html")
	errorHandler(err)

	q := r.URL.Query()

	count := int(0)
	count, err = strconv.Atoi(queryHandler(q.Get("count"), "10000"))
	errorHandler(err)

	maxPrice := int(0)
	maxPrice, err = strconv.Atoi(queryHandler(q.Get("maxPrice"), "1000000"))
	errorHandler(err)

	cards := generateCards(int32(count), uint32(maxPrice))

	err = file.Execute(w, cards)
	errorHandler(err)
}

func main() {
	niceMoney = []uint32{10, 20, 50, 100, 200, 500, 1000, 2000, 5000, 10000}

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	errorHandler(err)
}
