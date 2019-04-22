package train

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"team-project/database"
	"team-project/services/common"
	"team-project/services/data"
	"team-project/services/model"

	"github.com/go-zoo/bone"
)

var emptyResponse interface{}

func validateIfEmpty(t data.Train) error {
	if t.DepartureCity == "" || t.ArrivalCity == "" || t.DepartureTime == "" || t.DepartureDate == "" || t.ArrivalTime == "" || t.ArrivalDate == "" {
		return errors.New("Some incoming data is empty =(")
	}
	return nil
}

func nameIsValid(str string) bool {
	var validName = regexp.MustCompile(`^[a-zA-Z]+(?:[\s-][a-zA-Z]+)*$`)
	return validName.MatchString(str)
}

func timeIsValid(str string) bool {
	var validName = regexp.MustCompile(`^(2[0-3]|[01]?[0-9]):([0-5]?[0-9])$`)
	return validName.MatchString(str)
}

func dateIsValid(str string) bool {
	var validName = regexp.MustCompile(`^\s*(3[01]|[12][0-9]|0?[1-9])\.(1[012]|0?[1-9])\.((?:19|20)\d{2})\s*$`)
	return validName.MatchString(str)
}

func validate(t data.Train) error {
	if nameIsValid(t.ArrivalCity) == false || nameIsValid(t.DepartureCity) == false || dateIsValid(t.DepartureDate) == false || dateIsValid(t.ArrivalDate) == false || timeIsValid(t.ArrivalTime) == false || timeIsValid(t.DepartureTime) == false {
		return errors.New("Some data is invalid")
	}
	return nil
}

//GetTrains is a method
func GetTrains(w http.ResponseWriter, r *http.Request) {
	trains, err := database.Trains.GetAllTrains()
	if err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}

	/*for _, train := range trains {
		if err = validateIfEmpty(train); err != nil {
			common.RenderJSON(w, r, 404, emptyResponse)
			return
		}
		if err = validate(train); err != nil {
			common.RenderJSON(w, r, 404, emptyResponse)
			return
		}
	}*/

	common.RenderJSON(w, r, http.StatusOK, trains)
}

//GetSingleTrain is a method
func GetSingleTrain(w http.ResponseWriter, r *http.Request) {
	params := bone.GetAllValues(r)
	train, err := database.Trains.GetTrain(params["id"])
	if err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	if err = validateIfEmpty(train); err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	if err = validate(train); err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	common.RenderJSON(w, r, http.StatusOK, train)
}

//CreateTrain is a method
func CreateTrain(w http.ResponseWriter, r *http.Request) {
	t := data.Train{}
	json.NewDecoder(r.Body).Decode(&t)
	if err := validateIfEmpty(t); err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	if err := validate(t); err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	err := database.Trains.AddTrain(t)
	if err != nil {
		common.RenderJSON(w, r, http.StatusNoContent, emptyResponse)
		return
	}
	common.RenderJSON(w, r, http.StatusOK, t)
}

//UpdateTrain is a method
func UpdateTrain(w http.ResponseWriter, r *http.Request) {
	id, err := model.GetID(r)
	if err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	t := data.Train{}
	json.NewDecoder(r.Body).Decode(&t)
	if err = validateIfEmpty(t); err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	if err = validate(t); err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	t.ID = id
	err = database.Trains.UpdateTrain(t.ID, t.DepartureCity, t.ArrivalCity, t.DepartureDate, t.DepartureTime, t.ArrivalTime, t.ArrivalDate)
	if err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	train, err := database.Trains.GetTrain(t.ID.String())
	if err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	common.RenderJSON(w, r, 202, train)
}

//DeleteTrain is a method
func DeleteTrain(w http.ResponseWriter, r *http.Request) {
	id, err := model.GetID(r)
	if err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	err = database.Trains.DeleteTrain(id)
	if err != nil {
		common.RenderJSON(w, r, 404, emptyResponse)
		return
	}
	common.RenderJSON(w, r, 202, emptyResponse)
}
