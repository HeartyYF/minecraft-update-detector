package main

import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"bufio"
	"net/http"
	"time"
	"encoding/json"
)

var (
	inverval = 5
	url = "https://launchermeta.mojang.com/mc/game/version_manifest.json"
)

func check(filename string) (bool) {
	var exist = true;
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false;
	}
	return exist;
}

func main() {
	var release, snapshot string
	//检测存在性
	if check("release.txt") {
		f, err:= os.Open("release.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		s, _, _ := br.ReadLine()
		release = string(s)
	}
	if check("snapshot.txt") {
		f, err:= os.Open("snapshot.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		s, _, _ := br.ReadLine()
		snapshot = string(s)
	}

home:
	//获取manifest
	res, err := http.Get(url)
	if err != nil {
		now := time.Now()
		fmt.Printf("%d-%02d-%02d %02d:%02d:%02d: Failed to fetch manifest. Reason: %s\r\n",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second(),err)
		time.Sleep(time.Duration(inverval)*time.Second)
		goto home
	}
	body, err := ioutil.ReadAll(res.Body)
	f, err := os.Create("version_manifest.json")
	if err != nil {
		panic(err)
	}
	io.WriteString(f, string(body))

	//json映射
	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		panic(err)
	}
	crelease := m["latest"].(map[string]interface{})["release"].(string)
	csnapshot := m["latest"].(map[string]interface{})["snapshot"].(string)
	
	now := time.Now()
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d: ",now.Year(),now.Month(),now.Day(),now.Hour(),now.Minute(),now.Second())
	fmt.Printf("Release: %s    Snapshot: %s \r\n",crelease,csnapshot)
	
	//判断是否新版本
	if crelease != release {
		fmt.Printf("New Release Found: %s\r\n",crelease)
		f, err:= os.Create("release.txt")
		if err != nil {
			panic(err)
		}
		io.WriteString(f, crelease)
		release = crelease
	}
	if csnapshot != snapshot {
		fmt.Printf("New Snapshot Found: %s\r\n",csnapshot)
		f, err:= os.Create("snapshot.txt")
		if err != nil {
			panic(err)
		}
		io.WriteString(f, csnapshot)
		snapshot = csnapshot
	}
	
	time.Sleep(time.Duration(inverval)*time.Second)
	goto home
}
