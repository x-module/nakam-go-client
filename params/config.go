/**
 * Created by PhpStorm.
 * @file   config.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 14:06
 * @desc   config.go
 */

package params

const ServerKey = "defaultkey"
const ServerUrl = "http://192.168.1.187:7350"

var Nakama = struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Key     string `yaml:"key"`
	HttpKey string `yaml:"httpKey"`
}{
	Host:    "http://192.168.1.187",
	Port:    7350,
	Key:     "defaultkey",
	HttpKey: "defaulthttpkey",
}

const (
	PartyJoin  = 1
	SingleJoin = 2
)

var MatchConfig = struct {
	TotalCount int
	PartyCount int
	Delay      int
	PartySize  int
	Config     AddMatchMakerParams
}{
	TotalCount: 1003,
	PartyCount: 900,
	Delay:      30, // 延迟5秒后进入
	PartySize:  3,  // 三人队
	Config: AddMatchMakerParams{
		MinCount: 2,
		MaxCount: 30,
		Query:    "*",
		StringProperties: map[string]string{
			"MatchType":  "br",
			"GroupType":  "Squad",
			"MapTag":     "map_BattleRoyal",
			"ModeTag":    "gamemode_br",
			"LoadoutTag": "loadout_standard",
			"Region":     "NA",
			"GameType":   "GameType",
		},
		NumericProperties: map[string]float64{
			"IsLocalServer": 1,
			"MaxPlayerNum":  30,
			"MinPlayerNum":  10,
			"LobbyType":     1,
			"BotNum":        9,
		},
	},
}
