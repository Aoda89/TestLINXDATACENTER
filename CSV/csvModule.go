package CSV

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func ReadCSV(name string) {
	file, err := os.Open(name)

	if err != nil {
		panic("Ошибка открытия файла")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	searchCSV(records)

}

/*
Функция для разделения задачи поиска между горутинами.Разделяем на входе данные на блоки размером 100 в случае большо-
го объема данных.Каждая горутина ищет максимальное значение в своем блоке и записывает в канал затем итерируемся по каналу
и ищем максимальное значение среди результатов работы горутин.
*/
func searchCSV(records [][]string) {
	size := len(records)
	check := size % 10

	var buffer = 0
	if check != 0 {
		buffer = size/100 + 2
	} else {
		buffer = size/100 + 1
	}
	chanelPrice := make(chan map[int]string, buffer)
	chanelRating := make(chan map[int]string, buffer)
	var wg sync.WaitGroup

	for i := 1; i <= size; i += 100 {
		wg.Add(1)
		go func(startIndex int) {
			defer wg.Done()
			var priceResult = 0
			var ratingResult = 0
			var priceName = ""
			var ratingName = ""
			MapPrice := make(map[int]string)
			MapRating := make(map[int]string)

			for j := startIndex; j < startIndex+100 && j != size; j++ {

				split := strings.Split(records[j][0], ";")

				price, err := strconv.Atoi(split[1])
				if err != nil {
					panic(err)
				}
				rating, err := strconv.Atoi(split[2])
				if err != nil {
					panic(err)
				}

				if priceResult < price {
					priceResult = price
					priceName = split[0]
				}
				if ratingResult < rating {
					ratingResult = rating
					ratingName = split[0]
				}
			}

			MapPrice[priceResult] = priceName
			MapRating[ratingResult] = ratingName

			chanelPrice <- MapPrice
			chanelRating <- MapRating
		}(i)

	}
	wg.Wait()
	close(chanelPrice)
	close(chanelRating)

	var price = 0
	var rating = 0
	var priceName = ""
	var ratingName = ""

	for priceMap := range chanelPrice {
		for key, value := range priceMap {
			if price < key {
				price = key
				priceName = value
			}
		}
	}

	for ratingMap := range chanelRating {
		for key, value := range ratingMap {
			if rating < key {
				rating = key
				ratingName = value
			}
		}
	}

	fmt.Println("Товар с наибольшей ценой:")
	fmt.Println("Имя:  ", priceName)
	fmt.Println("Цена: ", price)
	fmt.Println("Товар с наибольшим рейтингом:")
	fmt.Println("Имя:  ", ratingName)
	fmt.Println("Цена: ", rating)

}
