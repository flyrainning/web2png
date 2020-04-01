package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/chromedp/chromedp"
)

var (
	cacheDir string
)

func init() {
	cacheDir = path.Join(GetWorkDir(), "cache")
	if !PathExists(cacheDir) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

}

func main() {

	//监听协议
	http.HandleFunc("/", screenshotFunc)
	//监听服务
	err := http.ListenAndServe("0.0.0.0:80", nil)

	if err != nil {
		fmt.Println("服务器错误")
	}

}

func screenshotFunc(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	url := strings.TrimPrefix(r.RequestURI, "/")
	url = strings.Replace(url, "http:/", "http://", 1)
	url = strings.Replace(url, "https:/", "https://", 1)
	if !strings.HasPrefix(url, "http") {
		if url == "clean" {
			if files, _, err := ListDir(cacheDir); err == nil {
				for _, v := range files {
					os.Remove(v)
				}
				w.Write([]byte("1"))
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fmt.Println("url:", url)

	hash := MD5(url)
	hashFile := path.Join(cacheDir, hash)

	isSync := true

	var buf []byte

	if PathExists(hashFile) {
		var err error
		buf, err = ReadFile(hashFile)
		if err == nil {
			isSync = false
		}
	}
	if isSync {
		// create context
		ctx, cancel := chromedp.NewContext(context.Background())
		// ctx, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222/devtools/browser/8d198c7b-a374-47c0-80cb-8506718709f3")
		defer cancel()
		// capture screenshot of an element

		if err := chromedp.Run(ctx, fullScreenshot(url, 100, &buf)); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		// if err := ioutil.WriteFile("fullScreenshot.png", buf, 0644); err != nil {
		// 	log.Fatal(err)
		// }
		WriteFile4byte(hashFile, &buf)
	}

	w.Header().Set("content-type", "image/png")
	w.Header().Set("cache-control", "no-cache, no-store, must-revalidate")
	w.Header().Set("expires", "0")
	w.WriteHeader(http.StatusOK)

	w.Write(buf)
}
