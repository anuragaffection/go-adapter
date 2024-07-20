package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"c2c.in/api/internal/models"
	"c2c.in/api/internal/services"
	"c2c.in/api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type TopicHttpHandler struct {
	ts *services.TopicService
}

func NewTopicHttpHandler(db *mongo.Database) *TopicHttpHandler {
	ts := services.NewTopicService(db)

	return &TopicHttpHandler{ts: ts}
}

func (th *TopicHttpHandler) RegisterServiceWithMux(mux *http.ServeMux) {
	basePath := "topics"
	mux.HandleFunc(fmt.Sprintf("POST /%s",basePath), th.CreateTopicHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s",basePath),th.GetAllTopicHandler)
}

func (th *TopicHttpHandler) CreateTopicHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the JSON request body into a models.Topic struct
	var topic models.Topic
	err:= json.NewDecoder(r.Body).Decode(&topic)
	if err!=nil{
		log.Println("Error in Decoding Json",err)
		utils.RespondWithError(w,http.StatusBadRequest,"Invalid Request Body")
		return 
	}
	id,err:= th.ts.CreatTopic(&topic)
	if err!=nil{
		log.Println("Error in creating Topic :",err)
		utils.RespondWithError(w,http.StatusInternalServerError,"Failed to Create Topic")
		return 
	}
	response := map[string]string{"id":id}
	utils.ResponseWithJson(w,http.StatusOK,response)

}

func(th *TopicHttpHandler) GetAllTopicHandler(w http.ResponseWriter,r *http.Request){
	topicNames,err:=th.ts.GetAllTopics()
	if err!=nil{
		log.Println("Error in Fetching Topics",err)
		utils.RespondWithError(w,http.StatusBadRequest,"Failed to Fetch Topics Names")
		return 
	}

	utils.ResponseWithJson(w,http.StatusOK,topicNames)
}
