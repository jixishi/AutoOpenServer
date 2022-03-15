package main

import (
	"Mindustry/tools"
	"bufio"
	"fmt"
	"github.com/adlane/exec"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

var ctx exec.ProcessContext
var r reader
var cmd string
var font string
var arg [4]string

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
func init() {
	fmt.Println("欢迎使用BE6-CLOUD一键开服器")
	var config tools.Config
	err := os.Chdir("./")
	err = os.MkdirAll("./config/mods", 0777)
	err = os.MkdirAll("./config/scripts", 0777)
	err = os.MkdirAll("./config/scripts/cache", 0777)
	config = tools.GetConfig()
	MVersion, MDownUrl := tools.MGetVersion(config.MindustryTagUrl)
	WVersion, WJarUrl, WZipUrl, WCaAUrl := tools.WGetVersion(config.WayZerTagUrl)
	fmt.Println("Mindustry当前版本：", config.MindustryVersion)
	fmt.Println("WayZer当前版本：", config.WayZerVersion)
	DownList := tools.NewDownloader("./")
	DownList.Concurrent = 3
	var Down bool
	var Mode = false
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
		DownList.AppendResource("Cache.zip", WCaAUrl)

		Down = true
		Mode = true
	}
	if Down {
		err = DownList.Start()
		tools.SaveConfig(config)
		if Mode {
			err = os.Rename("./WayZer.jar", "./config/mods/WayZer.jar")
			err = tools.DeCompressZip("./WayZer.zip", "./config/scripts")
			err = tools.DeCompressZip("./Cache.zip", "./config/scripts/cache")

		}
	}
	if err != nil {
		return
	}
}
func main() {

	//f, err := exec.LookPath("java")
	//if err != nil {
	//fmt.Println(err)
	//}
	//fmt.Println(f)
	time.Sleep(6 * time.Microsecond)
	fmt.Println("Server已经开启输入stop停止服务器!")

	if runtime.GOOS == "windows" {
		cmd = "java"
		arg = [4]string{"--add-opens", "java.base/java.net=ALL-UNNAMED", "--add-opens", "java.base/java.security=ALL-UNNAMED"}
		font = "GB18030"
	} else {
		cmd = "java"
		font = "UTF-8"
	}
	ctx = exec.InteractiveExec(cmd, arg[0], arg[1], arg[2], arg[3], "-jar", "./server.jar")

	go ctx.Receive(&r, 10*time.Second)
	input := bufio.NewScanner(os.Stdin)
	var cmds = [...]string{".help", ".restart", ".exit"}
	var deps = [...]string{"帮助信息", "重新启动服务端", "退出启动器"}
	go func() {
		for {
			fmt.Print("\r> ")
			input.Scan()
			if strings.Compare(strings.TrimSpace(input.Text()), "") == 0 {
				continue
			}
			if strings.Compare(strings.TrimSpace(input.Text()), ".help") == 0 {
				var page = 0
				fmt.Printf(">-------Help(%v)-------<\n", page)
				for i := 0; i < len(cmds); i++ {
					fmt.Printf("%v --%v\n", cmds[i], deps[i])
				}
			}
			if strings.Compare(strings.TrimSpace(input.Text()), ".restart") == 0 {
				ctx.Cancel()
				ctx.Stop()
				ctx := exec.InteractiveExec(cmd, arg[0], arg[1], arg[2], arg[3], "-jar", "./server.jar")
				r := reader{}
				go ctx.Receive(&r, 10*time.Second)
			}
			if strings.Compare(strings.TrimSpace(input.Text()), ".exit") == 0 {
				ctx.Cancel()
				ctx.Stop()
				os.Exit(0)
			}
			var m = input.Text() + "\n"
			err := ctx.Send(m)
			if err != nil {
				println(err)
			}
		}
	}()
	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}

type reader struct {
}

func (*reader) OnData(b []byte) bool {
	fmt.Print(ConvertByte2String(b, Charset(font)))
	return false
}

func (*reader) OnError(b []byte) bool {
	fmt.Print(ConvertByte2String(b, Charset(font)))
	return false
}

func (*reader) OnTimeout() {
	go ctx.Receive(&r, 10*time.Second)
}
