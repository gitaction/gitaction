package rpc

import (
	"compress/gzip"
	"fmt"
	"github.com/gitaction/martini"
	"github.com/gitaction/src/adapters/outbound/persistece"
	"github.com/gitaction/src/application/usecase/gitserver"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const SendPackRequestService = "git-receive-pack"
const FetchPackRequestService = "git-upload-pack"
const SendPackInfoPath = "/info/refs"
const ReceivePackPath = "/git-receive-pack"

func NewGitRouter() martini.Router {
	r := martini.NewRouter()
	r.Get("/([a-zA-Z]+).git" + SendPackInfoPath, InfoRefs)
	r.Post("/([a-zA-Z]+).git" +ReceivePackPath, ReceivePack)
	return r
}

func newWriteFlusher(w http.ResponseWriter) io.Writer {
	return writeFlusher{w.(interface {
		io.Writer
		http.Flusher
	})}
}

type writeFlusher struct {
	wf interface {
		io.Writer
		http.Flusher
	}
}

func (w writeFlusher) Write(p []byte) (int, error) {
	defer w.wf.Flush()
	return w.wf.Write(p)
}

func ReceivePack(res http.ResponseWriter, req *http.Request)  {
	body, err := getBody(req)
	if err != nil {
		http.Error(res, "receive pack read body error", 500)
		return
	}

	repoName := getRepoName(req.URL, ReceivePackPath)
	repository := persistece.NewRepoRepositoryImpl()
	setResponse(res, fmt.Sprintf("application/x-%s-result", SendPackRequestService))
	usecase := gitserver.NewGitSmartUseCase(repoName, SendPackRequestService, repository)
	err = usecase.ReceivePack(body, newWriteFlusher(res))
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func getBody(req *http.Request) (io.Reader, error) {
	body := req.Body
	var err error
	if req.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(req.Body)
	}
	return body, err
}

func setResponse(res http.ResponseWriter, ct string)  {
	res.Header().Add("Content-Type", ct)
	res.Header().Add("Cache-Control", "no-cache")
	res.WriteHeader(200)
}

func InfoRefs(res http.ResponseWriter, req *http.Request)  {
	serviceName := getService(req.URL)
	if !isServiceAvailable(serviceName) {
		http.Error(res, "Not Found", 404)
		return
	}
	repoName := getRepoName(req.URL, SendPackInfoPath)
	repository := persistece.NewRepoRepositoryImpl()

	setResponse(res, fmt.Sprintf("application/x-%s-advertisement", serviceName))

	if err := pktLine(res, fmt.Sprintf("# service=%s\n", serviceName)); err != nil {
		fmt.Printf("%+v\n", err)
	}
	
	usecase := gitserver.NewGitSmartUseCase(repoName, serviceName, repository)
	err := usecase.GetRefsInfo(res)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	if err := pktFlush(res); err != nil {
		fmt.Printf("%+v\n", err)
	}
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
func pktFlush(w io.Writer) error {
	_, err := fmt.Fprint(w, "0000")
	return err
}
