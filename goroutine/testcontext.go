package main

import (
	"time"
	"sync"
	"fmt"
	"context"
)

type Context interface {
	//当为该context 工作的work被取消时，deadline 将返回时间。在没有设定期限的情况下，
	//deadline 返回 ok == false. 联系的调用deadline 返回相同的结果
	Deadline() (deadline time.Time,ok bool)

	//当为该context 工作的work完成时，返回一个关闭的channel 。
	//如果这个context 不能被取消，那么Done 可能返回nil 。连续调用完成返回相同的值
	Done() <-chan struct{}

	//Err 在完成后返回一个 non-nil 值。如果context 被取消，或者在context 的deadline结束时，
	//如果context 被取消，Err 将被取消。没有定义Err 的其他值。连续调用结束后，用Err 返回相同的值
	Err() error

	//连续调用具有相同key 的值将返回相同的结果
	Value(key interface{}) interface{}
}


func main() {
	fmt.Println("sayhelloandbye")
	sayhelloandbye()
	fmt.Println()
	fmt.Println("contexthelloandbye")
	contexthelloandbye()
	fmt.Println()
	fmt.Println("ctxdeadline")
	ctxdeadline()
	fmt.Println()
	fmt.Println("ctxvalue")
	ctxvalue()
}

//同时打印问候和告别
func sayhelloandbye(){
	locale := func(done <-chan interface{}) (string,error){
		select {
		case <- done:
			return "",fmt.Errorf("canceled")
		case <-time.After(500*time.Millisecond):
		}
		return "EN/US",nil
	}

	genGreeting := func(done <-chan interface{}) (string,error) {
		switch locale,err := locale(done); {
		case err != nil:
			return "",err;
		case locale == "EN/US":
			return "hello",nil;
		}
		return "",fmt.Errorf("unsupported locale")
	}

	genFarewell := func(done <-chan interface{}) (string,error){
		switch locale,err := locale(done); {
		case err != nil:
			return "",err;
		case locale == "EN/US":
			return "goodbye",nil;
		}
		return "",fmt.Errorf("unsupported locale")
	}

	printGreeting := func(done <-chan interface{}) error {
		greeting , err := genGreeting(done)
		if err != nil{
			return err
		}
		fmt.Printf("%s world!\n",greeting)
		return nil
	}

	printFarewell := func(done <-chan interface{}) error {
		farewell,err := genFarewell(done)
		if err != nil {
			return err
		}
		fmt.Printf("%s world!\n",farewell)
		return nil
	}


	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(done);err != nil{
			fmt.Printf("%v",err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(done); err != nil{
			fmt.Printf("%v",err)
			return
		}
	}()

	wg.Wait()

}

//context包同步问候和告别
func contexthelloandbye(){

	locale := func(ctx context.Context) (string,error) {
		select {
		case <-ctx.Done():
			return "",ctx.Err()
		case <-time.After(1*time.Minute):
		}
		return "EN/US",nil
	}

	genFarewell := func(ctx context.Context) (string,error) {
		switch local,err := locale(ctx); {
		case err != nil:
			return "",err;
		case local == "EN/US":
			return "goodbye",nil
		}
		return "",fmt.Errorf("unsupported locale")
	}

	genGreeting := func(ctx context.Context) (string,error) {
		ctx,cancel := context.WithTimeout(ctx,500*time.Millisecond)
		defer cancel()
		switch local,err := locale(ctx); {
		case err != nil:
			return "",err;
		case local == "EN/US":
			return "hello",nil
		}
		return "",fmt.Errorf("unsupported locale")
	}

	printGreeting := func(ctx context.Context) error {
		greeting,err := genGreeting(ctx)
		if err != nil{
			return err
		}
		fmt.Printf("%s world!\n", greeting)
		return nil
	}

	printFarewell := func(ctx context.Context) error {
		farewell,err := genFarewell(ctx)
		if err != nil{
			return err
		}
		fmt.Printf("%s world\n", farewell)
		return nil
	}

	var wg sync.WaitGroup
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx);err != nil{
			fmt.Printf("cannot print greeting: %v\n",err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx);err != nil{
			fmt.Printf("cannot print farewell: %v\n",err)
		}
	}()
	wg.Wait()

	/*
	·ctx,cancel := context.WithCancel(context.Background())
	这里main 用context.BackGround() 创建一个Context ，并用context.WithCancel 包装它以允许取消

	·cancel()
	在这里，如果从打印问候语返回错误，main 将取消上下文.

	·ctx,cancel := context.WithTimeout(ctx,500*time.Millisecond)
	这里genGreeting 用context.WithTimeout 包装它的 Context
	这将会在 500 毫秒后自动取消返回的上下文，从而取消它传递上下文的任何子函数，即语言环境

	·return "",ctx.Err()
	这一行返回为什么Context 被取消的原因，这个错误会一直冒泡到main ，这会导致main 的sync.WaitGroup 也取消。

	*/

}

//使用context deadline
func ctxdeadline(){
	locale := func(ctx context.Context) (string,error) {
		if deadline,ok := ctx.Deadline();ok{
			if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0{
				return "",context.DeadlineExceeded
			}
		}
		select {
		case <-ctx.Done():
			return "",ctx.Err();
		case <-time.After(1*time.Minute):
		}
		return "EN/US",nil
	}

	genFarewell := func(ctx context.Context) (string,error) {
		switch local,err := locale(ctx); {
		case err != nil:
			return "",err
		case local == "EN/US":
			return "goodbye",nil;
		}
		return "",fmt.Errorf("unsupported locale")
	}

	genGreeting := func(ctx context.Context) (string,error) {
		ctx,cancel := context.WithTimeout(ctx,500*time.Millisecond)
		defer cancel()
		switch local,err := locale(ctx); {
		case err != nil:
			return "",err;
		case local == "EN/US":
			return "hello",nil
		}
		return "",fmt.Errorf("unsupported locale")
	}

	printGreeting := func(ctx context.Context) error {
		greeting,err := genGreeting(ctx)
		if err != nil {return err}
		fmt.Printf("%s world!\n",greeting)
		return nil
	}

	printFarewell := func(ctx context.Context) error{
		farewell , err := genFarewell(ctx)
		if err != nil{ return err}
		fmt.Printf("%s world!\n", farewell)
		return nil
	}

	var wg sync.WaitGroup

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreeting(ctx); err != nil{
			fmt.Printf("cannot print greeting: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell(ctx);err != nil{
			fmt.Printf("cannot print farewell: %v\n",err)
		}
	}()

	wg.Wait()

	/*
	·if deadline,ok := ctx.Deadline();ok{
	在这里我们检查我们的上下文是否提供了截止。
	如果确实如此，并且我们的系统时钟已超过截止时间，那么我们只会返回上下文包中定义的特定错误，即DeadlineExceeded
	*/
}

//context value
func ctxvalue(){
	HandleResponse := func(ctx context.Context) {
		fmt.Printf(
			"handling response for %v (%v)\n",
			ctx.Value("userID"),
			ctx.Value("authToken"),
			)
	}

	ProcessRequest := func(userID,authToken string) {
		ctx := context.WithValue(context.Background(),"userID",userID)
		ctx = context.WithValue(ctx,"authToken",authToken)
		HandleResponse(ctx)
	}

	ProcessRequest("jane","abc123")

	/*
	你使用的键值必须满足 golang 的可比性概念，也就是运算符 == 和 != 在使用时需要返回正确的结果

	返回值必须按期，才能从多个goroutine 访问
	*/

	//下面我们自定义键类型，防止上下文的冲突
	type ctxKey int

	const(
		ctxUserID ctxKey = iota
		ctxAuthToken
	)

	UserID := func(c context.Context) string{
		return c.Value(ctxUserID).(string)
	}

	AuthToken := func(c context.Context) string {
		return c.Value(ctxAuthToken).(string)
	}

	HandleResponse2 := func(ctx context.Context) {
		fmt.Printf(
			"handling response for %v (auth:%v)\n",
			UserID(ctx),
			AuthToken(ctx),
			)
	}

	ProcessRequest2 := func(userID,authToken string){
		ctx := context.WithValue(context.Background(),ctxUserID,userID)
		ctx = context.WithValue(ctx,ctxAuthToken,authToken)
		HandleResponse2(ctx)
	}

	ProcessRequest2("make","efg321")
}