/**
 * Created by PhpStorm.
 * @file   init.go
 * @author 李锦 <lijin@cavemanstudio.net>
 * @date   2023/6/6 10:10
 * @desc   init.go
 */

package services

import (
	"github.com/druidcaesa/gotool"
	"github.com/druidcaesa/gotool/openfile"
	"github.com/x-module/utils/utils/xlog"
	"os"
	"path"
)

func init() {
	InitializeLogger(Log{
		Path: "data",
		File: "system.log",
		Mode: "debug",
	})
}

// Log 日志设置
type Log struct {
	Path string `yaml:"path"`
	File string `yaml:"file"`
	Mode string `yaml:"mode"`
}

// InitializeLogger 初始化日志配置
func InitializeLogger(config Log) {
	if !gotool.FileUtils.Exists(config.Path) {
		err := os.MkdirAll(config.Path, os.ModePerm)
		if err != nil {
			panic("init system error. make log data err.path:" + config.Path)
		}
	}
	// 日志文件
	fileName := path.Join(config.Path, config.File)
	if !gotool.FileUtils.Exists(fileName) {
		openfile.Create(fileName)
		if !gotool.FileUtils.Exists(fileName) {
			panic("init system error. create log file err. log file:" + fileName)
		}
	}
	xlog.InitLogger(config.Path, config.File, config.Mode)
}
