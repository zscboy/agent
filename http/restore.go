package http

import (
	"encoding/json"
	"fmt"
	"github.com/open-falcon/agent/g"
	"github.com/open-falcon/agent/gentk"
	"github.com/toolkits/file"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	LISTEN_PORT              = "net.port.listen"
	PROBLEM                  = "PROBLEM"
	TAGS_WXSERVER_HTTPS_PORT = "port:4005"
	TAGS_WXSERVER_HTTP_PORT  = "port:4002"
	TAGS_WXSERVER_WEB_PORT   = "port:3001"
	TAGS_WXSERVER_DEV_PORT   = "port:17273"
	TAGS_REDIS_PORT          = "port:6379"
)

func configRestoreRoutes() {
	http.HandleFunc("/restore", func(w http.ResponseWriter, r *http.Request) {
		var uniqueId, tok = gentk.VerifyToken(r)
		if !tok {
			log.Println("token expired, path is:", r.URL.Path)
			//		replyError(w, codedef.ERR_TOKEN_EXPIRED)
			return
		}

		log.Println("uniqueId:", uniqueId)

		if r.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		bs, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		type Event struct {
			Endpoint string `json:"endpoint"`
			Metric   string `json:"metric"`
			Status   string `json:"status"`
			Step     string `json:"step"`
			Priority string `json:"priority"`
			Time     string `json:"time"`
			TplId    string `json:"tpl_id"`
			ExpId    string `json:"exp_id"`
			StraId   string `json:"stra_id"`
			Tags     string `json:"tags"`
		}

		var event Event
		err = json.Unmarshal(bs, &event)
		if err != nil {
			log.Println("onCreateDelayTask failed, json decode failed:", err)
			return
		}

		metric := event.Metric
		status := event.Status
		tags := event.Tags

		if metric != LISTEN_PORT || status != PROBLEM {
			return
		}

		switch tags {
		case TAGS_WXSERVER_HTTPS_PORT, TAGS_WXSERVER_HTTP_PORT, TAGS_WXSERVER_WEB_PORT, TAGS_WXSERVER_DEV_PORT:
			restartWxserver()
			break

		case TAGS_REDIS_PORT:
			restartRedis()
			break
		default:
			test()
			break
		}

	})
}

func restartWxserver() {
	log.Println("restartWxserver")
	dir := g.Config().Plugin.Dir
	parentDir := file.Dir(dir)
	cmd := exec.Command(dir + "/restartWxserver.sh")
	cmd.Dir = parentDir
	err := cmd.Run()
	if err != nil {
		log.Printf("run cmd in dir:%s fail. error: %s", dir, err)
		return
	}

	fmt.Println("run cmd " + dir + "/restartWxserver.sh")
}

func restartRedis() {
	dir := g.Config().Plugin.Dir
	parentDir := file.Dir(dir)
	cmd := exec.Command(dir + "/restartRedis.sh")
	cmd.Dir = parentDir
	err := cmd.Run()
	if err != nil {
		log.Printf("run cmd in dir:%s fail. error: %s", dir, err)
		return
	}

	fmt.Println("run cmd " + dir + "/restartRedis.sh")
}

func test() {
	dir := g.Config().Plugin.Dir
	parentDir := file.Dir(dir)
	cmd := exec.Command(dir + "/test.sh")
	cmd.Dir = parentDir
	err := cmd.Run()
	if err != nil {
		log.Printf("run cmd in dir:%s fail. error: %s", dir, err)
		return
	}

	fmt.Println("run cmd " + dir + "/test.sh")
}
