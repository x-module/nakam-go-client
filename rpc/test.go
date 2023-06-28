package main

import (
	"context"
	"fmt"
	"github.com/x-module/utils/utils"
	"go-client/pkg/nakama"
	"log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// create client
	cl := nakama.New(
		nakama.WithServerKey("defaultkey"),
		nakama.WithURL("http://192.168.1.187:7350"))
	if err := cl.AuthenticateEmail(ctx, "aabds@aa.com", "aasdfasdfasdfasdfasdaa", true, ""); err != nil {
		log.Fatal(err)
	}

	account, _ := cl.Account(ctx)
	utils.JsonDisplay(account)

	conn, err := cl.NewConn(ctx, append([]nakama.ConnOption{nakama.WithConnFormat("json")})...)
	if err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}

	if conn.Connected() == false {
		log.Fatalf("expected conn.Connected() == true")
	}

	conn.PartyPresenceEventHandler = func(ctx context.Context, msg *nakama.PartyPresenceEventMsg) {
		fmt.Printf("recive PartyPresenceEventHandler msg:%+v", msg)
	}
	conn.StreamPresenceEventHandler = func(ctx context.Context, msg *nakama.StreamPresenceEventMsg) {
		fmt.Printf("recive StreamPresenceEventHandler msg:%+v", msg)
	}

	conn.StreamDataHandler = func(ctx context.Context, msg *nakama.StreamDataMsg) {
		fmt.Printf("recive StreamDataHandler msg:%+v", msg)
	}
	conn.StatusPresenceEventHandler = func(ctx context.Context, msg *nakama.StatusPresenceEventMsg) {
		fmt.Printf("recive StatusPresenceEventHandler msg:%+v", msg)
	}

	conn.NotificationsHandler = func(ctx context.Context, msg *nakama.NotificationsMsg) {
		fmt.Printf("recive NotificationsHandler msg:%+v\n", msg)
	}
	conn.MatchmakerMatchedHandler = func(ctx context.Context, msg *nakama.MatchmakerMatchedMsg) {
		fmt.Printf("recive MatchmakerMatchedHandler msg:%+v\n", msg)
	}
	conn.MatchPresenceEventHandler = func(ctx context.Context, msg *nakama.MatchPresenceEventMsg) {
		fmt.Printf("recive MatchPresenceEventHandler msg:%+v", msg)
	}
	conn.ChannelPresenceEventHandler = func(ctx context.Context, msg *nakama.ChannelPresenceEventMsg) {
		fmt.Printf("recive ChannelPresenceEventHandler msg:%+v", msg)
	}
	conn.ChannelMessageHandler = func(ctx context.Context, msg *nakama.ChannelMessageMsg) {
		fmt.Printf("recive ChannelMessageHandler msg:%+v\n", msg)
	}
	conn.ErrorHandler = func(ctx context.Context, msg *nakama.ErrorMsg) {
		fmt.Printf("recive ErrorHandler msg:%+v", msg)
	}
	conn.ConnectHandler = func(context.Context) {
		log.Printf("connected")
	}
	conn.DisconnectHandler = func(_ context.Context, err error) {
		log.Printf("disconnected: %v", err)
	}
	conn.MatchDataHandler = func(_ context.Context, msg *nakama.MatchDataMsg) {
		fmt.Printf("recive msg:%+v", msg)
	}
	if err = conn.Open(ctx); err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}

	var result any
	// err = conn.Rpc(ctx, "player/attribute", nil, &result)
	// err = conn.Rpc(ctx, "player/equipment", nil, &result)
	// err = conn.Rpc(ctx, "config", nil, &result)
	// err = conn.Rpc(ctx, "goods", map[string]any{}, &result)
	// err = conn.Rpc(ctx, "player/charge", map[string]any{
	// 	"GoldCoin": 100,
	// }, &result)
	//
	err = conn.Rpc(ctx, "buy", map[string]any{
		"GoodsId":     12, // 商品ID
		"PaymentType": 1,  // 1 金币 2 代币
	}, &result)

	// err = conn.Rpc(ctx, "badge/set", map[string]any{
	// 	"BadgeId":   5,
	// 	"BadgeName": "世界规则",
	// 	"BadgeIcon": "asdfasdfasdfasdf", // 1 金币 2 代币
	// }, &result)
	//
	// err = conn.Rpc(ctx, "badge/get", map[string]any{
	// 	"pageIndex": 1,
	// 	"PlayerId":  "2969d8a5-a918-46f0-8138-746e8bba481e",
	// 	"pageSize":  10,
	// }, &result)

	if err != nil {
		panic(err)
	}

	utils.JsonDisplay(result)

	for {
	}
}
