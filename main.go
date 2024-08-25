package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/daodao97/xgo/xapp"
	"github.com/daodao97/xgo/xlog"
	"github.com/gin-gonic/gin"
)

//go:embed templates/*
var templatesFS embed.FS

var Version string

func main() {
	app := xapp.NewApp().
		AddStartup().
		AfterStarted(func() {
			xlog.Debug("version", xlog.String("version", Version))
		}).
		AddServer(xapp.NewHttp(xapp.Args.Bind, h))

	if err := app.Run(); err != nil {
		fmt.Printf("Application error: %v\n", err)
	}
}

func h() http.Handler {
	e := xapp.NewGin()

	// 使用嵌入的方式加载模板
	templ := template.Must(template.New("").ParseFS(templatesFS, "templates/*"))
	e.SetHTMLTemplate(templ)

	e.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "文本对比工具",
		})
	})

	e.GET("/tree", func(c *gin.Context) {
		d1 := c.Query("d1")
		d2 := c.Query("d2")

		files1, err1 := getFiles(d1)
		files2, err2 := getFiles(d2)

		files1Map := make(map[string]bool)
		files2Map := make(map[string]bool)
		allPaths := make(map[string]bool)

		for _, file := range files1 {
			files1Map[file] = true
			allPaths[file] = true
		}
		for _, file := range files2 {
			files2Map[file] = true
			allPaths[file] = true
		}

		sortedPaths := make([]string, 0, len(allPaths))
		for path := range allPaths {
			sortedPaths = append(sortedPaths, path)
		}
		sort.Strings(sortedPaths)

		c.HTML(http.StatusOK, "tree.html", gin.H{
			"title":     "文件树对比",
			"d1":        d1,
			"d2":        d2,
			"files1Map": files1Map,
			"files2Map": files2Map,
			"allPaths":  sortedPaths,
			"err1":      err1,
			"err2":      err2,
		})
	})

	e.GET("/diff", func(c *gin.Context) {
		file1 := c.Query("file1")
		file2 := c.Query("file2")

		content1, err1 := readFile(file1)
		content2, err2 := readFile(file2)

		c.HTML(http.StatusOK, "diff.html", gin.H{
			"title":    "文本对比",
			"file1":    file1,
			"file2":    file2,
			"content1": content1,
			"content2": content2,
			"err1":     err1,
			"err2":     err2,
		})
	})

	// file_content
	e.GET("/file_content", func(c *gin.Context) {
		file := c.Query("file")
		content, err := readFile(file)
		c.JSON(http.StatusOK, gin.H{
			"file":    file,
			"content": content,
			"err":     err,
		})
	})

	e.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return e.Handler()
}

func getFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("遍历文件夹错误: %v", err)
	}
	return files, nil
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件错误: %v", err)
	}
	return string(content), nil
}
