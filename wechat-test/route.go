package main

import (
	"io"
	"net/http"
	"time"
	"log"
	"github.com/leeeboo/wechat-new/wx"
	"regexp"
)


type Jobber interface {
	job(w http.ResponseWriter, r *http.Request)
}

type Base struct {
	Method   string
	Pattern  string
}
type Get Base
type Post Base

type Workers struct {
	num int
	Get
	Post
}

var GET = Get{"GET", "^/"}
var POST = Post{"POST", "^/"}
var workers Workers



func (g *Get) job(w http.ResponseWriter, r *http.Request)  {
	client, err := wx.NewClient(r, w, token)

	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	if len(client.Query.Echostr) > 0 {
		w.Write([]byte(client.Query.Echostr))
		return
	}

	w.WriteHeader(403)
	return
}


func (p *Post) job(w http.ResponseWriter, r *http.Request)  {
	client, err := wx.NewClient(r, w, token)


	if err != nil {
		log.Println(err)
		w.WriteHeader(403)
		return
	}

	client.Run()
	return
}


func init() {
	workers = Workers{2,GET, POST}
}

type httpHandler struct {
}

func (*httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	if workers.Get.Method == r.Method{
		if m, _ := regexp.MatchString(workers.Get.Pattern, r.URL.Path); m{
			workers.Get.job(w, r)
			go writeLog(r, t, "unmatch", "")
		}
	} else {
		if workers.Post.Method == r.Method{
			if m, _ := regexp.MatchString(workers.Post.Pattern, r.URL.Path); m{
				workers.Post.job(w, r)
				go writeLog(r, t, "unmatch", "")
			}
		}
	}
	io.WriteString(w, "")
	return
}
