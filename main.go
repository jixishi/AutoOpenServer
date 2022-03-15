package main

import (
	"Mindustry/tools"
	"bufio"
	"fmt"
	"github.com/adlane/exec"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"golang.org/x/text/encoding/simplifiedchinese"
	"os"
	"os/signal"
	"runtime"
	"strconv"
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
var info [2]string
var Debug bool = true
var (
	kernel32    *syscall.LazyDLL  = syscall.NewLazyDLL(`kernel32.dll`)
	proc        *syscall.LazyProc = kernel32.NewProc(`SetConsoleTextAttribute`)
	CloseHandle *syscall.LazyProc = kernel32.NewProc(`CloseHandle`)

	// FontC 给字体颜色对象赋值
	FontC Color = Color{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
)

type Color struct {
	black        int // 黑色
	blue         int // 蓝色
	green        int // 绿色
	cyan         int // 青色
	red          int // 红色
	purple       int // 紫色
	yellow       int // 黄色
	light_gray   int // 淡灰色（系统默认值）
	gray         int // 灰色
	light_blue   int // 亮蓝色
	light_green  int // 亮绿色
	light_cyan   int // 亮青色
	light_red    int // 亮红色
	light_purple int // 亮紫色
	light_yellow int // 亮黄色
	white        int // 白色
}

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB = 1 << (10 * iota)
	GB = 1 << (10 * iota)
	TB = 1 << (10 * iota)
	PB = 1 << (10 * iota)
)

func GetCpu() string {
	percent, _ := cpu.Percent(time.Second, false)
	return "C:" + strconv.FormatFloat(percent[0], 'f', 1, 64) + "%"
}

func GetMem() string {
	v, _ := mem.VirtualMemory()
	return "M:" + strconv.FormatFloat(v.UsedPercent, 'f', 1, 64) + "%)"
}
func CPrint(s string, i int) {
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(i))
	fmt.Print(s)
	_, _, err := CloseHandle.Call(handle)
	if err != nil {
		return
	}
}
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
	if !Debug {
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
	info[0] = GetCpu()
	info[1] = GetMem()
	input := bufio.NewScanner(os.Stdin)
	var cmds = [...]string{".help", ".restart", ".exit"}
	var deps = [...]string{"帮助信息", "重新启动服务端", "退出启动器"}
	go func() {
		for {
			CPrint("\r"+info[0]+info[1]+"> ", FontC.light_green)
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
	Msg := ConvertByte2String(b, Charset(font))
	if runtime.GOOS == "windows" {
		if find := strings.Contains(Msg, "I"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 34, Msg)
			if find := strings.Contains(Msg, "* "); find {
				//fmt.Printf("\x1b[%dm%v\x1b[0m", 33, Msg)
				CPrint(Msg, FontC.light_purple)
				return false
			}
			CPrint(Msg, FontC.light_cyan)
			return false
		}
		if find := strings.Contains(Msg, "E"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 33, Msg)
			CPrint(Msg, FontC.red)
			return false
		}
		CPrint(Msg, FontC.white)
		return false
	} else {
		fmt.Print(Msg)
		return false
	}
}

func (*reader) OnError(b []byte) bool {
	Msg := ConvertByte2String(b, Charset(font))
	if runtime.GOOS == "windows" {
		if find := strings.Contains(Msg, " 警告"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 32, Msg)
			CPrint(Msg, FontC.purple)
			return false
		}
		if find := strings.Contains(Msg, "INFO"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 34, Msg)
			CPrint(Msg, FontC.light_cyan)
			return false
		}
		if find := strings.Contains(Msg, "信息"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 32, Msg)
			CPrint(Msg, FontC.light_blue)
			return false
		}

		if find := strings.Contains(Msg, "WARN"); find {
			//fmt.Printf("\x1b[%dm%v\x1b[0m", 35, Msg)
			CPrint(Msg, FontC.purple)
			return false
		}
		CPrint(Msg, FontC.white)
		return false
	} else {
		fmt.Print(Msg)
		return false
	}
}

func (*reader) OnTimeout() {
	go ctx.Receive(&r, 5*time.Second)
	info[0] = GetCpu()
	info[1] = GetMem()

}
