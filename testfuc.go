package main

import (
"io"
"log"
"net/http"
"os/exec"
)

func reLaunch() {
	cmd := exec.Command("sh", "/app/wiki.sh")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = cmd.Wait()
}

func deployPage(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "<h1> Hello, this is my deploy server")

	//reLaunch()
}

func main() {
	http.HandleFunc("/deploy/wiki", deployPage)
	_ = http.ListenAndServe(":8000", nil)
}
