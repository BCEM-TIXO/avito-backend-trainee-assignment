package service

import (
	"context"
	"encoding/json"
	"net/http"
	"tender/models/bid"
	repeatable "tender/pkg/utils"

	// repeatable "tender/pkg/utils"

	"github.com/gorilla/mux"
)

func (s *Service) checkAuthorId(ctx context.Context, authorType string, id string) bool {
	if !repeatable.IsValidUUID(id) {
		return false
	}
	if authorType == "User" {
		_, err := s.userRepo.FindOneId(ctx, id)
		if err != nil {
			return false
		}
	}
	if authorType == "Organization" {
		_, err := s.organizationRepo.FindOne(ctx, id)
		if err != nil {
			return false
		}
	}
	return true
}

func (s *Service) GetBids(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenderId := mux.Vars(r)["tenderId"]
	userName := r.URL.Query().Get("username")
	if !s.checkUserInDB(ctx, userName) {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	bids, _ := s.bidRepo.FindAll(ctx, tenderId, &s.qm)
	bidsDTOs := bid.ToBidDTOs(bids)
	bytes, _ := json.Marshal(bidsDTOs)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (s *Service) PostBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var CreateBidDTO bid.CreateBidDTO
	json.NewDecoder(r.Body).Decode(&CreateBidDTO)
	b := CreateBidDTO.ToBid()

	if !s.checkAuthorId(ctx, b.AuthorType, b.AuthorId) {
		w.WriteHeader(401)
		bytes, _ := json.Marshal(Reason{Reason: "Пользователь не существует или некорректен."})
		w.Write(bytes)
		return
	}
	err := s.bidRepo.Create(ctx, &b)
	if err != nil {
		w.WriteHeader(404)
		bytes, _ := json.Marshal(Reason{Reason: err.Error()})
		w.Write(bytes)
		return
	}

	bidDTO := b.ToBidDTO()
	bytes, _ := json.Marshal(bidDTO)
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}
