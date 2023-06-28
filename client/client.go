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
	"github.com/x-module/utils/utils/random"
	"github.com/x-module/utils/utils/xlog"
	params2 "go-client/params"
	"go-client/services"
	"log"
	"sync"
	"time"
)

var loginCount = 0
var group sync.WaitGroup

func getClient(num int, clientService *services.ClientService) {
	start := int(time.Now().UnixNano())
	params2.MatchConfig.PartyCount = random.RandInt(35, 55)
	params2.MatchConfig.TotalCount = params2.MatchConfig.PartyCount + random.RandInt(15, 25)
	for i := 0; i < params2.MatchConfig.PartyCount; i++ {
		go func(num, i int) {
			group.Add(1)
			email := fmt.Sprintf("new-party_%d_%d@adminadminadmin.com", num, i)
			log.Println("login:", email)
			ctx := context.WithValue(context.Background(), "email", email)

			connect := clientService.Login(ctx, email, email, params2.PartyJoin)
			params2.OnlinePartyPlayerList.Store(email, connect)
			params2.OnlinePlayerList.Store(connect.UserId, connect)
			// log.Println("login success")
			group.Done()
		}(num, i+start)
	}
	for i := 0; i < params2.MatchConfig.TotalCount-params2.MatchConfig.PartyCount; i++ {
		go func(num, i int) {
			group.Add(1)
			email := fmt.Sprintf("new-single_%d_%d@adminadminadmin.com", num, i)
			ctx := context.WithValue(context.Background(), "email", email)
			connect := clientService.Login(ctx, email, email, params2.SingleJoin)
			params2.OnlineSinglePlayerList.Store(email, connect)
			params2.OnlinePlayerList.Store(connect.UserId, connect)
			group.Done()
		}(num, i+start)
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
func run(num int, ctx context.Context) {
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
	getClient(num, clientService)
	group.Wait()
	total := 0
	params2.OnlinePlayerList.Range(func(key, value any) bool {
		total++
		return true
	})
	log.Println("all player login success,count:", total)
	// 组件party
	groupPartyIds, singlePartyIds := makeParty(clientService)
	fmt.Println("group party count:", len(groupPartyIds))
	fmt.Println("single party count:", len(singlePartyIds))
	// matchmaker
	clientService.PartyMatchmaker(groupPartyIds)

	// 检查几秒后单个匹配的开始
	xlog.Debugf("delay party start\n")
	time.Sleep(time.Duration(params2.MatchConfig.Delay) * time.Second)
	xlog.Debugf("elay party start\n")

	clientService.PartyMatchmaker(singlePartyIds)
	for {
		fmt.Print(".")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done(): // ctx.Done()返回一个只读的channel，等待上级通知退出的信号
			params2.OnlinePartyPlayerList.Range(func(key, value any) bool {
				_ = value.(params2.Connect).Client.SessionLogout(context.Background())
				return true
			})
			break
		default:
		}
	}
}
func main() {
	i := 2
	for {
		i--
		timeLength := random.RandInt(50, 100)
		go func(timeLength int) {
			ctx, cancel := context.WithCancel(context.Background())
			go run(i, ctx)
			_ = ctx
			select {
			case d := <-time.After(time.Duration(timeLength) * time.Second):
				fmt.Println("time out")
				fmt.Println("current Time :", d)
				cancel()
				i++
			}
		}(timeLength)
		time.Sleep(time.Second)

		for {
			if i < 1 {
				time.Sleep(time.Second)
			} else {
				break
			}
		}
	}
}
