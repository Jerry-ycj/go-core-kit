package cmdkit

import (
	"bufio"
	"errors"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/spf13/cast"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type RunParams struct {
	Timeout int  `comment:"超时时间s"`
	Async   bool `comment:"异步处理返回值"`
}

func Run(command string, params ...RunParams) (string, error) {
	// todo 判断系统环境
	var param RunParams
	if len(params) == 0 {
		param = RunParams{}
	} else {
		param = params[0]
	}

	cmd := exec.Command("/bin/sh", "-c", command)
	// 程序退出时Kill子进程
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if !param.Async {
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			return "", err
		}
		if err = cmd.Start(); err != nil {
			return "", err
		}
		if param.Timeout > 0 {
			to := make(chan map[string]interface{})
			go func() {
				ret0, err2 := getRet(stdout, stderr, cmd)
				to <- map[string]interface{}{"ret": ret0, "err": err2}
			}()
			select {
			case <-time.After(time.Duration(param.Timeout) * time.Second):
				return "", errors.New("timeout")
			case m := <-to:
				ret := m["ret"].(string)
				var err error
				if m["err"] != nil {
					err = m["err"].(error)
				}
				return cast.ToString(ret), err
			}
		} else {
			ret, err := getRet(stdout, stderr, cmd)
			return ret, err
		}
	} else {
		if err := cmd.Start(); err != nil {
			return "", err
		}
	}
	return "", nil
}

func getRet(stdout io.ReadCloser, stderr io.ReadCloser, cmd *exec.Cmd) (string, error) {
	ret := ""
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		ret += line
	}
	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return ret, err
	}
	if len(bytesErr) != 0 {
		return ret, errors.New(string(bytesErr))
	}
	if err = cmd.Wait(); err != nil {
		return ret, err
	}
	return ret, nil
}

func LinuxCmd(name string, args ...string) error {
	command := exec.Command(name, args...)
	return command.Run()
}

func WinCmd(arg ...string) error {
	args := make([]string, 0, 3)
	args = append(args, "/C")
	args = append(args, arg...)
	command := &exec.Cmd{
		Path: "cmd",
		Args: args,
	}
	if filepath.Base("cmd") == "cmd" {
		if lp, err := exec.LookPath("cmd"); err != nil {
			logkit.Error("filePathErr")
		} else {
			command.Path = lp
		}
	}
	return command.Run()
}

func GetInput() string {
	//使用os.Stdin开启输入流
	//函数原型 func NewReader(rd io.Reader) *Reader
	//NewReader创建一个具有默认大小缓冲、从r读取的*Reader 结构见官方文档
	in := bufio.NewReader(os.Stdin)
	//in.ReadLine函数具有三个返回值 []byte bool error
	//分别为读取到的信息 是否数据太长导致缓冲区溢出 是否读取失败
	str, _, err := in.ReadLine()
	if err != nil {
		logkit.Error(err.Error())
		return ""
	}
	return string(str)
}
