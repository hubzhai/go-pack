package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

/*
文档
http://godoc.org/labix.org/v2/mgo

获取

go get gopkg.in/mgo.v2
连接

session, err := mgo.Dial(url)
切换数据库

db := session.DB("test")
切换集合

通过Database.C()方法切换集合（Collection）。

func (db Database) C(name string) Collection
插入

func (c *Collection) Insert(docs ...interface{}) error
c := session.DB("store").C("books")
err = c.Insert(book)
查询

func (c Collection) Find(query interface{}) Query
更新

c := session.DB("store").C("books")
err = c.Update(bson.M{"isbn": isbn}, &book)
查询所有

c := session.DB("store").C("books")

var books []Book
err := c.Find(bson.M{}).All(&books)

删除

c := session.DB("store").C("books")
err := c.Remove(bson.M{"isbn": isbn})
*/

// Person 人物数据库信息
type Person struct {
	Name  string
	Phone string
}

// PdkCardInfo .
type PdkCardInfo struct {
	ID         int
	WinCard    []int
	LoserCards [][]int
	WinPos     int
}

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// c := session.DB("jxjy_game_h5").C("pdk_winLibrary")
	c := session.DB("test").C("pdkWinLibrary")
	result := &PdkCardInfo{}

	for i, v := range WinCards {
		if err = c.Find(bson.M{"id": i}).One(result); err == mgo.ErrNotFound { // 没有找到就插入
			log.Println("insert", i, v)
			err = c.Insert(&PdkCardInfo{i, v, LoserCards[i], winPos[i]})
		} else {
			err = c.Update(bson.M{"id": result.ID}, &PdkCardInfo{result.ID, v, LoserCards[i], winPos[i]})
			log.Println("Update", result.ID, v)
		}
		if err != nil {
			log.Fatal(result.ID, err)
		}
	}
	var winCards [][]int
	var loseCards [][][]int
	var winPos []int
	var results []PdkCardInfo
	c.Find(bson.M{}).All(&results)
	winCards = make([][]int, len(results))
	loseCards = make([][][]int, len(results))
	for i, v := range results {
		winCards[i] = append(winCards[i], v.WinCard...)
		loseCards[i] = append(loseCards[i], v.LoserCards...)
		winPos = append(winPos, v.WinPos)
	}
	for k, v := range winCards {
		fmt.Printf("{")
		for i, v1 := range v {
			if i != 15 {
				fmt.Printf("%d, ", v1)
			} else {
				fmt.Printf("%d", v1)
			}
		}
		fmt.Printf("}, // %d\n", k)
	}
	for k, v := range loseCards {
		fmt.Printf("{")
		for _, v1 := range v {
			fmt.Printf("{")
			for i, v2 := range v1 {
				if i != 15 {
					fmt.Printf("%d, ", v2)
				} else {
					fmt.Printf("%d", v2)
				}

			}
			fmt.Printf("},")
		}
		fmt.Printf("}, // %d\n", k)
	}
	fmt.Printf("{\n")
	for k, v := range winPos {
		fmt.Printf("%d, ", v)
		fmt.Printf(" // %d\n", k)
	}
	fmt.Printf("}\n")
}
