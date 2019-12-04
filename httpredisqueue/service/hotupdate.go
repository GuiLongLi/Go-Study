package service

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

/************************** 热重启 ***************************/

var (
	listener net.Listener = nil

	graceful =  flag.Bool("graceful", false, "listen on fd open 3 (internal use only)")
)


//监听服务器
func Listenserver(server *http.Server){
	var err error

	//解析参数
	flag.Parse()

	//设置监听的对象(新建或已存在的socket描述符)
	if *graceful {
		//子进程监听父进程传递的 socket描述符
		log.Println("listening on the existing file descriptor 3")
		//子进程的 0 1 2 是预留给 标准输入 标准输出 错误输出
		//因此传递的socket 描述符应该放在子进程的 3
		f := os.NewFile(3,"")
		listener,err = net.FileListener(f)
		log.Printf( "graceful-reborn  %v %v  %#v \n", f.Fd(), f.Name(), listener)
	}else{
		//父进程监听新建的 socket 描述符
		log.Println("listening on a new file descriptor")
		listener,err = net.Listen("tcp",server.Addr)
		log.Printf("Actual pid is %d\n", syscall.Getpid())
	}
	if err != nil{
		log.Fatalf("listener error: %v\n",err)
	}
	go func(){
		err = server.Serve(listener)
		log.Printf("server.Serve err: %v\n",err)
		tcp,_ := listener.(*net.TCPListener)
		fd,_ := tcp.File()
		log.Printf( "first-boot  %v %v %#v \n ", fd.Fd(),fd.Name(), listener)
	}()
	//监听信号
	handleSignal(server)
	log.Println("signal end")
}

//处理信号
func handleSignal(server *http.Server){
	//把信号 赋值给 通道
	ch := make(chan os.Signal, 1)
	//监听信号
	signal.Notify(ch, syscall.SIGINT,syscall.SIGTERM,syscall.SIGUSR2)
	//阻塞主进程， 不停的监听系统信号
	for{
		//通道 赋值给 sig
		sig := <-ch
		log.Printf("signal receive: %v\n", sig)
		ctx,_ := context.WithTimeout(context.Background(),20*time.Second)
		switch sig{
		case syscall.SIGINT,syscall.SIGTERM:  //终止进程执行
			log.Println("shutdown")
			signal.Stop(ch)       //停止通道
			server.Shutdown(ctx)  //关闭服务器窗口
			log.Println("graceful shutdown")
			return
		case syscall.SIGUSR2:  //进程热重启
			log.Println("reload")
			err := reload()  //执行热重启
			if err != nil{
				log.Fatalf("listener error: %v\n",err)
			}
			//server.Shutdown(ctx)
			log.Println("graceful reload")
			return
		}
	}
}

//热重启
func reload() error{
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}
	//获取socket描述符
	currentFD, err := tl.File()
	if err != nil {
		return err
	}
	//设置传递给子进程的参数(包含 socket描述符)
	args := []string{"-graceful"}
	//args = append(args, "-continue")
	cmd := exec.Command(os.Args[0],args...)
	log.Println("os.Args[0]",os.Args[0])
	log.Printf("%+v\n",os.Args)
	log.Printf("args %+v \n",args)
	cmd.Stdout = os.Stdout  //标准输出
	cmd.Stderr = os.Stderr  //错误输出
	cmd.ExtraFiles = []*os.File{currentFD} //文件描述符

	err = cmd.Start()
	log.Printf("forked new pid %v: \n",cmd.Process.Pid)
	if err != nil{
		return err
	}
	return nil
}
/*
我们在父进程执行 cmd.ExtraFiles = []*os.File{f} 来传递 socket 描述符给子进程，子进程通过执行 f := os.NewFile(3, "") 来获取该描述符。值得注意的是，子进程的 0 、1 和 2 分别预留给标准输入、标准输出和错误输出，所以父进程传递的 socket 描述符在子进程的顺序是从 3 开始。
*/
