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
)

//英语
//苏教版
var start = 0
var listHost = "http://www.aoshu.com/zlk/dzkb/yy/bsdb/"
func BookScrape() {
	// 抓取小学所有页面放入 切片
	primaryPage := make([]map[string]string, 0 , 20)
	resource := listHost
    doc, err := goquery.NewDocument(resource)
    if err != nil {
      log.Fatal(err)
	}
	doc.Find("div.m-grade3 ul").Find("li a").Each(func(index int, s *goquery.Selection){
		eleText := utfString(s.Text())
		eleLink, _ := s.Attr("href")
		result := make(map[string]string , 2)
		result["text"] = eleText
		result["link"] = eleLink
		primaryPage = append(primaryPage, result)
	  })

	  //开始采集
	  fmt.Println("一共收集到采集页面map长度:", len(primaryPage))
	  restPage := primaryPage[start:len(primaryPage)]
	  for index, value := range restPage {
				start = index
				volume, gradeName := collectText(value["text"])
				if (gradeName == "") {
					continue
				}
				fmt.Println("开始采集页面map：", start, "，目录：", value["text"])
				imgJson, knowledgeName := collectImage(value["link"])
				 //写入数据库
				 ebook := model.MaterialEbookModel{}
				 ebook.CourseName = "英语"
				 ebook.MaterialName = "北师大"
				 ebook.Volume = volume
				 ebook.GradeName = gradeName
				 ebook.KnowledgeName = knowledgeName
				 ebook.QueryPage = start
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

    path := strings.Split(src, "/")
    name := strconv.Itoa(int(time.Now().Unix()))
    str := "null.jpg"
    if len(path) > 1 {
        str = path[len(path)-1]
    }
    imgName := name + "-" + str
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
	// <a>苏教版牛津小学英语一年级上册：contents</a>
	volume = 0
	slice1 := strings.Split(str, "：")
	// 章节
	materialName := slice1[0]
	
	// 上下册
	if strings.Contains(materialName, "上") {
		volume = int8(0)
	} else {
		volume = int8(1)
	}
	// 年级
	if strings.Contains(materialName, "一年级") {
		gradeName = "一年级"
	} else if strings.Contains(materialName, "二年级") {
		gradeName = "二年级"
	} else if strings.Contains(materialName, "三年级") {
		gradeName = "三年级"
	} else if strings.Contains(materialName, "四年级") {
		gradeName = "四年级"
	} else if strings.Contains(materialName, "五年级") {
		gradeName = "五年级"
	} else {
		gradeName = "六年级"
	}
	return int8(volume), gradeName
}

func collectImage(src string) (img []byte, name string) {
	doc, err := goquery.NewDocument(src)
	if err != nil {
		log.Fatal(err)
	}
	//knowledge <h1 class="tit-art tc">北师大版二年级英语课本上册：Unit1 Hello</h1>
	title := doc.Find("div.wrapper").Find("div.tk-article").Find("h1").First().Text()
	title = utfString(title)
	titleSlice := strings.Split(title, "：")
	knowledgeName := titleSlice[1]
	// 抓取图片
	imgLink := []string{}
	// 第一页
	doc.Find("div.content p").Find("img").Each(func(index int, s *goquery.Selection){
		imgsrc, ok := s.Attr("src")
		if ok != true {
		fmt.Println("img-src-err", err)
		} else {
		imgLink = append(imgLink, getImage(imgsrc))
		} 
	})
	//判断next page
	pages := doc.Find("div.content").Find("div.btn-pages a")
	fmt.Println("当前页面<a>总数", pages.Length())
	if pages.Length() <= 3 {
		imgJson, _ := json.Marshal(imgLink)
		return imgJson, knowledgeName
	}
	for index := range pages.Nodes {
		if index == 0  {
			continue
		}
		if (index + 2) >= pages.Length() {
			break
		}
		fmt.Println("开始抓取下一页图片")
		asrc, ok2 := pages.Eq(index).Attr("href")
		if ok2 == false {
			continue
		} 
		fmt.Println("下一页地址：", asrc)
		nextdoc, err := goquery.NewDocument(asrc)
		if err != nil {
			log.Fatal(err)
		}

		nextdoc.Find("div.content p").Find("img").Each(func(index int, s *goquery.Selection){
			imgsrc, ok := s.Attr("src")
			if ok != true {
			fmt.Println("img-src-err", err)
			} else {
			imgLink = append(imgLink, getImage(imgsrc))
			} 
		})
		fmt.Println("下一页抓取完毕")
		
	}
	imgJson, _ := json.Marshal(imgLink)
	return imgJson, knowledgeName
}