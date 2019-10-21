package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type People struct {
	id    uint32
	money float32
}

func main() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		bill := map[uint32]float32{
			0: float32(rand.Intn(100)),
			1: float32(rand.Intn(200)),
			2: float32(rand.Intn(300)),
			3: float32(rand.Intn(400)),
			4: float32(rand.Intn(500)),
			5: float32(rand.Intn(600)),
		}
		fmt.Println("账单", bill)
		_, money := balance(bill)
		x, _ := json.Marshal(money)
		fmt.Println("结余", string(x))
		fmt.Println("-------------------------")
	}
}

func balance(a map[uint32]float32) (err error, bill []map[uint32]map[uint32]float32) {

	var countMoney float32
	peoples := make([]People, 0)

	for _, money := range a {
		countMoney += money
	}

	peopleCount := float32(len(a))
	averageMoney := countMoney / peopleCount
	fmt.Printf("总额: %.2f, 人数: %.0f, 分账: %.2f\n", countMoney, peopleCount, averageMoney)

	for i, payMoney := range a {
		money := payMoney - averageMoney
		if money == 0 {
			delete(a, i)
		} else {
			m, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", money), 32)
			peoples = append(peoples, People{id: i, money: float32(m)})
		}
	}

	for i := 0; i < len(peoples)-1; i++ {
		for j := i + 1; j < len(peoples); j++ {
			if peoples[i].money > peoples[j].money {
				peoples[i], peoples[j] = peoples[j], peoples[i]
			}
		}
	}
	fmt.Println("清算", peoples)

	var count float32
	for len(peoples) > 1 { //0.01
		count = peoples[0].money + peoples[len(peoples)-1].money
		m, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", count), 32)
		count = float32(m)
		if count == 0 { //给刚好
			//fmt.Printf("== %d 给 %d : %.2f 元\n",peoples[0].id,peoples[len(peoples)-1].id,-peoples[0].money)
			bill = append(bill, map[uint32]map[uint32]float32{peoples[0].id: {peoples[len(peoples)-1].id: -peoples[0].money}})
			peoples = peoples[1 : len(peoples)-1]

		} else if count > 0 { //给不够
			//fmt.Printf(">0 %d 给 %d : %.2f 元\n",peoples[0].id,peoples[len(peoples)-1].id,-peoples[0].money)
			bill = append(bill, map[uint32]map[uint32]float32{peoples[0].id: {peoples[len(peoples)-1].id: -peoples[0].money}})
			peoples[len(peoples)-1].money = count
			peoples = peoples[1:]

		} else { //给有剩
			//fmt.Printf("<0 %d 给 %d : %.2f 元\n",peoples[0].id,peoples[len(peoples)-1].id,peoples[len(peoples)-1].money)
			bill = append(bill, map[uint32]map[uint32]float32{peoples[0].id: {peoples[len(peoples)-1].id: peoples[len(peoples)-1].money}})
			peoples[0].money = count
			peoples = peoples[0 : len(peoples)-1]
		}
	}
	return
}
