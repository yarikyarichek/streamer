package handler

import (
	"encoding/json"
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
	srv.mq.Start()
	http.HandleFunc("/messages", srv.messages)
	http.HandleFunc("/clear", srv.clear)
	return http.ListenAndServe(config.API_HOST+":"+config.API_PORT, nil)
}

func (srv *service) messages(w http.ResponseWriter, req *http.Request) {

	var resp interface{}

	switch req.Method {

	case "GET":

		var m presenter.GetMessageRequest

		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			srv.handleError(w, err)
			return
		}

		if len(bytes) > 0 {
			err = json.Unmarshal(bytes, &m)
			if err != nil {
				srv.handleError(w, err)
				return
			}
		}

		messeges, err := srv.repo.Get(m.ToMessage(), m.Offset, m.ValidateLimit())
		if err != nil {
			srv.handleError(w, err)
			return
		}

		resp = messeges

	case "POST":

		bytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			srv.handleError(w, err)
			return
		}

		if len(bytes) > 0 {
			var m presenter.CreateMessegeRequests
			err = json.Unmarshal(bytes, &m)
			if err != nil {
				srv.handleError(w, err)
				return
			}

			// non-multi-threaded option
			// err = srv.repo.Create(m.ToMessage())
			// if err != nil {
			// 	srv.handleError(w, err)
			// 	return
			// }
			// resp = &presenter.CreateMessegeResponse{Status: "created"}

			c := srv.mq.Query()
			for _, message := range m {
				c <- message.ToMessage()
			}

			resp = &presenter.CreateMessegeResponse{Status: "pushed"}

		} else {
			resp = &presenter.CreateMessegeResponse{Status: "can't parse body"}
		}

	default:
		resp = &presenter.CreateMessegeResponse{Status: "method not allowed"}
	}

	r, _ := json.Marshal(resp)
	w.Header().Add("Content-Type", "application/json")

	w.Write(r)

}

func (srv *service) clear(w http.ResponseWriter, req *http.Request) {

	var resp interface{}

	switch req.Method {
	case "POST":
		srv.mq.Clear()
		resp = &presenter.CreateMessegeResponse{Status: "cleared"}
	default:
		resp = &presenter.CreateMessegeResponse{Status: "method not allowed"}
	}

	r, _ := json.Marshal(resp)
	w.Write(r)
	w.Header().Add("Content-Type", "application/json")
}

func (srv *service) handleError(w http.ResponseWriter, err error) {
	r, _ := json.Marshal(&presenter.CreateMessegeResponse{Status: err.Error()})
	w.Write(r)
	w.Header().Add("Content-Type", "application/json")
}
