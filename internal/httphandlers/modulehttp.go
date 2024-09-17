package httphandlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

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
	mux.HandleFunc(fmt.Sprintf("GET /%s/{id}", basePath), mh.GetModuleByIdHandler)
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
	moduleDetails, err := mh.ms.GetAllModules()
	if err != nil {
		log.Println("Error in fetching Modules", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch Module Details")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, moduleDetails)
}

func (mh *ModuleHttpHandler) GetModuleByIdHandler(w http.ResponseWriter, r *http.Request) {

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		log.Println("Invalid URL format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	moduleID := pathSegments[len(pathSegments)-1]
	if moduleID == "" {
		log.Println("Missing module ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing module ID")
		return
	}

	if !isValidModuleID(moduleID) {
		log.Println("Invalid module ID format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	getModule, err := mh.ms.GetSpecificModule(moduleID)
	if err != nil {
		log.Println("Error in fetching Modules", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch Module ")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, getModule)
}

// isValidModuleID checks if the moduleID is a valid hexadecimal string of length 24
func isValidModuleID(id string) bool {
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}
