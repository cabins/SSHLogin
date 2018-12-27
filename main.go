package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	var csvFile string
	flag.StringVar(&csvFile, "f", "demo.csv", "指定csv文件名称")
	flag.Parse()

	// 读取csv文件
	// CSV文件的每一行内容格式为：ip,port,username,password,cmds
	// 其中cmds用分号隔开，例如ls -l;whoami;pwd;ifconfig
	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal("读取CSV文件出错： " + err.Error())
	}

	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal("解析CSV文件出错： " + err.Error())
	}

	for _, line := range lines {
		ip := line[0]
		port := line[1]
		username := line[2]
		password := line[3]
		cmds := strings.Split(line[4], ";")

		SSHRun(ip, port, username, password, cmds)
	}

	fmt.Println("Enter the <Return / Enter> key to exit...")
	fmt.Scanln()
}

func SSHRun(ip, port, username, password string, cmds []string) {
	fmt.Println("Running on ip " + ip + "...")

	clientConfig := ssh.ClientConfig{
		User:username,
		Auth:[]ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", ip+":" +port, &clientConfig)
	if err != nil {
		log.Fatal(ip + ":登陆服务器错误：" + err.Error())
	}

	defer client.Close()

	for _, cmd := range cmds {
		session, err := client.NewSession()
		if err != nil {
			log.Fatal(ip + ":创建session错误： " + err.Error())
		}
		defer session.Close()

		output, err := session.Output(cmd)
		if err != nil {
			log.Println(ip + ":" + cmd + ":运行命令出错： " + err.Error())
			WriteFile("error", ip, []byte("运行命令出错： "+err.Error()))
			break
		}

		WriteFile("output", ip, output)

		time.Sleep(time.Second * 1)
	}
}

func WriteFile(fileType string, ip string, data []byte) {
	// 之所以替换掉其中的冒号，是因为这个符号在Windows系统中被当做盘符分隔符
	t := strings.Replace(time.Now().Format(time.RFC3339), ":", "", -1)

	filename := fileType + "_" + ip + "_" + t + ".log"
	err := ioutil.WriteFile(filename, data, os.ModeAppend)
	if err != nil {
		log.Println("写入文件错误：" + filename + err.Error())
	}
}