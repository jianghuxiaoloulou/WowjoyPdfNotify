package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/Unknwon/goconfig"
	"github.com/wonderivan/logger"
)

// 定义全局变量
var (
	url           string // 通知服务器的url
	fileSavePath  string // pdf 保存在服务器的路径
	transportType int    // pdf 文件传输方式
	wg            sync.WaitGroup
)

type pdfPathMsg struct {
	Code int    `json:"code"`
	Path string `json:"path"`
	Msg  string `json:"msg"`
}

func main() {
	// 保存PDF文件路径的channel
	//pdfPath := make(chan string)
	// 获取可执行文件的路径
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)
	// 通过配置文件配置
	logger.SetLogger(dir + "/log.json")
	//2、读取配置文件
	logger.Debug(dir)
	readConfigFile(dir)
	// 获取pdf文件路径
	wg.Add(1)
	go getPdfPath()
	// 等待所有的任务完成
	wg.Wait()
	logger.Debug("程序结束")
}

// 获取pdf文件路径
func getPdfPath() {
	defer wg.Done()
	logger.Debug("开始执行获取pdf路径协程.....")
	arg_num := len(os.Args)
	logger.Debug(arg_num)
	if arg_num < 2 {
		logger.Debug("没有获取到文件路径")
		return
	} else {
		_dir := os.Args[1]
		logger.Debug("获取的文件路径是：", _dir)

		// 通过HTTP POST接口通知路径
		// 方法1：数据流传输
		// 方法2 ：共享路径传输
		switch transportType {
		case 1:
			postFile(_dir, url)
		case 2:
			var pdfpath pdfPathMsg
			pdfpath.Code = 1
			pdfpath.Path = _dir
			pdfpath.Msg = "successful"
			if bs, err := json.Marshal(pdfpath); err == nil {
				httpPostJson(bs, url)
			}
		}
	}
}

// 初始化全局变量
func readConfigFile(dir string) {
	logger.Debug("开始读取配置文件....", dir)
	cfg, err := goconfig.LoadConfigFile(dir + "/config.ini")
	if err != nil {
		logger.Debug("无法加载配置文件：%s", err)
	}
	url, _ = cfg.GetValue("PDFNotify", "url")
	fileSavePath, _ = cfg.GetValue("PDFNotify", "fileSavePath")
	transportType, _ = cfg.Int("PDFNotify", "transportType")
	logger.Debug("url: ", url)
	logger.Debug("fileSavePath: ", fileSavePath)
	logger.Debug("transportType: ", transportType)
}

// http post
func httpPostJson(jsonstr []byte, url string) {
	logger.Debug("通知服务端的json: ", string(jsonstr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonstr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	logger.Debug("接口返回数据是：")

	statuscode := resp.StatusCode
	logger.Debug("statuscode: ", statuscode)

	body, _ := ioutil.ReadAll(resp.Body)
	logger.Debug("body：", string(body))
}

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("file", filename)
	if err != nil {
		logger.Error("error writing to buffer")
		return err
	}
	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		logger.Error("error opening file")
		return err
	}

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	fh.Close()
	logger.Debug(targetUrl)
	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Debug(string(resp_body))
	err = os.Remove(filename)
	if err != nil {
		logger.Debug("文件删除失败 ", filename)
		logger.Error(err)
	} else {
		logger.Debug("文件上传成功")
	}
	return nil
}
