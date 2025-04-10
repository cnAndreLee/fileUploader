package main

import (
	_ "bytes"
	"log"
	"os"
	"path/filepath"
	_ "path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jlaffaye/ftp"
)

func GetEnv(key string, defaultvar string) string {
	v, exist := os.LookupEnv(key)
	var res string

	if !exist {
		res = defaultvar
		log.Println("loading default config---" + key + " : " + res)
	} else {
		res = v
		log.Println("loading   env   config---" + key + " : " + res)
	}

	return res
}

func main() {

	host := GetEnv("AFU_FTP_HOST", "127.0.0.1")
	port := GetEnv("AFU_FTP_PORT", "2222")
	user := GetEnv("AFU_FTP_USER", "jyiot")
	pass := GetEnv("AFU_FTP_PASS", "jyiot123")
	// path := GetEnv("AFU_FTP_PATH", "/gbiot/source/sw/61082107780/")
	path := GetEnv("AFU_FTP_PATH", "/home/")
	source := GetEnv("AFU_SOURCE_FILE_PATH", "C:\\Users\\Andre\\Desktop")
	log.Println(host, port, user, pass, path, source)

	c, err := ftp.Dial(host+":"+port, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Print(err)
	}
	defer c.Quit()

	err = c.Login(user, pass)
	if err != nil {
		log.Print(err)
	}

	// 创建一个新的fsnotify Watcher实例
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 设置要监控的目录路径
	dirPath := source // 请替换为你的目标目录路径

	// 开始监控目录
	err = watcher.Add(dirPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// 循环处理事件
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// 检查是否是新增文件的事件
			if event.Op&fsnotify.Create == fsnotify.Create {
				// 打开新文件
				// 必须等待200毫秒，否则无法打开文件
				time.Sleep(200 * time.Millisecond)
				file, err := os.Open(event.Name)
				if err != nil {
					log.Println("error opening file:", err)
					continue
				}

				err = c.Stor(path+filepath.Base(file.Name()), file)
				file.Close()
				if err != nil {
					log.Println("error upload file:", err)
					continue
				}

				log.Println("uploaded file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}

		// 休眠1秒
		//time.Sleep(1 * time.Second)
	}

}
