package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
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

// execute the shell scripts
func ReLaunch() {
	cmd := exec.Command("sh", shell)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
}

// verifySignature
func verifySignature(c *gin.Context) (bool, error) {
	PayloadBody, err := c.GetRawData()
	if err != nil {
		return false, err
	}
	// Get Header with X-Hub-Signature
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

func usage()  {
	_, _ = fmt.Fprintf(os.Stderr, `deploy version: deploy:1.0.5
Usage: deploy [-p port] [-path UriPath] [-sh DeployShell] [-pwd WebhookSecret]

Options:
`)
	flag.PrintDefaults()
}

func init() {
	// use flag to change args
	flag.StringVar(&port, "p", "8000", "listen and serve port")
	flag.StringVar(&secret, "pwd", "hongfeng", "deploy password")
	flag.StringVar(&path, "path", "/deploy/wiki", "uri serve path")
	flag.StringVar(&shell, "sh", "/app/wiki.sh", "deploy shell scritpt")
	flag.BoolVar(&h, "h", false, "show this help")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
	}
	// Disable Console Color, you don't need console color when writing the logs to file
	gin.DisableConsoleColor()
	// Logging to a file.
	f, _ := os.Create("/logs/gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// Use the following code if you need to write the logs to file and console at the same time.
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()
	router.GET(path, gitPush)
	_ = router.Run(":" + port)
}
