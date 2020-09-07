package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/suhas1294/go-mock/config"
	"github.com/suhas1294/go-mock/models"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"path"
	"regexp"
	"strings"
)

type mockController struct{}
type Message struct {
	Msg string
}
var MockController mockController


func (m mockController) CreateMock(w http.ResponseWriter, r *http.Request) {
	mock := models.Mock{}
	json.NewDecoder(r.Body).Decode(&mock)
	mock.Id = bson.NewObjectId()

	// check if exact mock exists
	if exist, id := checkIfMockAlreadyExist(r); exist {
		bs, _ := json.Marshal(Message{"Mock already exist with, id: " + id})
		w.Write(bs)
		return
	}

	// check if end point if not proper
	if strings.HasSuffix(mock.Endpoint, "/") {
		na, _ := json.Marshal(
			struct {
				Message string
			}{"trailing '/' not allowed in end point, it can have query params/hash"})
		w.Write(na)
		fmt.Fprint(w)
		return
	}
	// insert into db
	config.MongoSession.DB("devdb").C("mocks").Insert(mock)
	mj, err := json.Marshal(mock)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	fmt.Fprintf(w, "%s\n", mj)
}

func (m mockController) GetMock(w http.ResponseWriter, r *http.Request) {
	if vm, _ := verifyMethod(r,w, http.MethodGet); !vm{
		return
	}
	id := path.Base(r.URL.Path)
	fmt.Println("get request for id", id)
	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound) // 404
		return
	}
	oid := bson.ObjectIdHex(id)
	mock := models.Mock{}
	if err := config.MongoSession.DB("devdb").C("mocks").FindId(oid).One(&mock); err != nil {
		w.WriteHeader(404)
		return
	}
	mj, err := json.Marshal(mock)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(w, "%s\n", mj)
}

func (m mockController) Mock(w http.ResponseWriter, r *http.Request) {

	aq := r.URL.Query()

	keys := make([]string, len(aq))
	i := 0
	for k := range aq {
		keys[i] = k
		i++
	}
	fmt.Println(keys)


	fmt.Println()
	keys, ok := r.URL.Query()["key1"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	pr := r.URL.Path
	fp := regexp.MustCompile(`\/mock\/(.*)`)
	paths := fp.FindStringSubmatch(pr)
	ep := paths[1]

	var a []models.Mock
	conditions := bson.M{
		"endpoint": ep,
	}
	config.MongoSession.DB("devdb").C("mocks").Find(conditions).All(&a)
	fmt.Println(ep)
	bs, _ := json.MarshalIndent(a[0].Payload, "", "\t")
	w.Write(bs)
	fmt.Fprint(w)
}

func verifyMethod(r *http.Request, w http.ResponseWriter, em string) (bool, http.ResponseWriter){
	if r.Method != em {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		na, _ := json.Marshal(
			struct {
				Message string
			}{"Method Not allowed"})
		w.Write(na)
		return false, w
	}
	return true, w
}

func checkIfMockAlreadyExist(r *http.Request) (bool, string) {
	mock := models.Mock{}
	err := json.NewDecoder(r.Body).Decode(&mock)

	if err != nil {
		panic(err)
	}
	var ml []models.Mock

	conditions := bson.M{
		"endpoint": mock.Endpoint,
		"queryparams": mock.QueryParams,
		"headers": mock.Headers,
	}
	config.MongoSession.DB("devdb").C("mocks").Find(conditions).All(&ml)
	if len(ml) > 0 {
		return false, ml[0].Id.Hex()
	}
	return true, ""
}