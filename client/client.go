/**
 * Created by PhpStorm.
 * @file   client_muilt.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/5 10:59
 * @desc   client_muilt.go
 */

package main

import (
	"context"
	"fmt"
	"github.com/x-module/utils/utils/xlog"
	params2 "go-client/params"
	"go-client/services"
	"log"
	"sync"
	"time"
)

var loginCount = 0
var group sync.WaitGroup

func getClient(clientService *services.ClientService) {
	for i := 0; i < params2.MatchConfig.PartyCount; i++ {
		go func(i int) {
			group.Add(1)
			email := fmt.Sprintf("new-party_%d@adminadminadmin.com", i)
			ctx := context.WithValue(context.Background(), "email", email)

			connect := clientService.Login(ctx, email, email, params2.PartyJoin)
			params2.OnlinePartyPlayerList.Store(email, connect)
			params2.OnlinePlayerList.Store(connect.UserId, connect)
			group.Done()
		}(i)
	}
	for i := 0; i < params2.MatchConfig.TotalCount-params2.MatchConfig.PartyCount; i++ {
		go func(i int) {
			group.Add(1)
			email := fmt.Sprintf("new-single_%d@adminadminadmin.com", i)
			ctx := context.WithValue(context.Background(), "email", email)
			connect := clientService.Login(ctx, email, email, params2.SingleJoin)
			params2.OnlineSinglePlayerList.Store(email, connect)
			params2.OnlinePlayerList.Store(connect.UserId, connect)
			group.Done()
		}(i)
	}
}
func makeParty(clientService *services.ClientService) (map[string]params2.Connect, map[string]params2.Connect) {
	i := 0
	partyId := ""
	var groupPartyIdList = map[string]params2.Connect{}
	params2.OnlinePartyPlayerList.Range(func(key, value any) bool {
		conn := value.(params2.Connect)
		if i == 0 {
			partyId, _ = clientService.CreateParty(conn)
			log.Printf("key:%s create pary\n", key)
			groupPartyIdList[partyId] = conn
		}
		_ = clientService.JoinParty(conn, partyId)
		log.Printf("key:%s join pary\n", key)
		params2.PartyMembers[partyId] = append(params2.PartyMembers[partyId], conn.UserId)
		i++
		if i == params2.MatchConfig.PartySize {
			i = 0
		}
		return true
	})
	i = 0
	var singlePartyIdList = map[string]params2.Connect{}
	params2.OnlineSinglePlayerList.Range(func(key, value any) bool {
		conn := value.(params2.Connect)
		if i == 0 {
			partyId, _ = clientService.CreateParty(conn)
			log.Printf("key:%s create pary\n", key)
			singlePartyIdList[partyId] = conn
		}
		_ = clientService.JoinParty(conn, partyId)
		log.Printf("key:%s join pary\n", key)
		params2.PartyMembers[partyId] = append(params2.PartyMembers[partyId], conn.UserId)
		i++
		if i == params2.MatchConfig.PartySize {
			i = 0
		}
		return true
	})
	return groupPartyIdList, singlePartyIdList
}

func singleJoin() {

}
func main() {
	handler := services.NewHandlerService()
	handler.ConnectHandler = func(ctx context.Context) {
		loginCount++
	}
	handler.ConnectHandler = func(ctx context.Context) {
		// e := ctx.Value("email").(string)
		// fmt.Println("=========================== email:", e)
	}
	clientService := services.NewClientService(handler)

	// 客户端登录
	getClient(clientService)
	group.Wait()
	total := 0
	params2.OnlinePlayerList.Range(func(key, value any) bool {
		total++
		return true
	})
	log.Println("all player login success,count:", total)
	// // 组件party
	groupPartyIds, singlePartyIds := makeParty(clientService)
	fmt.Println("group party count:", len(groupPartyIds))
	fmt.Println("single party count:", len(singlePartyIds))
	// matchmaker
	clientService.PartyMatchmaker(groupPartyIds)

	// 检查几秒后单个匹配的开始
	xlog.Debugf("delay party start\n")
	xlog.Debugf("delay party start\n")
	xlog.Debugf("delay party start\n")
	xlog.Debugf("delay party start\n")
	xlog.Debugf("delay party start\n")
	xlog.Debugf("delay party start\n")
	time.Sleep(time.Duration(params2.MatchConfig.Delay) * time.Second)
	xlog.Debugf("==============delay party start\n")
	xlog.Debugf("==============delay party start\n")
	xlog.Debugf("==============delay party start\n")
	xlog.Debugf("==============delay party start\n")

	clientService.PartyMatchmaker(singlePartyIds)
	for {
	}
}
