/**
 * Created by PhpStorm.
 * @file   client.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 10:39
 * @desc   client.go
 */

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ascii8/nakama-go"
	"github.com/x-module/utils/global"
	"github.com/x-module/utils/utils"
	"github.com/x-module/utils/utils/request"
	"github.com/x-module/utils/utils/slice"
	"github.com/x-module/utils/utils/xlog"
	"go-client/params"
	"log"
	"strings"
)

type ClientService struct {
	handler *HandlerService
}

func NewClientService(handler *HandlerService) *ClientService {
	return &ClientService{
		handler: handler,
	}
}

// Login 玩家登录
func (c *ClientService) Login(ctx context.Context, email string, password string, playerType int) params.Connect {
	cl := nakama.New(nakama.WithServerKey(params.Nakama.Key), nakama.WithURL(fmt.Sprintf("%s:%d", params.Nakama.Host, params.Nakama.Port)))
	if err := cl.AuthenticateEmail(ctx, email, password, true, email); err != nil {
		xlog.WithField(global.ErrField, err).Fatal("login fail")
	}
	conn, err := cl.NewConn(ctx, append([]nakama.ConnOption{nakama.WithConnFormat("json")})...)
	if err != nil {
		xlog.WithField(global.ErrField, err).Fatal("create new connect fail")
	}
	if conn.Connected() == false {
		xlog.Fatal("websocket connect fail")
	}
	if err := conn.Close(); err != nil {
		xlog.WithField(global.ErrField, err).Fatal("websocket close fail")
	}
	if conn.Connected() == true {
		xlog.Fatal("websocket already connect after close ")
	}
	conn.ConnectHandler = c.handler.ConnectHandler
	conn.DisconnectHandler = c.handler.DisconnectHandler

	if err := conn.Open(ctx); err != nil {
		xlog.WithField(global.ErrField, err).Fatal("websocket connect err ")
	}
	if conn.Connected() == false {
		xlog.Fatal("websocket connect fail")
	}
	conn.MatchDataHandler = c.handler.MatchDataHandler
	conn.MatchmakerMatchedHandler = c.handler.MatchmakerMatchedHandler
	if err = conn.Open(ctx); err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}
	response, _ := cl.Account(context.Background())
	return params.Connect{
		Ctx:        ctx,
		UserId:     response.User.Id,
		PlayerType: playerType,
		Username:   response.User.Username,
		Conn:       conn,
		Client:     cl,
		Email:      email,
	}
}

func (c *ClientService) CreateParty(client params.Connect) (string, error) {
	party, err := client.Conn.PartyCreate(client.Ctx, true, 10)
	if err != nil {
		xlog.WithField(global.ErrField, err).Fatal("create party fail")
	}
	return party.PartyId, nil
}

func (*ClientService) JoinParty(client params.Connect, partyId string) error {
	err := client.Conn.PartyJoin(client.Ctx, partyId)
	if err != nil {
		xlog.WithField(global.ErrField, err).Fatal("join party fail")
	}
	return err
}

type MatchIdType string

func (m MatchIdType) isMatchmakerMatchedMsg_Id() {
	// TODO implement me
	panic("implement me")
}

func (c *ClientService) PartyMatchmaker(partyIds map[string]params.Connect) {
	for partyId, connect := range partyIds {
		go func(partyId string, connect params.Connect) {
			xlog.Debugf("party start matchmaker %s \n", partyId)
			// 尝试加入已经存在的比赛
			result, err := c.GroupMatchParty(partyId)

			if err != nil {
				return
			}
			if result.Matched > 0 { // 匹配成功
				xlog.Debugf("party match old match, partyId: %s \n", partyId)
				// utils.JsonDisplay(params.PartyMembers[partyId])
				var userList []*nakama.MatchmakerUserMsg
				for _, uid := range params.PartyMembers[partyId] {
					userList = append(userList, &nakama.MatchmakerUserMsg{
						Presence: &nakama.UserPresenceMsg{
							UserId: uid,
						},
					})
				}
				matchId := &nakama.MatchmakerMatchedMsg_MatchId{
					MatchId: result.MatchId,
				}
				c.handler.MatchmakerMatchedHandler(connect.Ctx, &nakama.MatchmakerMatchedMsg{
					Users: userList,
					Id:    matchId,
				})
			} else {
				xlog.Debug("start to join matchmaker...")
				_, err = nakama.PartyMatchmakerAdd(partyId, params.MatchConfig.Config.Query, params.MatchConfig.Config.MinCount, params.MatchConfig.Config.MaxCount).
					WithStringProperties(params.MatchConfig.Config.StringProperties).
					WithNumericProperties(params.MatchConfig.Config.NumericProperties).Send(context.Background(), connect.Conn)
				if err != nil {
					xlog.WithField(global.ErrField, err).Error("start matchmaker fail")
				} else {
					xlog.Debugf("party:%s join matchmaker success\n", partyId)
				}
			}
		}(partyId, connect)
	}
}

type PartyMatchResponse struct {
	Matched int    `json:"Matched,omitempty"` // 是否匹配成功 0 未匹配成功 1 匹配成功
	MatchId string `json:"MatchId,omitempty"` // 匹配的比赛ID
	// Data    any    `json:"Data"`
}

// GroupMatchParty 组队尝试加入已有的匹配
func (*ClientService) GroupMatchParty(partyId string) (PartyMatchResponse, error) {
	// 检查现有的是否有匹配的
	urlFormat := "%s:%d/v2/rpc/match/party?http_key=%s&unwrap"
	url := fmt.Sprintf(urlFormat, params.Nakama.Host, params.Nakama.Port, params.Nakama.HttpKey)
	matchParams := map[string]any{
		"Query":   strings.ReplaceAll(params.MatchConfig.Config.Query, "properties", "label"),
		"Members": slice.Unique(params.PartyMembers[partyId]),
	}
	res, err := request.NewRequest().Debug(true).Json().Post(url, matchParams)
	if err != nil {
		xlog.WithField(global.ErrField, err).Error("try join other match err")
		return PartyMatchResponse{}, err
	} else {
		var response PartyMatchResponse
		matchRes, _ := res.Content()
		fmt.Println("==============================match result ====================================")
		utils.JsonDisplay(matchParams)
		fmt.Println(matchRes)
		fmt.Println("==============================match result ====================================")
		err = json.Unmarshal([]byte(matchRes), &response)
		if err != nil {
			return PartyMatchResponse{}, err
		}
		return response, nil
	}
}
