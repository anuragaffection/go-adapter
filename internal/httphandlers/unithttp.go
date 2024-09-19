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

type UnitHttpHandler struct {
	us *services.UnitService
}

func NewUnitHttpHandler(db *mongo.Database) *UnitHttpHandler {
	us := services.NewUnitService(db)

	return &UnitHttpHandler{us: us}
}

func (uh *UnitHttpHandler) RegisterServiceWithMux(mux *http.ServeMux) {

	basePath := "units"
	mux.HandleFunc(fmt.Sprintf("POST /%s", basePath), uh.CreateUnitHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s", basePath), uh.GetAllUnitHandler)
	mux.HandleFunc(fmt.Sprintf("GET /%s/{id}",basePath),uh.GetUnitByIdHandler) 
}

func (uh *UnitHttpHandler) CreateUnitHandler(w http.ResponseWriter, r *http.Request) {

	// Parse the JSON request body into a models.Unit struct
	var unit models.Unit
	err := json.NewDecoder(r.Body).Decode(&unit)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid Request Body")
		return
	}

	log.Println(unit)
	// Call the service to create the unit in the database
	id, err := uh.us.CreatUnit(&unit)
	if err != nil {
		log.Println("Error in Creating the Unit :", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create unit")
		return
	}

	reponse := map[string]string{"id": id}
	utils.ResponseWithJson(w, http.StatusOK, reponse)
}

func (uh *UnitHttpHandler) GetAllUnitHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println("hi")
	unitNames, err := uh.us.GetAllUnits()
	if err != nil {
		log.Println("Error in Fetching Units :", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to Fetch Units Names")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, unitNames)
	
}

func(uh *UnitHttpHandler) GetUnitByIdHandler(w http.ResponseWriter,r *http.Request){
	pathSegments := strings.Split(r.URL.Path, "/")
	if len(pathSegments) < 3 {
		log.Println("Invalid URL format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL format")
		return
	}

	unitID := pathSegments[len(pathSegments)-1]
	if unitID == "" {
		log.Println("Missing unit ID")
		utils.RespondWithError(w, http.StatusBadRequest, "Missing unit ID")
		return
	}

	if !isValidUnitID(unitID) {
		log.Println("Invalid unit ID format")
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid unit ID format")
		return
	}
	getUnit, err := uh.us.GetSpecificUnit(unitID)
	if err != nil {
		log.Println("Error in fetching units", err)
		utils.RespondWithError(w, http.StatusBadGateway, "Failed to Fetch units ")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, getUnit)
}

// isValidUnitID checks if the moduleID is a valid hexadecimal string of length 24
func isValidUnitID(id string) bool {
	re := regexp.MustCompile(`^[a-fA-F0-9]{24}$`)
	return re.MatchString(id)
}