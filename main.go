package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const tarFirstDir = "file"
const splitChar = ","

var dir, host, skipDir string
var errorFlag bool

func Init() {
	flag.StringVar(&dir, "dir", "", "Compressed folder,[/var/www]")
	flag.StringVar(&host, "host", "", "Server IP port,[192.168.1.1:8888]")
	flag.StringVar(&skipDir, "skip", "", "skip folder")
	flag.BoolVar(&errorFlag, "err", true, "Exit in case of exception,Default true")
}
func main() {
	Init()
	flag.Parse()
	if dir == "" || host == "" {
		fmt.Println("upfile -dir /var/www -host 192.168.1.1:8888")
		fmt.Println("upfile -dir /var/www -host 192.168.1.1:8888 -err flase")
		fmt.Println("upfile -dir /var/www -host 192.168.1.1:8888 -skip \"/var/www/log,/var/www/upload\"")
		flag.PrintDefaults()
		os.Exit(0)
	}
	_, err := os.Stat(dir)
	if err != nil {
		log.Fatalln("-dir error,", err)
	}
	start(host, dir, "nc")
}

func start(host, dir, module string) {
	switch module {
	case "nc":
		conn, err := net.Dial("tcp", host)
		if err != nil {
			log.Fatalln("Connect "+host+"error,", err)
		}
		defer conn.Close()
		err = sendData(dir, conn)
		if err != nil {
			log.Fatalln("sendData error,", err)
		}
	default:
		log.Fatalln("功能还未开发")
	}

}

//发送压缩文件
func sendData(path string, conn net.Conn) error {
	gw := gzip.NewWriter(conn)
	tw := tar.NewWriter(gw)
	defer gw.Close()
	defer tw.Close()
	return filepath.Walk(path, func(fileName string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return ProcessError(err)
		}
		//添加排除文件夹
		if excludeFile(fileName, fileInfo, skipDir) {
			return filepath.SkipDir
		}
		//兼容windows
		fileName = strings.ReplaceAll(fileName, "\\", "/")
		fileHeader, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return ProcessError(err)
		}
		//替换绝对路径
		fileName = filepath.Clean(fileName)
		fileHeader.Name = tarFirstDir + "/" + strings.TrimPrefix(strings.ReplaceAll(fileName, filepath.Dir(fileName), ""), string(filepath.Separator))
		fileHeader.Format = tar.FormatGNU
		if err := tw.WriteHeader(fileHeader); err != nil {
			return ProcessError(err)
		}
		if !fileInfo.Mode().IsRegular() {
			return ProcessError(err)
		}
		fileRead, err := os.Open(fileName)
		defer fileRead.Close()
		if err != nil {
			return ProcessError(err)
		}
		_, err = io.Copy(tw, fileRead)
		if err != nil {
			return err
		}
		return nil
	})
}

//排除文件夹或文件类型
func excludeFile(filename string, fileInfo os.FileInfo, skip string) bool {
	if skip == "" {
		return false
	}
	skips := strings.Split(skip, splitChar)
	//排除文件夹
	if fileInfo.IsDir() {
		for _, str := range skips {
			if filepath.Dir(filename) == str || filename == str {
				return true
			}
		}
	} else { //排除文件
		for _, str := range skips {
			if !strings.HasPrefix(str, ".") {
				continue
			}
			if strings.HasSuffix(fileInfo.Name(), str) {
				return true
			}
		}
	}
	return false
}

func ProcessError(err error) error {
	if errorFlag {
		return err
	}
	return nil
}
