package main

import (
	"Mindustry/tools"
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
func main() {
	fmt.Println("欢迎使用BE6-CLOUD一键开服器")
	var config tools.Config
	config = tools.GetConfig()
	MVersion, MDownUrl := tools.MGetVersion(config.MindustryTagUrl)
	WVersion, WJarUrl, WZipUrl := tools.WGetVersion(config.WayZerTagUrl)
	fmt.Println("Mindustry当前版本：", config.MindustryVersion)
	fmt.Println("WayZer当前版本：", config.WayZerVersion)
	DownList := tools.NewDownloader("./")
	DownList.Concurrent = 3
	var Down bool
	if config.MindustryVersion != MVersion {
		fmt.Println("Mindustry最新版本：", MVersion, "有更新")
		config.MindustryVersion = MVersion
		DownList.AppendResource("server.jar", MDownUrl)
		Down = true
	}
	if config.WayZerVersion != WVersion {
		fmt.Println("WayZer最新版本：", WVersion, "有更新")
		config.WayZerVersion = WVersion

		DownList.AppendResource("WayZer.jar", WJarUrl)
		DownList.AppendResource("WayZer.zip", WZipUrl)
		Down = true
	}
	if Down {
		DownList.Start()
		tools.SaveConfig(config)
	}
	//f, err := exec.LookPath("java")
	//if err != nil {
	//fmt.Println(err)
	//}
	//fmt.Println(f)
	fmt.Println("Server已经开启输入stop停止服务器!")

	var cmd string
	var font string
	if runtime.GOOS == "windows" {
		cmd = "java"
		font = "GB18030"
	} else {
		cmd = "java"
		font = "UTF-8"
	}
	server := exec.Command(cmd, "-jar", "./server.jar")
	serverReader, err := server.StdoutPipe()
	serverWriter, err := server.StdinPipe()
	if err != nil {
		fmt.Printf("Error: Can not obtain the stdin pipe for command: %s\n", err)
		return
	}
	var sin = bufio.NewScanner(serverReader)
	if err != nil {
		fmt.Printf("create cmd stdoutpipe failed,error:%s\n", err)
		os.Exit(1)
	}
	err = server.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for sin.Scan() {
			cmdRe := ConvertByte2String(sin.Bytes(), Charset(font))
			fmt.Println(cmdRe)
		}
	}()
	input := bufio.NewScanner(os.Stdin)
	go func() {
		for {
			fmt.Print("\r> ")
			input.Scan()
			if strings.Compare(strings.TrimSpace(input.Text()), "") == 0 {
				continue
			}
			if strings.Compare(strings.TrimSpace(input.Text()), ".exit") == 0 {
				os.Exit(0)
			}
			_, err := serverWriter.Write(input.Bytes())
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()
	// Wait功能将等待，直到进程结束
	err = server.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 在进程终止后，*os.ProcessState 包含有关进程运行的简单信息
	fmt.Printf("PID: %d\n", server.ProcessState.Pid())
	fmt.Printf("程序运行时间: %dms\n",
		server.ProcessState.SystemTime()/time.Microsecond)
	fmt.Printf("成功退出: %t\n", server.ProcessState.Success())
}
