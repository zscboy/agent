package restore

import (
	"fmt"
	"g"
	"log"
	"net/http"
)

var (
	myMux = &myHttpServerMux{muxHandlers: http.NewServeMux(), ignoreToken: make(map[string]bool)}
)

type myHttpServerMux struct {
	muxHandlers *http.ServeMux
	ignoreToken map[string]bool
}

func (mux *myHttpServerMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	handler, pattern := mux.muxHandlers.Handler(r)
	log.Println("http pattern:", pattern)
	if handler == nil || pattern == "" {
		log.Printf("no handler for Request, path:%s\n", r.URL.Path)
		w.WriteHeader(404)
		return
	}

	ignore, ok := mux.ignoreToken[pattern]
	if ok && ignore {
		handler.ServeHTTP(w, r)
		return
	}

	var uniqueId, tok = VerifyToken(r)
	if !tok {
		log.Println("token expired, path is:", r.URL.Path)
		//		replyError(w, codedef.ERR_TOKEN_EXPIRED)
		return
	}

	r.Header.Add("uniqueId", uniqueId)
	handler.ServeHTTP(w, r)
}

func CreateHttpServer() {
	addr := g.Config().Http.Listen
	if addr == "" {
		return
	}
	portStr := fmt.Sprintf(":%d", g.Config().Http.Listen)
	s := &http.Server{
		Addr:    portStr,
		Handler: myMux,
		// ReadTimeout:    10 * time.Second,
		//WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 8,
	}

	go acceptRequest(s)
}

func acceptRequest(s *http.Server) {
	log.Printf("Http server listen at:%d\n", g.Config().Http.Listen)

	err := s.ListenAndServe()
	if err != nil {
		log.Println("Http server ListenAndServe failed:", err)
	}
}
