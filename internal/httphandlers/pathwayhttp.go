package httphandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"c2c.in/api/internal/models"
	"c2c.in/api/internal/services"
	"c2c.in/api/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)



type PathwayHttpHandler struct{
	ps *services.PathwayService
}

func NewPathwayHttpHandler(db *mongo.Database)* PathwayHttpHandler{
	ps:= services.NewPathwayService(db)
	return &PathwayHttpHandler{ps:ps}
}

func (ph *PathwayHttpHandler) RegisterServiceWithMux(mux *http.ServeMux){
	basePath := "/pathways"
	mux.HandleFunc(basePath,ph.CreatePathwayHandler)
}

func(ph *PathwayHttpHandler) CreatePathwayHandler(w http.ResponseWriter,r *http.Request){
	// Parse the JSON request body into a models.Pathway struct
	var pathway models.Pathway
	err:=json.NewDecoder(r.Body).Decode(&pathway)
	if err !=nil{
		log.Println("Error in Decoding Json",err)
		utils.RespondWithError(w,http.StatusBadRequest,"Invalid Request Body")
		return 
	}
	// Call the service to create the pathway in the database
	id,err:= ph.ps.CreatePathway(&pathway)
	if err!=nil{
		log.Println("Error in Creating Pathway",err)
		utils.RespondWithError(w,http.StatusInternalServerError,"Failed to Create the Pathway")
		return 
	}

	reponse := map[string]string{"id": id}
	utils.ResponseWithJson(w,http.StatusOK,reponse)
}