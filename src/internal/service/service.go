package service

import (
	"net/http"
	"tender/models/bid"
	bidDb "tender/models/bid/db"
	"tender/models/organization"
	organizationDb "tender/models/organization/db"
	"tender/models/responsible"
	responsibleDb "tender/models/responsible/db"
	"tender/models/tender"
	tenderDb "tender/models/tender/db"
	"tender/models/user"
	userDb "tender/models/user/db"
	postgresql "tender/pkg/client"

	"github.com/gorilla/mux"
)

const (
	ServiceType1 = "Construction"
	ServiceType2 = "Delivery"
	ServiceType3 = "Manufacture"

	Status1 = "Created"
	Status2 = "Published"
	Status3 = "Closed"
)

func getValidServiceTypesMap() map[string]bool {
	return map[string]bool{
		ServiceType1: true,
		ServiceType2: true,
		ServiceType3: true,
	}
}

func getValidStatusMap() map[string]bool {
	return map[string]bool{
		Status1: true,
		Status2: true,
		Status3: true,
	}
}

type Reason struct {
	Reason string `json:"reason"`
}

type Service struct {
	addres           string
	organizationRepo organization.Repository
	responsibleRepo  responsible.Repository
	tenderRepo       tender.Repository
	userRepo         user.Repository
	bidRepo          bid.Repository
	qm               tender.FindAllQueryModifier
	constants        constants
}

type constants struct {
	ValidStatus            map[string]bool
	ValidTenderServiceType map[string]bool
}

func NewService(psqlClient postgresql.Client, addres string) *Service {
	return &Service{
		addres:           addres,
		organizationRepo: organizationDb.NewRepository(psqlClient),
		responsibleRepo:  responsibleDb.NewRepository(psqlClient),
		tenderRepo:       tenderDb.NewRepository(psqlClient),
		userRepo:         userDb.NewRepository(psqlClient),
		bidRepo:          bidDb.NewRepository(psqlClient),
		constants: constants{
			ValidStatus:            getValidStatusMap(),
			ValidTenderServiceType: getValidServiceTypesMap(),
		},
	}
}

func (s Service) Run() {
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ping", s.Ping).Methods("GET")
	tendersRouter := apiRouter.PathPrefix("/tenders").Subrouter()
	tendersRouter.Use(setJSONContentType)
	tendersRouter.HandleFunc("", s.paginationMethodMW(s.GetTenders)).Methods("GET")

	tendersRouter.HandleFunc("/my", s.paginationMethodMW(s.GetMyTenders)).Methods("GET")
	tendersRouter.HandleFunc("/{tenderId}/status", s.checkTenderId(s.GetTenderStatus)).Methods("GET")
	tendersRouter.HandleFunc("/{tenderId}/status", s.checkTenderId(s.GetTenderStatus)).Methods("PUT")
	tendersRouter.HandleFunc("/{tenderId}/edit", s.checkTenderId(s.GetTenderStatus)).Methods("PATCH")
	tendersRouter.HandleFunc("/{tenderId}/rollback/{version}", s.checkTenderId(s.GetTenderStatus)).Methods("PUT")
	tendersRouter.HandleFunc("/new", s.PostTender).Methods("POST")

	bidsRouter := apiRouter.PathPrefix("/bids").Subrouter()

	bidsRouter.Use(setJSONContentType)

	bidsRouter.HandleFunc("/new", s.PostBid).Methods("POST")
	bidsRouter.HandleFunc("", s.paginationMethodMW(s.GetBids)).Methods("GET")

	r.Use()
	http.ListenAndServe(s.addres, r)
}
