#!/bin/bash

echo "请输入要创建的项目名:"
read dir_name
# 检查目录是否已存在
if [ -d "$dir_name" ]; then
    echo "目录 '$dir_name' 已经存在。"
else
    # 尝试创建目录
    mkdir "$dir_name"
    if [ $? -eq 0 ]; then
        echo "目录 '$dir_name' 创建成功。"
    else
        echo "无法创建目录 '$dir_name'。"
        exit 1
    fi
fi

# 改变当前终端会话到新创建的目录
cd "$dir_name"

# 初始化项目
go mod init

# 创建顶级目录
mkdir -p config
mkdir -p internal/{db/{migration,models},middleware,dao,service,utils}
mkdir -p pkg/{errors,logger}
mkdir -p routes/{api/{v1,v2},docs}
mkdir -p static templates

# 创建文件
touch config/config.go
touch internal/db/models/{base_model.go,user.go}
touch pkg/{errors/errors.go,logger/logger.go}
touch routes/healthcheck.go
touch main.go
touch config.tml
touch Dockerfile
touch Makefile
touch README.md
touch .gitignore

echo "项目目录初始化完毕"


# 安装cobra 
read -p "是否需要使用 cobra ? (y/n): " answer  
if [[ $answer =~ ^(yes|y|YES|Y)$ ]]; then  
    go get -u github.com/spf13/cobra@latest
    if [ $? == 0 ]
    then
        echo "cobra 安装成功"
        # 初始化命令行
        cobra-cli init
    else
        echo "cobra 安装失败"
    fi
elif [[ $answer =~ ^(no|n|NO|N)$ ]]; then  
    echo "跳过安装 cobra."  
else  
    echo "Invalid input. Please enter 'yes' or 'no'."  
fi

# 安装gin  
read -p "是否需要安装 gin? (y/n): " answer  
if [[ $answer =~ ^(yes|y|YES|Y)$ ]]; then  
    go get -u github.com/gin-gonic/gin 
    if [ $? == 0 ]
    then
        echo "gin 安装成功"
    else
        echo "gin 安装失败"
    fi
elif [[ $answer =~ ^(no|n|NO|N)$ ]]; then  
    echo "跳过安装 gin."  
else  
    echo "Invalid input. Please enter 'yes' or 'no'."  
fi


# 安装gorm
read -p "是否需要安装 gorm? (y/n): " answer  
if [[ $answer =~ ^(yes|y|YES|Y)$ ]]; then  
    go get -u gorm.io/gorm
    if [ $? == 0 ]
    then
        echo "gorm 安装成功"
    else
        echo "gorm 安装失败"
    fi 
elif [[ $answer =~ ^(no|n|NO|N)$ ]]; then  
    echo "跳过安装 gorm."  
else  
    echo "Invalid input. Please enter 'yes' or 'no'."  
fi

# 安装gorm驱动
read -p "是否需要安装 gorm 驱动? (y/n): " answer  
if [[ $answer =~ ^(yes|y|YES|Y)$ ]]; then  
    go get -u gorm.io/driver/postgres
    go get -u gorm.io/driver/mysql
    go get -u gorm.io/driver/sqlite
    go get -u gorm.io/driver/sqlserver
    if [ $? == 0 ]
    then
        echo "gorm 驱动 安装成功"
    else
        echo "gorm 驱动 安装失败"
    fi 
elif [[ $answer =~ ^(no|n|NO|N)$ ]]; then  
    echo "跳过安装 gorm 驱动."  
else  
    echo "Invalid input. Please enter 'yes' or 'no'."  
fi

if [ $? == 0 ]
then
    echo "执行成功"
else
    echo "执行失败"
fi
