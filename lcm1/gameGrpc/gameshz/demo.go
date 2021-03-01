package gamegrpc

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// GSRPC .
type GSRPC struct {
	Data    map[int32]*GamePool
	db      *mgo.Collection
	ShzCard map[int32][]int64
}

const (
	MaxZeroRateCard   = 10000
	MaxNoZeroRateCard = 2000
)

// GamePool 水池
type GamePool struct {
	Coin int64
}

// SHZCardInfo .
type SHZCardInfo struct {
	Card int64
	Rate int32
}

// UserData .
type UserData struct {
	GamefreeID int32
	Bet        int64
	ChangeCoin int64
}

// WinData .
type WinData struct {
	RandFloat float64
	MaxRate   int64
	MinRate   int64
}

// Prints .
func (gs *GSRPC) Prints(id *int32, ret *string) error {
	log.Println("GSRPC", gs.Data[*id])
	*ret = "true"
	return nil
}

// InitDB .
func (gs *GSRPC) InitDB(url *string, ret *int32) error {
	log.Println("InitDB")
	session, err := mgo.Dial(*url)
	if err != nil {
		return err
	}

	// defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	gs.db = session.DB("jxjy_log_hj").C("shz_Library") // 先写死
	for i := 0; i < 500; i++ {
		var results []SHZCardInfo
		gs.db.Find(bson.M{"rate": i}).All(&results)
		rl := len(results)
		if rl != 0 {
			if results[0].Rate == 0 {
				if rl > MaxZeroRateCard {
					cards := make([]int64, 0, MaxZeroRateCard) //TODO 随机
					for i := 0; i < MaxZeroRateCard; i++ {
						cards = append(cards, results[i].Card)
					}
					gs.ShzCard[results[0].Rate] = cards
				} else {
					cards := make([]int64, 0, rl)
					for _, v := range results {
						cards = append(cards, v.Card)
					}
					gs.ShzCard[results[0].Rate] = cards
				}
			} else {
				if rl > MaxNoZeroRateCard {
					cards := make([]int64, 0, MaxNoZeroRateCard) //TODO 随机
					for i := 0; i < MaxNoZeroRateCard; i++ {
						cards = append(cards, results[i].Card)
					}
					gs.ShzCard[results[0].Rate] = cards
				} else {
					cards := make([]int64, 0, rl)
					for _, v := range results {
						cards = append(cards, v.Card)
					}
					gs.ShzCard[results[0].Rate] = cards
				}
			}
			//log.Printf("倍率为%v 的个数为%v \n", i, rl)
			*ret = results[0].Rate
		}
	}
	return nil
}

// RetCard .
func (gs *GSRPC) RetCard(rate *int32, card *int64) error {
	if v, exist := gs.ShzCard[*rate]; exist {
		*card = v[rand.Intn(len(v))]
	} else {
		*card = 0
	}
	return nil
}

// 0 不中奖 1 2~8

// CalcGive .
func (gs *GSRPC) CalcGive(p *UserData, ret *WinData) error {

	if v, exist := gs.Data[p.GamefreeID]; exist {
		// calc
		if v.Coin < p.Bet {
			ret.RandFloat = 0.0
			return nil
		}
		cf := p.Bet * 100 / v.Coin
		if cf <= 5 { // 最好走配置表修改
			ret.RandFloat = 0.3
			ret.MaxRate = 8
			ret.MinRate = 2
			return nil
		}
		rand.Float64()
		// TODO 使用个人分控

		log.Println("CalcGive", *ret, v)
	} else {
		ret.RandFloat = 0.0
	}

	return nil
}

// ChangePool .
func (gs *GSRPC) ChangePool(p *UserData, ret *bool) error {
	*ret = true
	if v, exist := gs.Data[p.GamefreeID]; exist {
		// calc
		v.Coin += p.ChangeCoin
		log.Println("CalcGive 0 ", time.Now().UnixNano(), v)
	} else {
		pool := &GamePool{
			Coin: p.ChangeCoin,
		}
		gs.Data[p.GamefreeID] = pool
		log.Println("CalcGive 1 ", time.Now().UnixNano(), v)
	}
	return nil
}
