package JSON

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Magazine struct {
	Product string `json:"product"`
	Price   int    `json:"price"`
	Rating  int    `json:"rating"`
}

func ReadJSON(name string) {
	var instance []Magazine
	file, err := os.Open(name)

	if err != nil {
		panic("Ошибка открытия файла")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &instance)
	if err != nil {
		panic("Ошибка декодирования данных")
	}

	searchJSON(instance)

}

/*
Функция для разделения задачи поиска между горутинами.Разделяем на входе данные на блоки размером 100 в случае большо-
го объема данных.Каждая горутина ищет максимальное значение в своем блоке и записывает в канал затем итерируемся по каналу
и ищем максимальное значение среди результатов работы горутин.
*/
func searchJSON(instance []Magazine) {
	size := len(instance)
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

	for i := 0; i <= size; i += 100 {
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
				if priceResult < instance[j].Price {
					priceResult = instance[j].Price
					priceName = instance[j].Product
				}
				if ratingResult < instance[j].Rating {
					ratingResult = instance[j].Rating
					ratingName = instance[j].Product
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
