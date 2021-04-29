/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
		doFindAllGoroutine()

	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func doFindAllGoroutine() {
	targetDir := "/Users/zongxingye/Desktop/pondtest"
	//targetDir := "/Users/zongxingye/Desktop/ttt"
	Delimiter := string(os.PathSeparator)
	level := 1
	recursionAllFile(level, Delimiter, targetDir)
}
func recursionAllFile(level int, Delimiter string, fileDir string) {
	files, _ := ioutil.ReadDir(fileDir)
	for _, onefile := range files {
		if onefile.IsDir() {
			// 是目录则递归
			recursionAllFile(level+1, Delimiter, fileDir+Delimiter+onefile.Name())
		} else {
			// 不是目录，判断文件格式是否为.go
			if strings.Contains(onefile.Name(), ".go") {
				// 是go文件，读取文件内容
				findGoroutine(fileDir + Delimiter + onefile.Name())
			} else {
				fmt.Printf("%s不是go文件跳过检测\n", onefile.Name())
			}
		}
	}
}

func findGoroutine(path string) (contain bool) {
	goRoutine := []byte("go func")
	comment := []byte("//")
	commentWithLinkBreak := []byte("`")
	file, err := os.Open(path)

	if nil != err {
		fmt.Println(err)
	}
	defer file.Close()
	// 开始逐行扫描

	input := bufio.NewScanner(file)
	var row = 0  // 记录一下当前的行数
	var mark = 0 // 记录换行字符串的情况
	//buffers := make([]byte, 0) // 记录文件名和行

	for input.Scan() {
		row += 1
		info := input.Bytes() // 拿到info内容
		if bytes.Contains(info, comment) {
			//fmt.Println("注释跳过")
			continue
		}
		if mark == 0 && bytes.Contains(info, commentWithLinkBreak) {
			//fmt.Println("换行注释跳过")
			mark = 1
			continue
		}
		if mark == 1 {
			continue
		}
		if mark == 1 && bytes.Contains(info, commentWithLinkBreak) {
			//fmt.Println("结束换行注释跳过")
			continue
		}

		if bytes.Contains(info, goRoutine) {
			//写入log文件
			logs, _ := os.OpenFile("check.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic("读写日志失败")
				return false
			}
			defer logs.Close()
			log.SetOutput(logs)
			log.Printf("*%s~文件第^%d&行有写协程\n", path, row)
			fmt.Printf("!!!!!!!%s文件第%d行有写协程\n", path, row)

			//buffers = append(buffers, []byte(fmt.Sprintf("%s,%d\n", path, row))...)
			//fmt.Println("123132", buffers)
			contain = true
			continue
		}
	}
	//file2, err := os.OpenFile(
	//	"check"+path+".txt",
	//	os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
	//	0666,
	//)
	//nums, err := file2.Write(buffers)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("nums", nums)
	//fmt.Println("bbb", buffers)
	return contain

}
