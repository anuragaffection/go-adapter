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

	// log.Println(unit)
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

	unitNames, err := uh.us.GetAllUnits()
	if err != nil {
		log.Println("Error in Fetching Units :", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to Fetch Units Names")
		return
	}

	utils.ResponseWithJson(w, http.StatusOK, unitNames)

}
