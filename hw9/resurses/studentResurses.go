package resursestype

import "C"
import (
	"GoLangProjector/hw6/hw9/converter"
	"GoLangProjector/hw6/hw9/entity"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type StudentResource struct {
	S *entity.Storage
	C converter.Converter
}

func (sR *StudentResource) GetUserById(w http.ResponseWriter, r *http.Request) {
	idUser := r.PathValue("id")
	id, err := strconv.Atoi(idUser)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	student := sR.S.GetStudent(id)
	dtoStudent := sR.C.Convert(student)
	err = json.NewEncoder(w).Encode(dtoStudent)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
