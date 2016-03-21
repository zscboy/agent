package http

import (
	"encoding/json"
	"github.com/open-falcon/agent/gentk"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

var (
	LISTEN_PORT         = "net.port.listen"
	PROBLEM             = "PROBLEM"
	WXSERVER_HTTPS_PORT = "4005"
	WXSERVER_HTTP_PORT  = "4002"
	WXSERVER_WEB_PORT   = "3001"
	WXSERVER_DEV_PORT   = "17273"
	REDIS_PORT          = "6379"
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
		case WXSERVER_HTTPS_PORT:
		case WXSERVER_HTTP_PORT:
		case WXSERVER_WEB_PORT:
		case WXSERVER_DEV_PORT:
			restartWxserver()
			break

		case REDIS_PORT:
			restartRedis()
			break
		default:
			test()
			break
		}

	})
}

func restartWxserver() {
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

	fmt.Println("run cmd " + dir + "/restartWxserver.sh")
}

func test() {
	cmd := "./test.sh"
	out, err := exec.Command("/home/riguang/open-falcon/src/github.com/open-falcon/agent/test.sh").Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}
