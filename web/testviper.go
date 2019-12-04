package main

import (
	"fmt"
	"web/config"
)

func main() {
	vipConfig,error := config.Init()                                  //vipConfig是配置
	fmt.Printf("config.init error是%v\n", error)
	//fmt.Printf("config.init vipConfig是%v\n",vipConfig,)
	database := config.GetDatabaseInfo()
	fmt.Printf("直接获取common[database]配置是%v\n", database)
	fmt.Printf("直接获取common[database][host]配置是%v\n", database["host"])

	//因为我们不知道 vipConfig 的下级是什么类型的数据，所以这里使用了interface{}
	//因此所有的类型、任意动态的内容都可以解析成 interface{}。
	for key,val := range vipConfig.(map[string]interface{}){        //循环接口类型，获取配置信息
		fmt.Printf("vipConfig 的key是%v val是%v\n",key,val )

		switch val.(type) {                                          //判断val的类型
		case map[string]interface{}:                                //如果是 interface接口类型
			for ke,va := range val.(map[string]interface{}){        //循环接口类型，获取配置信息
				fmt.Printf("vipConfig 的ke是%v va是%v\n",ke,va )

				switch va.(type) {                                    //判断va的类型
				case map[string]interface{}:                         //如果是 interface接口类型
					for k,v := range va.(map[string]interface{}){   //循环接口类型，获取配置信息
						fmt.Printf("vipConfig 的k是%v v是%v\n",k,v )
					}
				}
			}
		}
	}

	//viper可以获取服务器的环境变量
	GO111MODULE := config.GetEnvInfo("GO111MODULE")
	fmt.Printf("GO111MODULE的值是%v\n",GO111MODULE)

}

