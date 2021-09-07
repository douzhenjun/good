/*
 * @Description:
 * @version:
 * @Company: iwhalecloud
 * @Author: Dou
 * @Date: 2021-02-07 16:32:07
 * @LastEditors: Dou
 * @LastEditTime: 2021-02-07 16:32:07
 */

package utils

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func LoggerInfo(info ...interface{}) {
	golog.Info(info...)
}

func LoggerError(err error) {
	if err != nil {
		golog.Error(err)
	}
}

func LoggerErrorP(err *error) {
	LoggerError(*err)
}

func ErrorContains(err error, target string) bool {
	return strings.Contains(err.Error(), target)
}

/*
S2JMap Json字符串转为Map类型
*/
func S2JMap(s string) (map[string]interface{}, error) {
	var jm = make(map[string]interface{})
	err := json.Unmarshal(Str2bytes(s), &jm)
	return jm, err
}

/*
S2JMap2 Json字符串转为[]Map类型
*/
func S2JMap2(s string) ([]map[string]interface{}, error) {
	var jm2 = make([]map[string]interface{}, 0)
	err := json.Unmarshal(Str2bytes(s), &jm2)
	return jm2, err
}

/*
处理json格式的字符串, 将其转为raw类型, 序列化时可直接将字符串原样输出
不处理时, 字符串中的 " 会被加上转义符 \
*/
func RawJson(rawString string) (rawJson json.RawMessage) {
	_ = json.Unmarshal(Str2bytes(rawString), &rawJson)
	return
}

func RandCode(n int) string {
	randBytes := make([]byte, n/2)
	_, _ = rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

/*
查找字符串中的数字, 返回第一个匹配
*/
func NumberInString(s string) int {
	var regexpNumber = regexp.MustCompile("[0-9]+")
	match := regexpNumber.Find([]byte(s))
	number, _ := strconv.Atoi(string(match))
	return number
}

/*
通过指针转换[]byte和string, 无内存拷贝
*/
func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

/**
 * 根据时间类型字符串返回对应时间戳
 */
func GetTimeByDataString(timgStr string) int64 {
	timeTemplate1 := "2006-01-02 15:04:05"                               //常规类型
	stamp, _ := time.ParseInLocation(timeTemplate1, timgStr, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	return stamp.Unix()                                                  //输出：1546926630
}

/**
 * 获取当前时间字符串
 */
func GetStringDatetime() string {
	dataStr := fmt.Sprintf("%d", time.Now().Unix())
	return dataStr
}

/**
 * 判断某个路径是否存在
 * 返回两个值：第一个值为路径是否存在；第二个值返回error
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 按天生成日志文件
func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

// 创建打开文件
func newLogFile() *os.File {
	filename := todayFilename()
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func NewRequestLogger() (h iris.Handler, close func() error) {
	const deleteFileOnExit = false
	close = func() error { return nil }
	c := logger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Columns: true,
	}
	logFile := newLogFile()
	close = func() error {
		err := logFile.Close()
		if deleteFileOnExit {
			err = os.Remove(logFile.Name())
		}
		return err
	}
	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
		output := logger.Columnize(now.Format("2006/01/02 - 15:04:05"), latency, status, ip, method, path, message, headerMessage)
		logFile.Write([]byte(output))
	}
	h = logger.New(c)
	return
}

// 保留两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
	return value
}

func ResolveTime(seconds int) string {
	//天
	day := seconds / 86400
	//小时
	hour := seconds % 86400 / 3600
	//每分钟秒数
	minute := seconds % 86400 % 3600 / 60
	//每分钟秒数
	new_seconds := seconds % 86400 % 3600 % 60
	if day > 0 {
		return strconv.Itoa(day) + "天" + strconv.Itoa(hour) + "时" + strconv.Itoa(minute) + "分" + strconv.Itoa(new_seconds) + "秒"
	} else if day == 0 && hour > 0 {
		return strconv.Itoa(hour) + "时" + strconv.Itoa(minute) + "分" + strconv.Itoa(new_seconds) + "秒"
	} else if day == 0 && hour == 0 && minute > 0 {
		return strconv.Itoa(minute) + "分" + strconv.Itoa(new_seconds) + "秒"
	} else {
		return strconv.Itoa(seconds) + "秒"
	}
}

func ResolveSize(size int, unit string) string {
	return strconv.Itoa(size) + ""
}

//  数组去重
func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}

// MB的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024))
	} else {
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024))
	}
}

// MB的单位转换 保留两位小数
func Formattime(fileSize int64) (size string) {
	if fileSize < 60 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fs", float64(fileSize)/float64(1))
	} else if fileSize < (60 * 60) {
		return fmt.Sprintf("%.2fmin", float64(fileSize)/float64(60))
	} else if fileSize < (60 * 60 * 60) {
		return fmt.Sprintf("%.2fhour", float64(fileSize)/float64(60*60))
	} else if fileSize < (60 * 60 * 60 * 24) {
		return fmt.Sprintf("%.2fday", float64(fileSize)/float64(60*60*60))
	} else if fileSize < (60 * 60 * 60 * 24 * 30) {
		return fmt.Sprintf("%.2fmonth", float64(fileSize)/float64(60*60*60*24))
	} else {
		return fmt.Sprintf("%.2fyear", float64(fileSize)/float64(60*60*60*24*30))
	}
}
