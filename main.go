package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var dirName string
	fmt.Print("请输入要创建的项目名: ")
	fmt.Scanln(&dirName)

	var template string
	fmt.Print("是否使用模板文件: y/n:	")
	fmt.Scanln(&template)

	// 检查目录是否已存在
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		log.Fatalf("目录 '%s' 已经存在。", dirName)
	} else {
		// 尝试创建目录
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			log.Fatalf("无法创建目录 '%s': %v", dirName, err)
		} else {
			fmt.Printf("目录 '%s' 创建成功。\n", dirName)
		}
	}

	//修改模板包名
	err := replaceInFile("./template/main.go", "app", dirName)
	if err != nil {
		log.Fatalf("替换包名失败 '%s': %v", dirName, err)
	}
	_ = replaceInFile("./template/migrate.go", "app", dirName)

	// 使用绝对路径进行后续操作
	projectPath, _ := filepath.Abs(dirName)

	// 创建多级目录
	dirsToCreate := []string{
		"config",
		"internal/db/models",
		"internal/messaging",
		"internal/middleware",
		"internal/redis",
		"internal/dao",
		"internal/controller",
		"internal/service",
		"internal/utils",
		"pkg/errors",
		"pkg/logger",
		"proto",
		"routes/api/v1",
		"routes/api/v2",
		"routes/api/docs",
		"static",
		"log",
		"templates",
	}
	for _, dir := range dirsToCreate {
		fullPath := filepath.Join(projectPath, dir)
		err := os.MkdirAll(fullPath, 0755)
		if err != nil {
			log.Fatalf("创建目录结构失败: %v", err)
		}
	}

	// 创建空文件
	filesToTouch := []string{
		"main.go",
		"Dockerfile",
		"Makefile",
		"README.md",
		".gitignore",
	}
	for _, file := range filesToTouch {
		fullPath := filepath.Join(projectPath, file)
		_, err := os.Create(fullPath)
		if err != nil {
			log.Fatalf("创建文件失败: %v", err)
		}
	}
	fmt.Println("项目目录初始化完毕")

	commands := [][]string{
		{"go", "mod", "init", filepath.Base(projectPath)},

		{"go", "get", "-u", "github.com/spf13/cobra@latest"},
		{"cobra-cli", "init"},

		//拷贝模板文件
		{"cp", "../template/config.toml", "./"},
		{"cp", "../template/config.go", "./config/"},
		{"cp", "../template/healthcheck.go", "./routes/"},
		{"cp", "../template/prometheus.go", "./routes/"},
		{"cp", "../template/router.go", "./routes/"},
		{"cp", "../template/logger.go", "./pkg/logger/"},
		{"cp", "../template/consumer.go", "./internal/messaging/"},
		{"cp", "../template/producer.go", "./internal/messaging/"},
		{"cp", "../template/producer_test.go", "./internal/messaging/"},
		{"cp", "../template/user.go", "./internal/db/models/"},
		{"cp", "../template/migrate.go", "./internal/db/"},
		{"cp", "../template/client.go", "./internal/redis"},
		{"cp", "../template/client_test.go", "./internal/redis"},
		{"cp", "../template/main.go", "./"},
		{"cp", "../template/Makefile", "./"},
		{"cp", "../template/Dockerfile", "./"},

		{"go", "mod", "tidy"},
	}

	if template == "y" {
		for _, cmdArgs := range commands {
			cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			cmd.Dir = projectPath // 为每条命令都设置相同的工作目录

			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("执行命令失败: %v, 错误信息: %s\n", cmd.Args, err)
				log.Printf("命令输出: %s\n", output)
			} else {
				log.Printf("命令 '%s' 执行成功\n", strings.Join(cmd.Args, " "))
			}
		}
	}

}

func replaceInFile(filePath string, oldStr, newStr string) error {
	// 打开文件读取
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建临时文件用于写入替换后的内容
	tmpFile, err := os.Create(filePath + ".tmp")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(tmpFile)

	// 逐行读取并替换
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		// 如果读到文件尾部则跳出循环
		if err == io.EOF {
			break
		}

		// 替换行内内容
		newLine := strings.ReplaceAll(line, oldStr, newStr)
		if _, err := writer.WriteString(newLine); err != nil {
			return err
		}
	}

	// 刷新缓冲区确保数据写入磁盘
	if err := writer.Flush(); err != nil {
		return err
	}

	// 替换文件
	if err := os.Rename(filePath+".tmp", filePath); err != nil {
		return err
	}

	return nil
}
