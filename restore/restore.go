package restore

import (
	"codedef"
	"fmt"
	"gentk"
	"github.com/garyburd/redigo/redis"
	"log"
	"net/http"
	"rhelper"
)

func restore(w http.ResponseWriter, r *http.Request) {
	//var account = r.Header.Get("account")
	//var password = r.Header.Get("password")
	//var id = r.Header.Get("id")
	log.Printf("onLgin")

	//log.Printf("login, account:%s, password:%s\n", account, password)

	var tk = gentk.GenTK(id)
	log.Printf("login, tk:%s\n", tk)

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("tk", tk)
	w.Write([]byte(fmt.Sprintf("{\"tk\":\"%s\", \"id\":\"%s\"}", tk, id)))
}

func init() {
	myMux.muxHandlers.HandleFunc("/restore", restore)
}
