package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type BililiveRecorderElement struct {
	XMLName xml.Name `xml:"BililiveRecorder"`
	Version string   `xml:"version,attr"`
}
type BililiveRecorderRecordInfoElement struct {
	XMLName        xml.Name `xml:"BililiveRecorderRecordInfo"`
	Name           string   `xml:"name,attr"`
	Title          string   `xml:"title,attr"`
	Areanameparent string   `xml:"areanameparent,attr"`
	Areanamechild  string   `xml:"areanamechild,attr"`
	Start_time     string   `xml:"start_time,attr"`
	Roomid         int      `xml:"roomid,attr"`
	Shortid        int      `xml:"shortid,attr"`
}
type DamakuElement struct {
	XMLName xml.Name `xml:"d"`
	Attr    string   `xml:"p,attr"`
	User    string   `xml:"user,attr"`
	Message string   `xml:",chardata"`
}
type GiftElement struct {
	XMLName   xml.Name `xml:"gift"`
	Timestamp float64  `xml:"ts,attr"`
	User      string   `xml:"user,attr"`
	Uid       int      `xml:"uid,attr"`
	Giftname  string   `xml:"giftname,attr"`
	GiftCount int      `xml:"giftcount,attr"`
}
type GuardElement struct {
	XMLName   xml.Name `xml:"guard"`
	Timestamp float64  `xml:"ts,attr"`
	User      string   `xml:"user,attr"`
	Uid       int      `xml:"uid,attr"`
	Giftname  int      `xml:"level,attr"`
	GiftCount int      `xml:"count,attr"`
}

type SuperChatElement struct {
	XMLName   xml.Name `xml:"sc"`
	Timestamp float64  `xml:"ts,attr"`
	User      string   `xml:"user,attr"`
	Uid       int      `xml:"uid,attr"`
	Price     int      `xml:"price,attr"`
	Time      int      `xml:"time,attr"`
	Value     string   `xml:",chardata"`
}
type BililiveRecorderXmlStyleElement struct {
	XMLName xml.Name `xml:"BililiveRecorderXmlStyle"`
	Style   string   `xml:",innerxml"`
}
type RootElement struct {
	XMLName                    xml.Name                          `xml:"i"`
	Chatserver                 string                            `xml:"chatserver"`
	Chatid                     string                            `xml:"chatid"`
	Mission                    string                            `xml:"mission"`
	Maxlimit                   string                            `xml:"maxlimit"`
	State                      string                            `xml:"state"`
	Real_name                  string                            `xml:"real_name"`
	Source                     string                            `xml:"source"`
	BililiveRecorder           BililiveRecorderElement           `xml:"BililiveRecorder"`
	BililiveRecorderRecordInfo BililiveRecorderRecordInfoElement `xml:"BililiveRecorderRecordInfo"`
	BililiveRecorderXmlStyle   BililiveRecorderXmlStyleElement   `xml:"BililiveRecorderXmlStyle"`
	Damakus                    []DamakuElement                   `xml:"d"`
	Gifts                      []GiftElement                     `xml:"gift"`
	Guards                     []GuardElement                    `xml:"guard"`
	SuperChats                 []SuperChatElement                `xml:"sc"`
}

func main() {
	fileCount := len(os.Args) - 1
	if fileCount < 2 {
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "需要拼接%d个弹幕文件\n", fileCount)
	fmt.Fprintf(os.Stderr, "正在读取第1个弹幕文件\n")
	file, _ := os.Open(os.Args[1]) // For read access.
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	root := RootElement{}
	xml.Unmarshal(data, &root)
	fmt.Fprintf(os.Stderr, "获取第1个弹幕文件基准时钟:%s\n", root.BililiveRecorderRecordInfo.Start_time)
	StartTime, _ := time.Parse(time.RFC3339, root.BililiveRecorderRecordInfo.Start_time)
	StartTimeStamp := StartTime.UnixMilli()
	fmt.Fprintf(os.Stderr, "解析基准时钟:%d\n", StartTimeStamp)
	for i := 2; i <= fileCount; i++ {
		fmt.Fprintf(os.Stderr, "正在读取第%d个弹幕文件\n", i)
		tmpFile, _ := os.Open(os.Args[i])
		defer tmpFile.Close()
		tmpData, _ := ioutil.ReadAll(tmpFile)
		tmpRoot := RootElement{}
		xml.Unmarshal(tmpData, &tmpRoot)
		fmt.Fprintf(os.Stderr, "获取第%d个弹幕文件基准时钟:%s\n", i, tmpRoot.BililiveRecorderRecordInfo.Start_time)
		tStartTime, _ := time.Parse(time.RFC3339, tmpRoot.BililiveRecorderRecordInfo.Start_time)
		tStartTimeStamp := tStartTime.UnixMilli()
		fmt.Fprintf(os.Stderr, "解析基准时钟:%d\n", tStartTimeStamp)
		fmt.Fprintln(os.Stderr, "正在转换弹幕时间戳")
		for _, Danmaku := range tmpRoot.Damakus {
			Attributes := strings.Split(Danmaku.Attr, ",")
			tTimeStamp, _ := strconv.ParseInt(Attributes[4], 10, 64)
			tTimeOffset := float64(tTimeStamp-StartTimeStamp) / 1000
			sTimeOffset := fmt.Sprintf("%.3f", tTimeOffset)
			Attributes[0] = sTimeOffset
			Danmaku.Attr = strings.Join(Attributes, ",")
			root.Damakus = append(root.Damakus, Danmaku)
		}
		fmt.Fprintln(os.Stderr, "正在转换礼物时间戳")
		for _, Gift := range tmpRoot.Gifts {
			tTimeStamp := int64(Gift.Timestamp*1000) + tStartTimeStamp
			tTimeOffset := float64(tTimeStamp-StartTimeStamp) / 1000
			Gift.Timestamp = tTimeOffset
			root.Gifts = append(root.Gifts, Gift)
		}
		fmt.Fprintln(os.Stderr, "正在转换舰长时间戳")
		for _, Guard := range tmpRoot.Guards {
			tTimeStamp := int64(Guard.Timestamp*1000) + tStartTimeStamp
			tTimeOffset := float64(tTimeStamp-StartTimeStamp) / 1000
			Guard.Timestamp = tTimeOffset
			root.Guards = append(root.Guards, Guard)
		}
		fmt.Fprintln(os.Stderr, "正在转换SuperChat时间戳")
		for _, SuperChat := range tmpRoot.SuperChats {
			tTimeStamp := int64(SuperChat.Timestamp*1000) + tStartTimeStamp
			tTimeOffset := float64(tTimeStamp-StartTimeStamp) / 1000
			SuperChat.Timestamp = tTimeOffset
			root.SuperChats = append(root.SuperChats, SuperChat)
		}
	}
	output, _ := xml.MarshalIndent(&root, "", "  ")
	os.Stdout.Write(output)
}
