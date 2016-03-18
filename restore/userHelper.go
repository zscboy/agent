package restore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func replyError(w http.ResponseWriter, ec int) {
	w.Header().Set("error", fmt.Sprintf("%d", ec))
	w.Write(nil)
}

func replyBuffer(w http.ResponseWriter, buf []byte) {
	w.Header().Add("Content-Length", strconv.Itoa(len(buf)))
	w.WriteHeader(200)
	w.Write(buf)
}

func replyJSON(w http.ResponseWriter, j interface{}) {
	buf, _ := json.Marshal(j)
	replyBuffer(w, buf)
}
