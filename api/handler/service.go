package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yarikyarichek/streamer/api/presenter"
	"github.com/yarikyarichek/streamer/config"
	"github.com/yarikyarichek/streamer/infostructure/repository"
	"github.com/yarikyarichek/streamer/usercase/mq"
)

type service struct {
	mq   mq.Service
	repo repository.Service
}

func NewService(mqSrv mq.Service, repo repository.Service) Service {
	return &service{mqSrv, repo}
}

func (srv *service) Start() error {
	http.HandleFunc("/messages", srv.messages)
	http.HandleFunc("/clear", srv.clear)
	return http.ListenAndServe(config.API_HOST+":"+config.API_PORT, nil)
}

func (srv *service) messages(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":

		fmt.Println(1)

		var m presenter.GetMessageRequest

		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(2)

		if len(bytes) > 0 {
			fmt.Println(21)
			err = json.Unmarshal(bytes, &m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		fmt.Println(3)

		messeges, err := srv.repo.Get(m.ToMessage(), m.Offset, m.ValidateLimit())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Println(4, *messeges)

		srv.handleResponse(w, messeges)

	case "POST":

		var m presenter.CreateMessegeRequests
		err := json.NewDecoder(req.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = srv.repo.Create(m.ToMessage())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		srv.handleResponse(w, &presenter.CreateMessegeResponse{Status: "created"})

	default:
		http.Error(w, "Sorry, only GET and POST methods are supported.", http.StatusBadRequest)
		return
	}

}

func (srv *service) clear(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		srv.mq.Clear()

	default:
		http.Error(w, "Sorry, only POST method is supported.", http.StatusBadRequest)
	}
}

func (srv *service) handleResponse(w http.ResponseWriter, resp interface{}) {
	r, _ := json.Marshal(resp)
	w.Write(r)
	w.Header().Set("Content-Type", "application/json")
}
