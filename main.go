package main

import (
	"bytes"
	"forward/gocron"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func RecoverFromPanic() {
	if r := recover(); r != nil {
		log.Println("Recovered from panic: ", r)
	}
}
func getWireGuardStat() string {
	defer RecoverFromPanic()
	var out bytes.Buffer
	cmd := exec.Command("cmd.exe", "/C", "sc query WireGuardTunnel$w", "qqq")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return ""
	} else {
		return out.String()
	}
}
func run() {
	str := getWireGuardStat()
	//log.Println(str)
	log.Println(strings.Index(strings.ToUpper(str), "STOPPED"))
	if strings.Index(strings.ToUpper(str), "STOPPED") > 0 {
		log.Println("stopped,send run cmd to wireGuard service...")
		var out bytes.Buffer
		cmd := exec.Command("cmd.exe", "/C", "sc start WireGuardTunnel$w", "qqq")
		cmd.Stdout = &out
		_ = cmd.Run()
	} else {
		log.Println("already running...")
		gocron.Clear()
		gocron.Every(5).Minutes().Do(runWatch)
		//gocron.Every(3).Seconds().Do(runWatch)
	}
}
func runWatch() {
	str := getWireGuardStat()
	//log.Println(str)
	log.Println("Watching..")
	if strings.Index(strings.ToUpper(str), "STOPPED") > 0 {
		log.Println("runWatch: stopped,send run cmd to wireGuard service...")
		var out bytes.Buffer
		cmd := exec.Command("cmd.exe", "/C", "sc start WireGuardTunnel$w", "qqq")
		cmd.Stdout = &out
		_ = cmd.Run()
	} else {
		log.Println("already running...")
	}
}
func main() {
	f, err := os.OpenFile("wireGuardWatch.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return
	}
	defer func() {
		f.Close()
	}()

	// 组合一下即可，os.Stdout代表标准输出|
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	gocron.Every(15).Second().Do(run)
	<-gocron.Start()
}
