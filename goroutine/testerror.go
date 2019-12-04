package main

import (
	"os"
	"log"
	"fmt"
	"os/exec"
	"runtime/debug"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	err := runJob("1")
	if err != nil{
		msg := "there was an unexpected issue; please report this as a bug."
		if _, ok := err.(IntermediateErr);ok{
			msg = err.Error()
		}
		handleError(1,err,msg)
	}

	/*
	·if _, ok := err.(IntermediateErr);ok{
	在这里，我们检查一下异常是否是预期的类型。
	如果是，那么可以确定这是一个结构完整的异常，我们只要简单的将其中的消息传递给用户即可

	·handleError(1,err,msg)
	在这一行，我们将日志和异常消息与一个ID 绑定在一起。
	我们可以使用一个自增ID ,或者用GUID 来保证ID 的唯一性

	*/
}


//明确异常发生的原因，发生的时间，发生的位置，对用户友好的信息，告诉用户如何获取更多信息

//异常类型
type MyError struct {
	Inner error
	Message string
	StrackTrace string
	Misc map[string]interface{}
}

func(err MyError) Error() string {
	return err.Message
}
func wrapError(
	err error,
	messagef string,
	msgArgs ...interface{},
	) MyError{
	return MyError{
		Inner: err,
		Message: fmt.Sprintf(messagef,msgArgs...),
		StrackTrace: string(debug.Stack()),
		Misc: make(map[string]interface{}),
	}
	/*
	·Inner: err,
	在这里，我们存储了我们正在包装的异常。
	通常我们会希望能够找到最底层的异常，以便在需要时可以调查发生的异常

	·StrackTrace: string(debug.Stack()),
	这行代码在创建异常时记录堆栈的轨迹。
	过于复杂的错误类型经过wrapError 包装后可能会省略一些栈帧

	·Misc: make(map[string]interface{}),
	在这里，我们创建一个可以存储各种杂项的变量。
	我们可以将并发ID ,堆栈轨迹的hash 或可能有助于诊断异常的其他上下文信息存储在这里
	*/

}


//底层模块
type LowLevelErr struct {
	error
}
func isGloballyExec(path string) (bool,error){
	info,err := os.Stat(path)
	if err != nil{
		return false,LowLevelErr{(wrapError(err,err.Error()))}
	}
	return info.Mode().Perm()&0100 == 0100,nil

	/*
	·return false,LowLevelErr{(wrapError(err,err.Error()))}
	在这里，我们用自定义的异常来调用os.Stat 中的原始异常。
	在这种情况下，我们可以用这个异常传递信息，而不用对它做任何修饰
	*/
}

//中间模块
type IntermediateErr struct {
	error
}
func runJob(id string) error{
	const jobBinPath = "/bad/job/binary"
	isExecutable,err := isGloballyExec(jobBinPath)

	if err != nil{
		return IntermediateErr{wrapError(
			err,
			"cannot run job %q: requisite binaries not available",
			id,
			)}
	} else if isExecutable == false{
		return wrapError(
			nil,
			"job binary is not executable",
			id,
			)
	}
	return exec.Command(jobBinPath, "--id="+id).Run()

	/*
	·return exec.Command(jobBinPath, "--id="+id).Run()
	这里我们传递来自底层模块的异常。
	因为我们的体系结构决定，我们需要考虑从其他模块传递来的错误，而不是将它们用我们自己的错误类型包装，这里会存在一些问题，后面会提到

	·return IntermediateErr{wrapError(
	在这里，我们使用精心设计的异常信息。
	在这种情况下，我们想隐藏异常的底层细节，因为我们觉得这对我们模块的调用者来说并不重要
	*/
}

func handleError(key int,err error,message string){
	log.SetPrefix(fmt.Sprintf("[logID: %v]: ",key))
	log.Printf("%#v\n",err)
	fmt.Printf("[%v] %v\n",key,message)

	/*
	在这里，我们记录下异常的所以内容，以备有人需要深入了解发生的事情。
	*/
}