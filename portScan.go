package main

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

func checkPort(ip net.IP, port int,wg *sync.WaitGroup) {
	tcpAddr := net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	ch := make(chan bool)
	timeout := make(chan bool)
	go func(){
		time.Sleep(3*time.Second)
		timeout <- true
	}()
	go func(){
		conn, err := net.DialTCP("tcp", nil, &tcpAddr)
		ch <- true
		if err == nil {
			fmt.Printf("ip: %v port: %v \n",ip,port)
			defer func(){
				if conn != nil{
					e := conn.Close()
					if e !=nil{
						fmt.Println(e)
					}
				}
			}()
		}
	}()
	select {
	case <- timeout:
		wg.Done()
	case <- ch:
		wg.Done()
	}
}
func checkIp(ip string) bool{
	if net.ParseIP(ip) == nil {
		fmt.Println("非法ip地址")
		return false
	}else {
		return true
	}
}
func main() {

	startTime := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(65534)
	ip := os.Args[1]
	if checkIp(ip) {
		for port := 1; port <= 65534;port++  {
			go checkPort(net.ParseIP(ip),port,&wg)
		}
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("执行时间%v",endTime.Sub(startTime))
}