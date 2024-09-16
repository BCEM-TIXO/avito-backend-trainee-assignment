package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	repeatable "tender/pkg/utils"

	"github.com/gorilla/mux"
)

func (s *Service) paginationMethodMW(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset := r.URL.Query().Get("offset")
		limit := r.URL.Query().Get("limit")
		if offset == "" {
			offset = "0"
		}
		if limit == "" {
			limit = "5"
		}
		var err error
		s.qm.Pagination.Limit, err = strconv.Atoi(limit)
		if err != nil || s.qm.Pagination.Limit < 0 {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
		s.qm.Pagination.Offset, err = strconv.Atoi(offset)
		if err != nil || s.qm.Pagination.Offset < 0 {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
		next(w, r)
	})
}

func setJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Service) checkTenderId(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		tenderId := mux.Vars(r)["tenderId"]
		if !repeatable.IsValidUUID(tenderId) {
			w.WriteHeader(400)
			bytes, _ := json.Marshal(Reason{Reason: "Неверный формат запроса или его параметры."})
			w.Write(bytes)
			return
		}
		_, err := s.tenderRepo.FindOne(ctx, tenderId)
		if err != nil {
			w.WriteHeader(404)
			bytes, _ := json.Marshal(Reason{Reason: "Тендер не найден."})
			w.Write(bytes)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Service) checkUserInDB(ctx context.Context, userName string) bool {
	_, err := s.userRepo.FindOne(ctx, userName)
	return err == nil
}

func (s *Service) CheckUserAuth(ctx context.Context, userName string, orgId string) (bool, error) {
	user, err := s.userRepo.FindOne(ctx, userName)
	fmt.Println(user.Id, orgId)
	if err != nil {
		panic(err)
		// return false, err
	}
	resp, err := s.responsibleRepo.FindOne(ctx, user.Id, orgId)
	if err != nil {
		panic(err)
		// return false, err
	}
	if resp.Id == "" {
		panic("ID empt")
		// return true, nil
	}
	return true, nil
}
