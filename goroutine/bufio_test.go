package main

import (
	"testing"
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"bufio"
)

/*
测试缓冲和未缓冲的写入速度
缓冲写入比未缓冲写入更快。
这是因为在bufio.Writer中，写入在内部的数据，排队到缓冲区中，知道已经积累了足够的块为止，然后块被写出。
这个过程通常称为分块。
*/

//未缓冲
func BenchmarkUnbufferedWrited(b * testing.B){
	performWrite(b,tmpFileOrFatal())
}

//缓冲
func BenchmarkBufferedWrite(b *testing.B){
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b,bufio.NewWriter(bufferedFile))
}

func tmpFileOrFatal() *os.File{
	file,err := ioutil.TempFile("","tmp")
	if err != nil{
		fmt.Printf("error: %v\n",err)
	}
	return file
}

func performWrite(b *testing.B,writer io.Writer){
	repeat := func(
		done <-chan interface{},
		vals ...interface{},
		) <-chan interface {}{
		repeatStream := make(chan interface{})
		go func() {
			defer close(repeatStream)
			for{
				for val := range vals{
					select{
					case <- done:
						return
					case repeatStream<-uint8(val): //将val 的byte类型转换成 uint8 类型
					}
				}
			}
		}()
		return repeatStream
	}
	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
		)<-chan interface{}{
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i:=num;i>0||i!=-1;{
				if i != -1{
					i--
				}
				select {
				case <-done:
					return
				case takeStream<-<-valueStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done,repeat(done,byte(0)),b.N){
		writer.Write([]byte{bt.(byte)})
	}
}

