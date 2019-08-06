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
	cmd := exec.Command("sh", "/app/wiki.sh")
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
	signature := MySha1(PayloadBody)
	return XHubSignature == signature, nil
}

// hmac-sha1
func MySha1(payloadBody []byte) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write(payloadBody)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func main() {
	// use flag to change args
	flag.StringVar(&port, "p", "8000", "listen and sever port")
	flag.StringVar(&secret, "p", "hongfeng", "deploy password")
	flag.Parse()

	router := gin.Default()
	router.GET("/deploy/wiki", gitPush)
	_ = router.Run(":" + port)
}

