package service

import (
	"encoding/json"
	"net/http"
	"tender/models/tender"
	repeatable "tender/pkg/utils"

	"github.com/gorilla/mux"
)

func (s *Service) PostTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var tenderDTO tender.CreateTenderDTO
	err := json.NewDecoder(r.Body).Decode(&tenderDTO)
	if err != nil {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
		w.Write(bytes)
		return
	}
	if !repeatable.IsValidUUID(tenderDTO.OrganizationId) {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
		w.Write(bytes)
		return
	}
	if !s.checkUserInDB(ctx, tenderDTO.CreatorUsername) {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	validUser, err := s.CheckUserAuth(ctx, tenderDTO.CreatorUsername, tenderDTO.OrganizationId)
	if !validUser {
		w.WriteHeader(403)
		bytes, _ := json.Marshal(Reason{Reason: "Недостаточно прав для выполнения действия."})
		w.Write(bytes)
		return
	}
	if err != nil {
		w.WriteHeader(403)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	createdTender := tenderDTO.ToTender()
	s.tenderRepo.Create(ctx, &createdTender)
	bytes, _ := json.Marshal(createdTender)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Service) GetTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceTypes := r.URL.Query()["service_type"]
	for _, seserviceType := range serviceTypes {
		if !s.constants.ValidTenderServiceType[seserviceType] {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
	}
	s.qm.ServiceTypes = serviceTypes
	tenders, err := s.tenderRepo.FindAll(ctx, &s.qm)
	if err != nil {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	tendersDTOs := tender.ToTenderDTOs(tenders)
	bytes, _ := json.Marshal(tendersDTOs)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Service) GetMyTenders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userName := r.URL.Query().Get("username")
	if s.checkUserInDB(ctx, userName) {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	user, err := s.userRepo.FindOne(ctx, userName)
	if user.Id == "" {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	if err != nil {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	resp, _ := s.responsibleRepo.FindAll(ctx, user.Id)
	if len(resp) == 0 {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	s.qm.Organization_id = make([]string, len(resp))
	for i, v := range resp {
		s.qm.Organization_id[i] = v.OrganizationId
	}
	tenders, err := s.tenderRepo.FindAll(ctx, &s.qm)
	if err != nil {
		w.WriteHeader(403)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	tendersDTOs := tender.ToTenderDTOs(tenders)
	bytes, _ := json.Marshal(tendersDTOs)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Service) GetTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userName := r.URL.Query().Get("username")
	tenderId := mux.Vars(r)["tenderId"]
	tender, err := s.tenderRepo.FindOne(ctx, tenderId)
	if err != nil {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	validUser, err := s.CheckUserAuth(ctx, userName, tender.OrganizationId)
	if !validUser {
		w.WriteHeader(403)
		bytes, _ := json.Marshal(Reason{Reason: "Недостаточно прав для выполнения действия."})
		w.Write(bytes)
		return
	}
	if err != nil {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tender.Status))
}

func (s *Service) ChangeTenderStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderId := mux.Vars(r)["tenderId"]
	userName := r.URL.Query().Get("username")
	status := r.URL.Query().Get("status")
	if !s.constants.ValidStatus[status] {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
		w.Write(bytes)
		return
	}
	tender, err := s.tenderRepo.FindOne(ctx, tenderId)
	if err != nil {
		w.WriteHeader(400)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}
	validUser, _ := s.CheckUserAuth(ctx, userName, tender.OrganizationId)
	if !validUser {
		w.WriteHeader(403)
		bytes, _ := json.Marshal(Reason{Reason: "Недостаточно прав для выполнения действия."})
		w.Write(bytes)
		return
	}
	s.tenderRepo.Update(ctx, &tender, map[string]interface{}{"status": status})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tender.Status))
}

func (s *Service) EditTender(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var tenderDTO tender.UpdateTenderDTO
	err := json.NewDecoder(r.Body).Decode(&tenderDTO)
	if err != nil {
		panic(err)
	}
	var jsonData []byte
	r.Body.Read(jsonData)
	var fields map[string]interface{}
	json.Unmarshal(jsonData, &fields)
	userName := r.URL.Query().Get("username")
	tenderId := mux.Vars(r)["tenderId"]

	validUser, err := s.CheckUserAuth(ctx, userName, tenderId)
	if !validUser {
		panic("User invalid")
	}
	if err != nil {
		panic(err)
	}

	updatedTender := tenderDTO.ToTender()
	// s.tenderRepo.Update(ctx, &updatedTender)
	bytes, _ := json.Marshal(updatedTender)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
