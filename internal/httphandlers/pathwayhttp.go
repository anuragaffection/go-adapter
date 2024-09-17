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

type PathwayHttpHandler struct {
	ps *services.PathwayService
}

func NewPathwayHttpHandler(db *mongo.Database) *PathwayHttpHandler {
	ps := services.NewPathwayService(db)
	return &PathwayHttpHandler{ps: ps}
}

func (ph *PathwayHttpHandler) RegisterServiceWithMux(mux *http.ServeMux) {
	basePath := "pathways"
	mux.HandleFunc(fmt.Sprintf("POST /%s", basePath), ph.CreatePathwayHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s", basePath),ph.GetAllPathwayHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s/{id}", basePath), ph.GetPathwayByIdHandler)
}

func (ph *PathwayHttpHandler) CreatePathwayHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body into a models.Pathway struct
	var pathway models.Pathway
	err := json.NewDecoder(r.Body).Decode(&pathway)
	if err != nil {
		log.Println("Error in Decoding Json", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}
	// Call the service to create the pathway in the database
	id, err := ph.ps.CreatePathway(&pathway)
	if err != nil {
		log.Println("Error in Creating Pathway", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to Create the Pathway")
		return
	}

	reponse := map[string]string{"id": id}
	utils.ResponseWithJson(w, http.StatusOK, reponse)
}

func (ph *PathwayHttpHandler) GetAllPathwayHandler(w http.ResponseWriter, r *http.Request) {
	pathwayNames, err := ph.ps.GetAllPathway()
	if err != nil {
		log.Println("Error in Fetching Pathways", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to Fetch Pathway Names")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, pathwayNames)
}

func (ph *PathwayHttpHandler) GetPathwayByIdHandler(w http.ResponseWriter,r *http.Request){

	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		log.Println("Invalid URL format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	pathwayID := pathSegments[len(pathSegments)-1]
	if pathwayID == "" {
		log.Println("Missing pathway ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing pathway ID")
		return
	}

	if !isValidPathwayID(pathwayID) {
		log.Println("Invalid module ID format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid module ID format")
		return
	}

	getPathway, err := ph.ps.GetSpecificPathway(pathwayID)
	if err != nil {
		log.Println("Error in fetching Modules", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch Module ")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, getPathway)
}

// isValiPathwayID checks if the moduleID is a valid hexadecimal string of length 24
func isValidPathwayID(id string) bool {
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}