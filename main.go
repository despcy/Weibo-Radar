package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"

	"github.com/buger/jsonparser"
)

type Card struct {
	Id         string
	ScreenName string
	Desc1      string
	Desc2      string
	Info       UserInfo
}

type UserInfo struct {
	ScreenName   string
	Sex          string
	Location     string
	Birthday     string
	SexOri       string
	SingleStatus string
	Intro        string
	Labels       string
	Study        string
	Work         string
}

var latitude = "38.12345"
var lontitude = "115.12345"
var Cookie = ""
var count = 0
var writer *csv.Writer

func WriteToCSV(data Card) {
	writer.Write([]string{data.Id, data.ScreenName, data.Desc1, data.Desc2, data.Info.ScreenName, data.Info.Sex, data.Info.Location, data.Info.Birthday, data.Info.SexOri, data.Info.SingleStatus, data.Info.Intro, data.Info.Labels, data.Info.Study, data.Info.Work})
}

func main() {

	file, _ := os.Create("result.csv")
	defer file.Close()
	writer = csv.NewWriter(file)

	defer writer.Flush()
	writer.Write([]string{"id", "screenName", "Desc1", "Desc2", "ScreenName", "Sex", "Location", "Birthday", "SexOri", "SingleStatus", "Intro", "Labels", "Study", "Work"})
	// Instantiate default collector
	c := colly.NewCollector(

		//colly.AllowedDomains("weibo.cn"),
		colly.Debugger(&debug.LogDebugger{}),
	)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {

	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)

	})
	c.OnHTML("body", func(e *colly.HTMLElement) {

		card := e.Request.Ctx.GetAny("userData").(Card)

		infostr := e.ChildText("div.c:nth-child(7)")

		//	fmt.Println(regexp.MustCompile("标签:.*").FindString(infostr))

		userInfo := UserInfo{
			ScreenName:   regexp.MustCompile("昵称\\:.*?\\:").FindString(infostr),
			Sex:          regexp.MustCompile("性别\\:.*?\\:").FindString(infostr),
			Location:     regexp.MustCompile("地区\\:.*?\\:").FindString(infostr),
			Birthday:     regexp.MustCompile("生日\\:.*?\\:").FindString(infostr),
			SexOri:       regexp.MustCompile("性取向\\:.*?\\:").FindString(infostr),
			SingleStatus: regexp.MustCompile("感情状况\\:.*?\\:").FindString(infostr),
			Intro:        regexp.MustCompile("简介:.*?:").FindString(infostr),
			Labels:       regexp.MustCompile("标签:.*").FindString(infostr),
			Study:        e.ChildText("div.c:nth-child(9)"),
			Work:         e.ChildText("div.c:nth-child(11)"),
		}
		card.Info = userInfo
		fmt.Println(card)
		WriteToCSV(card)

	})

	c.OnResponse(func(r *colly.Response) {

		jsonparser.ArrayEach(r.Body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			userID, _, _, err := jsonparser.Get(value, "user", "id")
			//	fmt.Println(string(userID))
			screenName, _, _, err := jsonparser.Get(value, "user", "screen_name")
			//	fmt.Println(string(screenName))
			desc1, _, _, err := jsonparser.Get(value, "desc1")
			//	fmt.Println(string(desc1))
			desc2, _, _, err := jsonparser.Get(value, "desc2")
			//fmt.Println(string(desc2))
			//fmt.Print("----------")
			card := Card{
				Id:         "https://www.weibo.com/" + string(userID),
				ScreenName: string(screenName),
				Desc1:      string(desc1),
				Desc2:      string(desc2),
			}
			r.Ctx.Put("userData", card)

			r.Headers.Set("cookie", Cookie)

			c.Request("GET", "https://weibo.cn/"+string(userID)+"/info", nil, r.Ctx, *r.Headers)
			count++

			fmt.Println(count)

			time.Sleep(2000 * time.Millisecond)
		}, "cards", "[0]", "card_group")

	})

	// Start scraping
	for i := 1; i < 200; i++ {
		time.Sleep(1000 * time.Millisecond)

		c.Visit("https://api.weibo.cn/2/guest/cardlist?lat=" + latitude + "&lon=" + lontitude + "&page=" + strconv.Itoa(i) + "&count=20&containerid=2317120015_111_1")
	}

}
