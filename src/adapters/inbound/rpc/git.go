package rpc

import (
	"github.com/gitaction/martini"
	"net/http"
	"net/url"
)

const SendPackRequestService = "git-receive-pack"
const FetchPackRequestService = "git-upload-pack"

func NewGitRouter() martini.Router {
	r := martini.NewRouter()
	r.Get("/info/refs", InfoRefs)
	return r
}

func InfoRefs(res http.ResponseWriter, req *http.Request)  {
	serviceName := getService(req.URL)
	if !isServiceAvailable(serviceName) {
		http.Error(res, "Not Found", 404)
	}
	
	
	res.WriteHeader(200)
}

func getService(url *url.URL) string {
	return url.Query().Get("service")
}

func isServiceAvailable(service string) bool {
	return service == SendPackRequestService || service == FetchPackRequestService
}