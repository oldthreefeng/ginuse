package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os/exec"
)

var (
	secret string
	port   string
	path   string
	shell    string
)

//签名对了才能执行Relaunch
func gitPush(c *gin.Context) {
	matched, _ := verifySignature(c)
	if !matched {
		err := "Signatures did not match"
		c.String(http.StatusForbidden, err)
		fmt.Println(err)
		return
	}
	fmt.Println("Signatures is matched ~")
	ReLaunch()
	c.String(http.StatusOK, "OK")
}

// 执行部署脚本
func ReLaunch() {
	cmd := exec.Command("sh", shell)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
}

// 验证签名
func verifySignature(c *gin.Context) (bool, error) {
	PayloadBody, err := c.GetRawData()
	if err != nil {
		return false, err
	}
	// 获取请求头的签名信息
	XHubSignature := c.GetHeader("X-Hub-Signature")
	signature := getSha1Code(PayloadBody)
	fmt.Println(signature)
	return XHubSignature == signature, nil
}

// hmac-sha1
func getSha1Code(payloadBody []byte) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write(payloadBody)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func main() {
	// use flag to change args
	flag.StringVar(&port, "p", "8000", "listen and serve port")
	flag.StringVar(&secret, "pwd", "hongfeng", "deploy password")
	flag.StringVar(&path, "path", "/deploy/wiki", "url serve path")
	flag.StringVar(&shell, "sh", "/app/wiki.sh", "deploy shell scritpt")
	flag.Parse()

	router := gin.Default()
	router.GET(path, gitPush)
	_ = router.Run(":" + port)
}
