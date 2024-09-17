package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"regexp"

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
	mux.HandleFunc(fmt.Sprintf("GET /%s/{id}",basePath),th.GetTopicByIdHandler)
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

func(th *TopicHttpHandler) GetTopicByIdHandler(w http.ResponseWriter,r *http.Request){
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		log.Println("Invalid URL format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	topicID := pathSegments[len(pathSegments)-1]
	if topicID == "" {
		log.Println("Missing topic ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing topic ID")
		return
	}

	if !isValidTopicID(topicID) {
		log.Println("Invalid topic ID format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid topic ID format")
		return
	}
	getTopic, err := th.ts.GetSpecificTopic(topicID)
	if err != nil {
		log.Println("Error in fetching Topic", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch Topic ")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, getTopic)
}

// isValidTopicID checks if the moduleID is a valid hexadecimal string of length 24
func isValidTopicID(id string) bool {
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}