package main

import "fmt"

func main() {
	var divisor = 100
	dividend := 10

	result,error := divisionFunc(divisor,dividend)
	fmt.Printf("除数%v 被除数%v 的结果是%v 错误信息%v\n",divisor,dividend,result,error )

	//把 被除数修改为 0
	dividend = 0
	result,error = divisionFunc(divisor,dividend)
	fmt.Printf("除数%v 被除数%v 的结果是%v 错误信息%v\n",divisor,dividend,result,error )
}

//error类型是一个接口类型，这是它的定义：
/*
type error interface {
   Error() string
}
*/
//定义一个 除法 结构体
type division struct {
	divisor int    //除数
	dividend int   //被除数
}

//声明除法错误函数
func (division division) divisionError() string{
	str := `
    停止运行，被除数不能为0
    除数: %d
    被除数: 0
`
	return fmt.Sprintf(str,division.divisor)
}

//处理除法函数
func divisionFunc (divisor int,dividend int) (result int,error string){
	if(dividend == 0){ //当被除数是0 时，抛出错误信息
		errorData := division{
			divisor:divisor,
			dividend:dividend,
		}
		error = errorData.divisionError()
		return 0,error
	}else{
		return divisor / dividend,""
	}
}