package main

import (
	"log"
	"math"
	"math/rand"

	"games.agamestudio.com/jxjydbgold/gamerule/watermargin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// SHZCardInfo .
type SHZCardInfo struct {
	Card int64
	Rate int32
}

// int64 19位整数
func main() {
	log.Println("SHZ ")
	session, err := mgo.Dial("192.168.1.166:27017")
	// session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	c := session.DB("jxjy_log_hj").C("shz_Library")
	c.EnsureIndex(mgo.Index{Key: []string{"card"}, Background: true, Sparse: true})
	c.EnsureIndex(mgo.Index{Key: []string{"rate"}, Background: true, Sparse: true})
	cards := make([]int64, 0, 1000000)
	for i := 0; i < 1000000; {
		// ToDo 随机牌型
		result := &SHZCardInfo{}
		card := Shuffle()
		var rate int32
		log.Println("Shuffle", i, card)
		// 验证牌型是否重复
		for _, v := range cards {
			if v == card {
				goto Exit
			}
		}
		cards = append(cards, card)
		i++
		// 计算倍率
		rate = CalcRate(card)
		// 插入DB
		result.Card = card
		result.Rate = rate
		err = c.Insert(result)

	Exit: // 跳过
	}
	for i := 0; i < 5000; i++ {
		var results []SHZCardInfo
		c.Find(bson.M{"rate": i}).All(&results)
		rl := len(results)
		if rl != 0 {
			log.Printf("倍率为%v 的个数为%v \n", i, rl)
		}
	}
	log.Println("END ")
}

// Shuffle .
func Shuffle() int64 {
	ret := int64(0)
	rd := int64(0)
	for i := 14; i >= 0; i-- {
		if i > 2 && i < 12 { // 3~ 11 中间三列
			rd = rand.Int63n(8) * int64(math.Pow10(i))
		} else { // 1,5列
			rd = rand.Int63n(9) * int64(math.Pow10(i))
		}
		ret += rd
	}
	return ret
}

// CalcRate .
func CalcRate(card int64) int32 {
	cards := make([]int32, 15)
	for i := 14; i >= 0; i-- {
		j := 14 - i
		cards[j] = int32((card / int64(math.Pow10(i))) % 10)
	}

	log.Println("CalcRate ", cards)
	rate, num := watermargin.CalcuScore(cards)
	log.Println("CalcuScore ", rate, num)
	return rate
}
