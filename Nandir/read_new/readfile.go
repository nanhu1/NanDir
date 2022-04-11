package Nandir

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func Readfile(filename string) []string {
	var s []string
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("字典加载失败！")
		return nil
	}
	reader := bufio.NewReader(file)
	for {
		url, err := reader.ReadString('\n') //注意是字符
		str1 := strings.Replace(url, "\n", "", -1)
		str := strings.Replace(str1, "\r", "", -1)
		if err == io.EOF {
			file.Close()
		}
		if err != nil {
			break
		}
		s = append(s, str)
		fmt.Println(str)
	}
	return s
}

type Nanlimit struct {
	Num int
	C   chan struct{}
}

func NewNan(num int) *Nanlimit {
	return &Nanlimit{
		Num: num,
		C:   make(chan struct{}, num),
	}
}

func (g *Nanlimit) Run(f func()) {
	g.C <- struct{}{}
	go func() {
		f()
		<-g.C
	}()
}
