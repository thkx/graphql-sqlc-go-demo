package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/documentloaders"
)

func main() {
	// 打开 PDF
	file, err := os.Open("test.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 获取文件大小
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := fi.Size()

	// 创建 PDF loader
	loader := documentloaders.NewPDF(file, size)

	// 加载
	docs, err := loader.Load(context.Background())
	if err != nil {
		log.Fatal("PDF 加载失败:", err)
	}

	fmt.Println("✅ PDF 加载成功！")
	fmt.Println("页数:", len(docs))
	fmt.Println("内容:", docs[0].PageContent)
}
