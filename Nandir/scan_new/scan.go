package Nandir

import (
	"crypto/tls"
	"github.com/gookit/color"
	"github.com/go-resty/resty/v2"
	"Nandir/read_new"
	"strings"
	"sync"
	"time"
)

type Webinfo struct {
	StatusCode int
	//Title      string
	Server  string
	Powered string
	Body    string
	Res     string //成功的结果
	Bodylen int    //返回包长度
}

func Nanscan(url string, dir string, errtime int64, timesleep int64) (Webinfo, error) {
	time.Sleep(time.Duration(timesleep) * time.Second) //设置延时时间
	var t string
	var Web Webinfo
	t = url
	a := t[len(t)-1:]
	if strings.Contains(t, "https://") {
		//log.Println("https")
	} else {
		t = "http://" + url
	}
	if a != "/" {
		t = url + "/" //判断结尾是否为/，如果不是，那就加上
	}
	client := resty.New().SetTimeout(time.Duration(errtime) * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //忽略https证书错误，设置超时时间
	client.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
	resp, err := client.R().EnableTrace().Get(t + dir) //开始请求扫描
	if err != nil {
		//log.Println(err)
		return Web, err
	}
	str := resp.Body()
	body := string(str)
	//re1 := regexp.MustCompile(title) //正则取标题
	//titlename := re1.FindAllStringSubmatch(body, 1)
	//fmt.Println(body)
	Web.StatusCode = resp.StatusCode()
	Web.Powered = resp.Header().Get("X-Powered-By")
	//Web.Title = titlename[0][1]
	Web.Server = resp.Header().Get("server")
	Web.Body = body
	Web.Res = t + dir
	Web.Bodylen = len(body)
	return Web, nil
}

func NanDirScan(url string, filename string, num int, errtime int64, timesleep int64, openerrstr int, errstr string) {
	dicall := read_new.Readfile(filename)
	r, err := Nanscan(url, "", errtime, 0)
	if err != nil {
		color.Warn.Println("[Err] 目标访问错误，可能被ban了！")
		return
	}
	color.Red.Println("[Info] NanDir|Try| By PuPp1T.; 版本:1.0")
	if openerrstr == 1 {
		color.Red.Println("[Info] 已开启自定义错误关键词:", errstr)
	}
	color.Red.Println("[Info] 目标地址:", url)
	color.Red.Println("[Info] 当前线程:", num)
	color.Red.Println("[Info] 超时时间:", errtime)
	color.Red.Println("[Info] 目标相关容器:", r.Server, r.Powered)
	color.Red.Println("[Info] 加载字典数量:", len(dicall))
	color.Red.Println("[Info] 开始扫描中.")
	color.Red.Println("---------------------------------------")

	g := read_new.NewNan(num) //设置线程数量
	wg := &sync.WaitGroup{}
	beg := time.Now()
	for i := 0; i < len(dicall); i++ {
		wg.Add(1)
		task := dicall[i]
		g.Run(func() {
			respBody, err := Nanscan(url, task, errtime, timesleep)
			if err != nil {
				//color.Warn.Println("目标访问错误，可能被ban了！")
				wg.Done()
				return
			}
			if strings.Contains(respBody.Body, errstr) == false {
				if respBody.StatusCode == 200 {
					color.Info.Println("[200] ", respBody.Res+"   [len]", respBody.Bodylen)
					//writefile.Write(url, "[200] "+respBody.Res+"\n")
				} else if respBody.StatusCode == 403 {
					color.Warn.Println("[403] ", respBody.Res+"   [len]", respBody.Bodylen)
					//writefile.Write(url, "[403] "+respBody.Res+"\n")
				} else if respBody.StatusCode == 302 {
					color.Warn.Println("[302] ", respBody.Res+"   [len]", respBody.Bodylen)
					//writefile.Write(url, "[302] "+respBody.Res+"\n")
				}
			}
			wg.Done()
		})
	}
	wg.Wait()
	color.Red.Printf("[info] 扫描完成！当前用时: %fs", time.Now().Sub(beg).Seconds())
}
