package resursestype

import (
	"GoLangProjector/hw9/converter"
	"GoLangProjector/hw9/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ClassResource struct {
	S  *entity.Storage
	Cc *converter.ClassConverter
}

func (cR *ClassResource) GetAllClasses(w http.ResponseWriter, r *http.Request) {
	classes := cR.S.GetAllClasses()
	dtoClasses := cR.Cc.Convert(classes)

	err := json.NewEncoder(w).Encode(dtoClasses)
	if err != nil {
		fmt.Printf("Failed to encode: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (cR *ClassResource) GetClass(w http.ResponseWriter, r *http.Request) {
	idClass := r.PathValue("id")
	id, err := strconv.Atoi(idClass)
	if err != nil {
		http.Error(w, "Invalid class ID", http.StatusBadRequest)
		return
	}

	class := cR.S.GetClass(id)
	err = json.NewEncoder(w).Encode(class)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
