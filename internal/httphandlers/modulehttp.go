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

type ModuleHttpHandler struct {
	ms *services.ModuleService
}

func NewModuleHttpHandler(db *mongo.Database) *ModuleHttpHandler {
	ms := services.NewModuleService(db)

	return &ModuleHttpHandler{ms: ms}
}

func (mh *ModuleHttpHandler) RegisterServiceWithMux(mux *http.ServeMux) {
	basePath := "modules"
	mux.HandleFunc(fmt.Sprintf("POST /%s", basePath), mh.CreateModuleHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s", basePath), mh.GetAllModuleHandler)
}

func (mh *ModuleHttpHandler) CreateModuleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a models.Module struct
	var module models.Module
	err := json.NewDecoder(r.Body).Decode(&module)
	if err != nil {
		log.Println("Error in Decoding Json", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	id, err := mh.ms.CreateModule(&module)
	if err != nil {
		log.Println("Error in Creating Module", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to Create Module")
		return
	}
	response := map[string]string{"id": id}
	utils.ResponseWithJson(w, http.StatusOK, response)
}

func (mh *ModuleHttpHandler) GetAllModuleHandler(w http.ResponseWriter, r *http.Request) {
	moduleNames, err := mh.ms.GetAllModules()
	if err != nil {
		log.Println("Error in fetching Modules", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch Module Names")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, moduleNames)
}
