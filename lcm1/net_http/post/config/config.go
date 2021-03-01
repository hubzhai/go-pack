package config

// Bodypos body对应下标
type Bodypos int

const (
	// EditMemberName 修改昵称
	EditMemberName Bodypos = iota
	// AddMemberGoldByID 添加玩家保险柜金币数量
	AddMemberGoldByID
	// RestoreExchangeCoin 返还兑换订单中的金币
	RestoreExchangeCoin
	// GetMemberGoldByID 获取用户金币数量
	GetMemberGoldByID
	// CreateShortMessage 邮件
	CreateShortMessage
	// AddCoinByID addGold
	AddCoinByID
	// AddDiamondByID .
	AddDiamondByID
)

// Bodys 内容
var Bodys = []string{
	`{"Param":{"ID":201625,"Name":"zxx","Platform":"1"}}`,
	`{"Param":{"ID":206866,"Gold":1000000,"BillNo":1,"Desc":"100","Platform":"1","Tag":1,"LogType":1}}`, // `{"Param":{"ID":201625,"Gold":1000000,"BillNo":23,"Desc":"100","Platform":"1","Tag":1,"LogType":1}}`,
	`{"Param":{"ID":206866,"Coin":10000,"BillNo":2,"Platform":"1","Tag":1}}`,                            //`{"Param":{"ID":201625,"Coin":100,"BillNo":11,"Platform":"1","Tag":1}}`,
	`{"Param":{"ID":206866,"Patform":"1"}}`,                                                             //`{"Param":{"ID":201625,"Patform":"1"}}`,
	`{"Param":{"NoticeTitle":"测试","NoticeContent":"zxx","Platform":"3","SrcSnid":201625,"DestSnid":201625,"MessageType":4}}`,
	`{"Param":{"ID":206866,"Gold":100000,"BillNo":5,"Desc":"100","Platform":"1","LogType":1}}`,
	`{"Param":{"ID":203762,"Gold":2000,"GoldExt":10,"BillNo":115,"Desc":"100"}}`,
	//`{"Param":{"SnId":203678,"Rmb":3000,"Name":"400000金币","BillNo":122,"payType":"1","ShopId":"39"}}`,
	`{"Param":{"SnId":203879,"Rmb":600,"Name":"60000金币","BillNo":151,"payType":"1","ShopId":"11"}}`,
}

// Urls .
var Urls = []string{
	"http://127.0.0.1:9595/api/Member/EditMemberName?sign=d96c9982b89b206d1f34859ab4119e90&ts=111111",
	"http://127.0.0.1:9696/api/Member/AddMemberGoldById?sign=b88d7d98beb3fbbeb95a55d214d3f951&ts=111111",   //"http://127.0.0.1:9595/api/Member/AddMemberGoldById?sign=dea14c5262690ee07455a33ac92d49c4&ts=111111",
	"http://127.0.0.1:9696/api/Member/RestoreExchangeCoin?sign=64a09c1123c196b0ebbaa5d663eb6fe2&ts=111111", //"http://127.0.0.1:9595/api/Member/RestoreExchangeCoin?sign=09ecbca33bd86ff67c1a32526ae37121&ts=111111",
	"http://127.0.0.1:9696/api/Member/GetMemberGoldById?sign=8850fa545ece7d11a92655457eea0171&ts=111111",   //"http://127.0.0.1:9595/api/Member/GetMemberGoldById?sign=88c58d4333861d39a2d41b8f7472e207&ts=111111",
	"http://127.0.0.1:9595/api/Game/CreateShortMessage?sign=789eed0ef2e2f74b6ea1f234eba0a55f&ts=111111",
	"http://127.0.0.1:9696/api/Game/AddCoinById?sign=cbf31df03ae802bd24242bc79411faad&ts=111111",
	"http://127.0.0.1:9696/api/Game/AddDiamondById?sign=d3091b0cd7d2f821a2e17d02d853696f&ts=111111",
	"http://127.0.0.1:9696/api/Game/RechargeById?sign=a4c4d38eeb83d620c94c9c31eda7aa7f&ts=111111",
}

/*
AddScene
*/
