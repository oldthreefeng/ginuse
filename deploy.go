package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/exec"
)

var (
	secret string
	port   string
	path   string
	shell  string
	h      bool
)

// return true then deploy
func gitPush(c *gin.Context) {
	matched, _ := VerifySignature(c)
	if !matched {
		err := "Signatures did not match"
		c.String(http.StatusForbidden, err)
		log.Warn(err)
		return
	}
	log.Info("Signatures is matched ~")
	//return 200 first
	c.String(http.StatusOK, "OK")
	ReLaunch(shell)
}

// aliyun code
func gitPushCode(c *gin.Context) {
	c.String(http.StatusOK, "ok")
	ReLaunch("/app/images.sh")
}

// execute the shell scripts
func ReLaunch(cmdStr string) {
	cmd := exec.Command("sh", cmdStr)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
}

// verifySignature
func VerifySignature(c *gin.Context) (bool, error) {
	PayloadBody, err := c.GetRawData()
	if err != nil {
		return false, err
	}
	// Get Header with X-Hub-Signature
	XHubSignature := c.GetHeader("X-Hub-Signature")
	signature := getSha1Code(PayloadBody)
	log.Info(signature)
	return XHubSignature == signature, nil
}

// hmac-sha1
func getSha1Code(payloadBody []byte) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write(payloadBody)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `deploy version: deploy:1.1.19
Usage: deploy [-p port] [-path UriPath] [-sh DeployShell] [-pwd WebhookSecret]

Options:
`)
	flag.PrintDefaults()
}

func defaultPage(g *gin.Context) {
	firstName := g.DefaultQuery("firstName", "test")
	lastName := g.Query("lastName")
	g.String(http.StatusOK, "Hello %s %s, This is My deploy Server~", firstName, lastName)
}

func init() {
	// use flag to change args
	flag.StringVar(&port, "p", "8000", "listen and serve port")
	flag.StringVar(&secret, "pwd", "hongfeng", "deploy password")
	flag.StringVar(&path, "path", "/deploy/wiki", "uri serve path")
	flag.StringVar(&shell, "sh", "/app/w.sh", "deploy shell scritpt")
	flag.BoolVar(&h, "h", false, "show this help")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}
	// Disable Console Color, you don't need console color when writing the logs to file
	gin.DisableConsoleColor()
	// Logging to a file.
	var f *os.File
	f, _ = os.OpenFile("logs/gin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET("/", defaultPage)
	router.POST(path, gitPush)
	router.POST("/aliyun/code", gitPushCode)
	_ = router.Run(":" + port)
}
