/**
 * Created by PhpStorm.
 * @file   params.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/6 10:07
 * @desc   params.go
 */

package params

import (
	"context"
	"github.com/ascii8/nakama-go"
	"sync"
)

var OnlinePartyPlayerList = sync.Map{}

// var OnlinePlayerList = map[string]Connect{}
var OnlinePlayerList = sync.Map{}
var OnlineSinglePlayerList = sync.Map{}

type Connect struct {
	Ctx              context.Context
	PlayerType       int
	Number           int            `json:"number"`
	UserId           string         `json:"user_id"`
	Username         string         `json:"username"`
	Client           *nakama.Client `json:"-"`
	Conn             *nakama.Conn   `json:"-"`
	Email            string         `json:"email"`
	MatchmakerConfig string         `json:"matchmaker_config"`
}

var PartyMembers = map[string][]string{}
var StartMatchDsService = sync.Map{}
