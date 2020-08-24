package util

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var getaway string=""
var getawayStart bool = false
var w sync.WaitGroup
func linuxGeteway() () {
	sysType := runtime.GOOS
	if "linux" == sysType && ""==getaway {
		path, _ := os.Getwd()
		ch := exec.Command("sh", "-c", "chmod 777 "+path+"/module/getaway")
		o, err := ch.Output()
		if err != nil {
			fmt.Println(o)
			fmt.Println(err)
		}
		cmd := exec.Command(path+"/module/getaway")
		cmd.Wait()
		// 执行命令，并返回结果
		output, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		// 因为结果是字节数组，需要转换成string
		//fmt.Println("网关:"+string(output))
		getaway = string(output)
	}


	fmt.Println("网关=====================:"+string(getaway))

}

func whilelinuxGeteway(){
	sysType := runtime.GOOS
	if getawayStart==false && "linux" == sysType {
		path, _ := os.Getwd()
		ch := exec.Command("sh", "-c", "chmod 777 "+path+"/module/getaway")
		o, err := ch.Output()
		if err != nil {
			fmt.Println(o)
			fmt.Println(err)
		}
		linuxGeteway()
	}
	if getawayStart==true || "linux" != sysType{
		return
	}
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	// 时间打印
	go func(c context.Context, cel context.CancelFunc) {
		getawayStart=true
		sum := 0
		for {
			sum ++
			if sum > 10{
				break
			}else{
				linuxGeteway()
				fmt.Println("线程循环:"+getaway)
				time.Sleep(5 * 1e9)
			}
		}
	}(ctx, cancel)

}
