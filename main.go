package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ascii8/nakama-go"
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
	// // retrieve account
	// account, err := cl.Account(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("%+v", account)

	conn, err := cl.NewConn(ctx, append([]nakama.ConnOption{nakama.WithConnFormat("json")})...)
	if err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}

	if conn.Connected() == false {
		log.Fatalf("expected conn.Connected() == true")
	}
	if err := conn.Close(); err != nil {
		log.Fatalf("expected on error, got: %v", err)
	}
	if conn.Connected() == true {
		log.Fatalf("expected conn.Connected() == false")
	}
	connectCh := make(chan bool, 1)
	conn.ConnectHandler = func(context.Context) {
		log.Printf("connected")
		connectCh <- true
	}
	disconnectCh := make(chan error, 1)
	conn.DisconnectHandler = func(_ context.Context, err error) {
		log.Printf("disconnected: %v", err)
		disconnectCh <- err
	}
	if err := conn.Open(ctx); err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}
	if conn.Connected() == false {
		log.Fatalf("expected conn.Connected() == true")
	}

	if conn.Connected() != true {
		log.Fatalf("expected conn.Connected() == true")
	}

	dataCh := make(chan *nakama.MatchDataMsg, 1)
	conn.MatchDataHandler = func(_ context.Context, msg *nakama.MatchDataMsg) {
		fmt.Printf("recive msg:%+v", msg)
		dataCh <- msg
	}
	conn.MatchmakerMatchedHandler = func(ctx context.Context, msg *nakama.MatchmakerMatchedMsg) {
		b, _ := json.Marshal(*msg)
		fmt.Printf("recive MatchmakerMatchedHandler msg:%+v\n", string(b))
	}

	if err := conn.Open(ctx); err != nil {
		log.Fatalf("expected no error, got: %v", err)
	}
	create, err := conn.PartyCreate(ctx, true, 10)
	if err != nil {
		log.Fatalf("PartyCreate error, got: %v", err)
	}
	fmt.Println("party id:", create.GetPartyId())
	if err != nil {
		log.Fatalf("MatchmakerAdd error, got: %v", err)
	}

	matchmakerAdd, err := conn.PartyMatchmakerAdd(ctx, create.PartyId, "*", 2, 10)
	if err != nil {
		log.Fatalf("PartyMatchmakerAdd error, got: %v", err)
	}

	fmt.Println("-----------:", matchmakerAdd.String())

	select {
	case <-ctx.Done():
		log.Fatalf("context closed: %v", ctx.Err())
	case msg := <-dataCh:
		if s, exp := string(msg.Data), "hello world"; s != exp {
			log.Printf("expected %q, got: %q", exp, s)
		}
	}

	for {
	}
}
