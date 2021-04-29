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
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// replaceCmd represents the replace command
var replaceCmd = &cobra.Command{
	Use:   "replace",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("replace called")
		arrList := readCheckLog("check.txt") // 返回二维的字符串数组，【【name，start】】

		for i := 0; i < len(arrList); i++ {
			temp, err := strconv.Atoi(arrList[i][1])
			if err != nil {
				println(err)
				return
			}
			replaceArr := findStartEnd(arrList[i][0], temp) // 返回一个一维数组【name，start，end】

			tempStart, err := strconv.Atoi(replaceArr[1])
			if err != nil {
				println(err)
				return
			}
			fmt.Println(tempStart)

			tempEnd, err := strconv.Atoi(replaceArr[2])
			if err != nil {
				println(err)
				return
			}
			fmt.Println(tempEnd)

			replace(replaceArr[0], tempStart-1, tempEnd-1)
			time.Sleep(time.Second)
		}
	},
}

func init() {
	rootCmd.AddCommand(replaceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// replaceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// replaceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 找到start 和End
func findStartEnd(fileName string, start int) []string {
	file, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()
	fmt.Println("succ")

	goRoutine := []byte("go func")

	bracketSlice1 := byte('{')
	bracketSlice2 := byte('}')
	bracketCount := 1
	flag := 0
	row := 0
	//rowLast := 0
	mark1 := make([]int, 0) // 携程开始
	mark2 := make([]int, 0) // 携程收尾

	input := bufio.NewScanner(file)
	for input.Scan() {
		row += 1
		info := input.Bytes() // 拿到info内容

		if flag == 0 && bytes.Contains(info, goRoutine) {

			fmt.Printf("!!!!!!文件第%d行有写协程\n", row)
			mark1 = append(mark1, row)
			flag = 1

			continue
		}
		if flag == 1 {

			for _, v := range info {
				//fmt.Println(-11, v)
				if v == bracketSlice1 {
					bracketCount++
					continue
				}
				if v == bracketSlice2 {
					bracketCount--
					continue
				}
				if bracketCount == 0 {
					fmt.Println("sa", bracketCount, row)
					mark2 = append(mark2, row)
					flag = 0
					bracketCount = 1
					break
				}

			}

			continue
		}

	}

	l := len(mark1)
	finalArr := make([]string, 0)
	for i := 0; i < l; i++ {
		if mark1[i] == start {
			// 生产一个[filename,start,end的数组]
			finalArr = append(finalArr, fileName, strconv.Itoa(mark1[i]), strconv.Itoa(mark2[i]))
		}
	}
	return finalArr
}

// 替换操作
func replace(fileName string, start, end int) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	row := 0
	newBuffer := make([]byte, 0)

	str1 := "RecoverGo(func() {\n"
	str2 := "    } )\n"

	for {
		slice, err := r.ReadBytes('\n')
		if row == start {
			newBuffer = append(newBuffer, []byte(str1)...)
		} else if row == end {
			newBuffer = append(newBuffer, []byte(str2)...)
		} else {
			newBuffer = append(newBuffer, slice...)
		}

		fmt.Println(row, slice)
		row += 1
		if err == io.EOF { // 如果读取到文件末尾
			break
		}
	}

	file2, err := os.OpenFile(
		fileName,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		//os.O_WRONLY|os.O_CREATE|syscall.O_DIRECTORY,
		0666,
	)
	nums, err := file2.Write(newBuffer)
	if err != nil {
		return
	}
	file2.Close()
	fmt.Println(nums)
}

// replace同时替换
func replaceOnce(fileName string, arr [][]int) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	r := bufio.NewReader(file)
	row := 0
	newBuffer := make([]byte, 0)
	//add :=make([]byte,0)
	str1 := "RecoverGo(func() {\n"
	str2 := "    } )\n"
	ptr1 := 0 // start指针，指向arr数组中的start
	ptr2 := 0 // end指针，指向arr数组中的end
	for {
		slice, err := r.ReadBytes('\n')
		if row == arr[0][ptr1] {
			newBuffer = append(newBuffer, []byte(str1)...)
			if ptr1 < len(arr[0]) {
				ptr1++
			}

		} else if row == arr[1][ptr2] {
			newBuffer = append(newBuffer, []byte(str2)...)
			if ptr2 < len(arr[1]) {
				ptr2++
			}
		} else { // 不是go func 直接append
			newBuffer = append(newBuffer, slice...)
		}

		fmt.Println(row, slice)
		row += 1
		if err == io.EOF { // 如果读取到文件末尾
			break
		}
	}
	fmt.Println(newBuffer)
	file2, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		//os.O_WRONLY|os.O_CREATE|syscall.O_DIRECTORY,
		0666,
	)
	nums, err := file2.Write(newBuffer)
	if err != nil {
		return
	}
	//file2.Close()
	fmt.Println(nums)
}

//func readCheckLog(name string) (fileName string, start int) {
//读日志txt，返回
func readCheckLog(name string) [][]string {
	file, err := os.Open(name)

	if err != nil {
		return nil
	}
	r := bufio.NewReader(file)
	arr := make([][]string, 0)
	for {
		num1 := 0
		num2 := 0
		num3 := 0
		num4 := 0
		tempArr := make([]string, 0)
		slice, err := r.ReadBytes('\n')
		for p, v := range slice {
			if '*' == v {
				num1 = p
			}
			if '~' == v {
				num2 = p
			}
			if '^' == v {
				num3 = p
			}
			if '&' == v {
				num4 = p
			}
			//fmt.Println(num1, num2, num3, num4)
			if num1*num2*num3*num4 != 0 {
				tempName := slice[num1+1 : num2]
				tempRow := slice[num3+1 : num4]
				tempArr = append(tempArr, string(tempName), string(tempRow))
				fmt.Println(string(tempName), string(tempRow))
				break
			}

		}
		if len(tempArr) == 2 {
			arr = append(arr, tempArr)
		}
		if err == io.EOF { // 如果读取到文件末尾
			break
		}
	}
	fmt.Println(arr)
	return arr
}
