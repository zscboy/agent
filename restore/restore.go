package restore

import (
	"log"
	"net/http"
	"os/exec"
)

func restore(w http.ResponseWriter, r *http.Request) {
	//var account = r.Header.Get("account")
	//var password = r.Header.Get("password")
	//var id = r.Header.Get("id")
	cmd := "./restartWxserver.sh"
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)

	//log.Printf("login, account:%s, password:%s\n", account, password)

	//	var tk = GenTK(id)
	//log.Printf("login, tk:%s\n", tk)

	w.Header().Set("Content-Type", "text/plain")
	//w.Header().Set("tk", tk)
	w.Write([]byte("hello"))
}

func init() {
	myMux.muxHandlers.HandleFunc("/restore", restore)
}
