package rpc

import (
	"fmt"
	"github.com/gitaction/martini"
	"github.com/sunwei/gitaction/src/adapters/outbound/persistece"
	"github.com/sunwei/gitaction/src/application/usecase/gitserver"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const SendPackRequestService = "git-receive-pack"
const FetchPackRequestService = "git-upload-pack"
const SendPackInfoPath = "/info/refs"
const ReceivePackPath = "/git-upload-pack"

func NewGitRouter() martini.Router {
	r := martini.NewRouter()
	r.Get(SendPackInfoPath, InfoRefs)
	r.Post(ReceivePackPath, ReceivePack)
	return r
}

func ReceivePack(res http.ResponseWriter, req *http.Request)  {
	
}

func InfoRefs(res http.ResponseWriter, req *http.Request)  {
	serviceName := getService(req.URL)
	if !isServiceAvailable(serviceName) {
		http.Error(res, "Not Found", 404)
		return
	}
	repoName := getRepoName(req.URL, SendPackInfoPath)
	repository := persistece.NewRepoRepositoryImpl()

	res.Header().Add("Content-Type", fmt.Sprintf("application/x-%s-advertisement", serviceName))
	res.Header().Add("Cache-Control", "no-cache")
	res.WriteHeader(200)

	if err := pktLine(res, fmt.Sprintf("# service=%s\n", serviceName)); err != nil {
		fmt.Printf("%+v\n", err)
	}
	
	usecase := gitserver.NewGitSmartUseCase(repoName, serviceName, repository)
	chunksInfo, err := usecase.GetRefsInfo()
	fmt.Printf("%s\n", chunksInfo)
	_, _ = res.Write(chunksInfo)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	//if err := pktFlush(res); err != nil {
	//	fmt.Printf("%+v\n", err)
	//}
}

func getRepoName(url *url.URL, path string) string {
	//fmt.Printf("%v - %v\n", url, path)
	//return "test123"
	return strings.TrimSuffix(strings.TrimPrefix(strings.TrimSuffix(url.Path, path), "/"), ".git")
}

func getService(url *url.URL) string {
	return url.Query().Get("service")
}

func isServiceAvailable(service string) bool {
	return service == SendPackRequestService || service == FetchPackRequestService
}

func pktLine(w io.Writer, s string) error {
	_, err := fmt.Fprintf(w, "%04x%s", len(s)+4, s)
	return err
}

//TODO: check is this necessary, MacOS & Linux - ubuntu
//func pktFlush(w io.Writer) error {
//	_, err := fmt.Fprint(w, "0000")
//	return err
//}
