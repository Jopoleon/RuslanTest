package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type FileServer struct {
	s StorageOperator
}

func (fs *FileServer) UploadFile(w http.ResponseWriter, r *http.Request) {
	logrus.Info("UploadFile request with content-length: ", r.Header.Get("Content-Length"))
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Printf("%+v \n", errors.WithStack(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("%+v \n", errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = fs.s.Put(handler.Filename, b)
	if err != nil {
		fmt.Printf("%+v \n", errors.WithStack(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SUCCESS"))
}

//func (server *AlertServer) AlertHTTP(w http.ResponseWriter, r *http.Request) {
//	logrus.Info("Alert request with content-length: ", r.Header.Get("Content-Length"))
//
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		logrus.Errorf("AlertHTTP ioutil.ReadAll(req.Body) error: %v", err)
//		http.Error(w, fmt.Sprintf("AlertHTTP ioutil.ReadAll(req.Body) error: %v", err), http.StatusBadRequest)
//		return
//	}
//
//	defer r.Body.Close()
//	var al db.Alert
//
//	err = json.Unmarshal(body, &al)
//	if err != nil {
//		logrus.Errorf("AlertHTTP  json.Unmarshal(body, &alert) error: %v", err)
//		http.Error(w, fmt.Sprintf("AlertHTTP  json.Unmarshal(body, &alert) error: %v.", err), http.StatusInternalServerError)
//		return
//	}
//	alreadyRegistered, err := db.IsRegistered(al.SystemId, db.NewRegMnr(server.Mgr.Db))
//	//reg := db.NewRegMnr(server.Mgr.Db)
//	//alreadyRegistered, err := reg.IsExist(al.SystemId)
//	if err != nil {
//		logrus.Errorf("AlertHTTP db.IsRegistered error: %v", err)
//		http.Error(w, fmt.Sprintf("AlertHTTP db.IsRegistered error: %v.", err), http.StatusInternalServerError)
//		return
//	}
//	if alreadyRegistered {
//		err = server.Mgr.Put(al)
//		if err != nil {
//			logrus.Errorf("AlertHTTP db.Put(alert) error: %v", err)
//			http.Error(w, fmt.Sprintf("AlertHTTP db.Put(alert) error: %v.", err), http.StatusInternalServerError)
//			return
//		}
//		w.WriteHeader(http.StatusOK)
//	} else {
//		logrus.Errorf("AlertHTTP Unregistered system ID : %+v", al)
//		http.Error(w, fmt.Sprintf("Client with system ID = %v , is not registered", al.SystemId), http.StatusBadRequest)
//	}
//}
