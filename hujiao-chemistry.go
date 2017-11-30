package main

import (
	"time"
  "learn/app/library/database"
  "fmt"
  "log"
  "goquery-book/model"
  "github.com/PuerkitoBio/goquery"
  "github.com/djimenez/iconv-go"
  "strconv"
  "strings"
  "net/http"
  "encoding/json"
  "os"
  "io/ioutil"
  "io"
  "bytes"
  "math/rand"
)

//语文 英语 物理 历史
//北师大版

var startPage = "http://www.czhxzx.com/hxzx/dzkb/hjbdzkb/0948203156.html"
var host = "http://www.czhxzx.com"
func BookScrape() {
	
	  for j :=0; j < 100; j++ {
				fmt.Println("开始采集页面map：", j)
				fmt.Println(startPage)
				doc, err := goquery.NewDocument(startPage)
				if err != nil {
				  log.Fatal(err)
				}
				
				title := doc.Find("div.maincontent").Find("div.showCon").Find("h1").Eq(0).Text()
				if !strings.Contains(title, "2013年新沪教版化学九年级下册") {
					fmt.Println("2013年新沪教版化学九年级下册----------采集完毕")
					break
				}
				//获取下一页
				next := doc.Find("div.maincontent").Find("div.other0").Find("ul li").Last().Find("a").First()
				nextHref, _ := next.Attr("href")
				startPage = host + nextHref

				knowledgeSlice := strings.Split(title, "《")
				knowledgeSlice2 := strings.Split(knowledgeSlice[1], "》")
				knowledgeName := knowledgeSlice2[0]

				volume, gradeName := collectText(title)
				imgJson := collectImage(doc)
				 //写入数据库
				 ebook := model.MaterialEbookModel{}
				 ebook.CourseName = "化学"
				 ebook.MaterialName = "人教版"
				 ebook.Volume = volume
				 ebook.GradeName = gradeName
				 ebook.KnowledgeName = knowledgeName
				 ebook.QueryPage = j
				 ebook.ImgLink = string(imgJson)
			 
				 db := database.DB()
				 db.Create(&ebook)
				 time.Sleep(3)
	  }

  fmt.Println("采集完毕---------------------------")
}

func main() {
  BookScrape()
}

func utfString(str string) string {
  word, err:= iconv.ConvertString(str,  "gbk","utf-8")
  if err != nil {
    word = ""
    return word
  }
  return word
}

func getImage(src string) string {

	name := strconv.Itoa(int(time.Now().Unix()))
	randNumber := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := strconv.Itoa(randNumber.Intn(10000))
	imgName := name + "-" + randNum + ".jpg"
	str := "null.jpg"
    out, err := os.Create("images/" + imgName)
    defer out.Close()
    resp, err := http.Get(src)
    if err != nil {
      fmt.Println("获取图片地址错误", err)
      return str
    }
    defer resp.Body.Close()
    byteString, _ := ioutil.ReadAll(resp.Body)
    io.Copy(out, bytes.NewReader(byteString))
    return imgName
}

func collectText(str string) (volume int8, gradeName string) {
	// <a>苏教版一年级语文上册——目录</a> 或者 <a>苏教版二年级上册语文：识字4</a> 或者 <a>小学六年级苏教版语文下册 生字表</a>
	volume = 0
	
	// 上下册
	if strings.Contains(str, "上册") {
		volume = int8(0)
	} else {
		volume = int8(1)
	}
	// 年级
	if strings.Contains(str, "七年级") {
		gradeName = "七年级"
	} else if strings.Contains(str, "八年级") {
		gradeName = "八年级"
	}else {
		gradeName = "九年级"
	}
	return int8(volume), gradeName
}

func collectImage(doc *goquery.Document) []byte {
	// 抓取图片
	imgLink := []string{}
	// 第一页
	doc.Find("div.maincontent").Find("div.showCon").Find("div.article img").Each(func(index int, s *goquery.Selection){
		imgsrc, ok := s.Attr("src")
		if ok != true {
		fmt.Println("img-采集失败")
		} else {
		imgLink = append(imgLink, getImage(host + imgsrc))
		} 
	})
	//判断是否为空
	// if len(imgLink)  == 0 {
	// 	img := doc.Find("div.min_right").Find("div.home_c").Children().Eq(3).Find("img").First();
	// 	imgsrc, ok := img.Attr("src")
	// 	if ok != true {
	// 	fmt.Println("img-src-err2", err)
	// 	} else {
	// 	imgLink = append(imgLink, getImage(host + imgsrc))
	// 	} 
	// }
	imgJson, _ := json.Marshal(imgLink)
	return imgJson
}