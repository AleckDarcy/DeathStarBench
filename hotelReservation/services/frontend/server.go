package frontend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/delimitrou/DeathStarBench/hotelreservation/dialer"
	"github.com/delimitrou/DeathStarBench/hotelreservation/registry"
	profile "github.com/delimitrou/DeathStarBench/hotelreservation/services/profile/proto"
	recommendation "github.com/delimitrou/DeathStarBench/hotelreservation/services/recommendation/proto"
	reservation "github.com/delimitrou/DeathStarBench/hotelreservation/services/reservation/proto"
	search "github.com/delimitrou/DeathStarBench/hotelreservation/services/search/proto"
	user "github.com/delimitrou/DeathStarBench/hotelreservation/services/user/proto"
	"github.com/delimitrou/DeathStarBench/hotelreservation/tls"
	"github.com/delimitrou/DeathStarBench/hotelreservation/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/AleckDarcy/ContextBus"
	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
	cb_http "github.com/AleckDarcy/ContextBus/third-party/go/net/http"
	"github.com/delimitrou/DeathStarBench/hotelreservation/services/context_bus"
)

// cheat go import
var _ = cb_http.DefaultServeMux
var _ = tracing.TracedServeMux{}

// Server implements frontend service
type Server struct {
	searchClient         search.SearchClient
	profileClient        profile.ProfileClient
	recommendationClient recommendation.RecommendationClient
	userClient           user.UserClient
	reservationClient    reservation.ReservationClient
	KnativeDns           string
	IpAddr               string
	Port                 int
	Tracer               opentracing.Tracer
	Registry             *registry.Client

	CBConfig *cb_configure.ServerConfigure
}

// Run the server
func (s *Server) Run() error {
	if s.Port == 0 {
		return fmt.Errorf("Server port must be set")
	}

	log.Info().Msg("Initializing gRPC clients...")
	if err := s.initSearchClient("srv-search"); err != nil {
		return err
	}

	if err := s.initProfileClient("srv-profile"); err != nil {
		return err
	}

	if err := s.initRecommendationClient("srv-recommendation"); err != nil {
		return err
	}

	if err := s.initUserClient("srv-user"); err != nil {
		return err
	}

	if err := s.initReservation("srv-reservation"); err != nil {
		return err
	}
	log.Info().Msg("Successfull")

	log.Trace().Msg("frontend before mux")

	context_bus.Set(s.CBConfig, context_bus.SetConfigureForTesting)
	mux := cb_http.NewServeMux() // context bus server mux
	//mux := tracing.NewServeMux(s.Tracer)
	mux.Handle("/", http.FileServer(http.Dir("services/frontend/static")))
	mux.Handle("/reset", cb_http.NewHandlerFunc(s.resetHandler))
	mux.Handle("/hotels", cb_http.NewHandlerFunc(s.searchHandler))
	mux.Handle("/recommendations", cb_http.NewHandlerFunc(s.recommendHandler))
	mux.Handle("/user", cb_http.NewHandlerFunc(s.userHandler))
	mux.Handle("/reservation", cb_http.NewHandlerFunc(s.reservationHandler))

	log.Trace().Msg("frontend starts serving")

	tlsconfig := tls.GetHttpsOpt()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: mux,
	}
	if tlsconfig != nil {
		log.Info().Msg("Serving https")
		srv.TLSConfig = tlsconfig
		return srv.ListenAndServeTLS("x509/server_cert.pem", "x509/server_key.pem")
	} else {
		log.Info().Msg("Serving http")
		return srv.ListenAndServe()
	}
}

func (s *Server) initSearchClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.searchClient = search.NewSearchClient(conn)
	return nil
}

func (s *Server) initProfileClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.profileClient = profile.NewProfileClient(conn)
	return nil
}

func (s *Server) initRecommendationClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.recommendationClient = recommendation.NewRecommendationClient(conn)
	return nil
}

func (s *Server) initUserClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.userClient = user.NewUserClient(conn)
	return nil
}

func (s *Server) initReservation(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.reservationClient = reservation.NewReservationClient(conn)
	return nil
}

func (s *Server) getGprcConn(name string) (*grpc.ClientConn, error) {
	log.Info().Msg("get Grpc conn is :")
	log.Info().Msg(s.KnativeDns)
	log.Info().Msg(fmt.Sprintf("%s.%s", name, s.KnativeDns))
	if s.KnativeDns != "" {
		return dialer.Dial(
			fmt.Sprintf("%s.%s", name, s.KnativeDns),
			dialer.WithTracer(s.Tracer))
	} else {
		return dialer.Dial(
			name,
			dialer.WithTracer(s.Tracer),
			dialer.WithBalancer(s.Registry.Client),
		)
	}
}

func (s *Server) resetHandler(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("reset databases starts")
	ctx := r.Context()
	s.searchClient.ResetDB(ctx, &search.NearbyRequest{})
	s.profileClient.ResetDB(ctx, &profile.Request{})
	s.recommendationClient.ResetDB(ctx, &recommendation.Request{})
	s.reservationClient.ResetDB(ctx, &reservation.Request{})
	s.userClient.ResetDB(ctx, &user.Request{})

	log.Info().Msg("reset databases ends")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "ok"})
}

func (s *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	// Context Bus
	cbCtx, cbOK := ContextBus.FromContext(r.Context())

	// Context Bus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.1",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "starts searchHandler",
			Paths:   nil,
		})
	} else {
		log.Info().Msg("starts searchHandler")
	}

	params := r.URL.Query()

	// in/out dates from query params
	inDate, outDate := params.Get("inDate"), params.Get("outDate")
	if inDate == "" || outDate == "" {
		http.Error(w, "Please specify inDate/outDate params", http.StatusBadRequest)
		return
	}

	// lan/lon from query params
	sLat, sLon := params.Get("lat"), params.Get("lon")
	if sLat == "" || sLon == "" {
		http.Error(w, "Please specify location params", http.StatusBadRequest)
		return
	}

	Lat, _ := strconv.ParseFloat(sLat, 32)
	lat := float32(Lat)
	Lon, _ := strconv.ParseFloat(sLon, 32)
	lon := float32(Lon)

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.2",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "starts searchHandler querying downstream",
			Paths:   nil,
		})
	} else {
		log.Info().Msg("starts searchHandler querying downstream")
	}

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.3",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: fmt.Sprintf("SEARCH [lat: %v, lon: %v, inDate: %v, outDate: %v", lat, lon, inDate, outDate),
			Paths:   nil,
		})
	} else {
		log.Info().Msgf("SEARCH [lat: %v, lon: %v, inDate: %v, outDate: %v", lat, lon, inDate, outDate)
	}

	// search for best hotels
	searchResp, err := s.searchClient.Nearby(ctx, &search.NearbyRequest{
		Lat:       lat,
		Lon:       lon,
		InDate:    inDate,
		OutDate:   outDate,
		CBPayload: cbCtx.Payload(), // set ContextBus payload
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.4",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "SearchHandler gets searchResp",
			Paths:   nil,
		})
	} else {
		log.Info().Msg("SearchHandler gets searchResp")
	}

	//for _, hid := range searchResp.HotelIds {
	//	log.Trace().Msgf("Search Handler hotelId = %s", hid)
	//}

	// grab locale from query params or default to en
	locale := params.Get("locale")
	if locale == "" {
		locale = "en"
	}

	reservationResp, err := s.reservationClient.CheckAvailability(ctx, &reservation.Request{
		CustomerName: "",
		HotelId:      searchResp.HotelIds,
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   1,
		CBPayload:    cbCtx.Payload(), // set ContextBus payload
	})
	if err != nil {
		log.Error().Msg("SearchHandler CheckAvailability failed: " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.5",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: fmt.Sprintf("searchHandler gets reserveResp.HotelId = %s", reservationResp.HotelId),
			Paths:   nil,
		})
	} else {
		log.Info().Msgf("searchHandler gets reserveResp.HotelId = %s", reservationResp.HotelId)
	}

	// hotel profiles
	profileResp, err := s.profileClient.GetProfiles(ctx, &profile.Request{
		HotelIds:  reservationResp.HotelId,
		Locale:    locale,
		CBPayload: cbCtx.Payload(), // set ContextBus payload
	})
	if err != nil {
		log.Error().Msg("SearchHandler GetProfiles failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "frontend.searchHandler.6",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "searchHandler gets profileResp",
			Paths:   nil,
		})
	} else {
		log.Info().Msg("searchHandler gets profileResp")
	}

	json.NewEncoder(w).Encode(geoJSONResponse(profileResp.Hotels))
}

func (s *Server) recommendHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	sLat, sLon := r.URL.Query().Get("lat"), r.URL.Query().Get("lon")
	if sLat == "" || sLon == "" {
		http.Error(w, "Please specify location params", http.StatusBadRequest)
		return
	}
	Lat, _ := strconv.ParseFloat(sLat, 64)
	lat := float64(Lat)
	Lon, _ := strconv.ParseFloat(sLon, 64)
	lon := float64(Lon)

	require := r.URL.Query().Get("require")
	if require != "dis" && require != "rate" && require != "price" {
		http.Error(w, "Please specify require params", http.StatusBadRequest)
		return
	}

	// recommend hotels
	recResp, err := s.recommendationClient.GetRecommendations(ctx, &recommendation.Request{
		Require: require,
		Lat:     float64(lat),
		Lon:     float64(lon),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// grab locale from query params or default to en
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "en"
	}

	// hotel profiles
	profileResp, err := s.profileClient.GetProfiles(ctx, &profile.Request{
		HotelIds: recResp.HotelIds,
		Locale:   locale,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(geoJSONResponse(profileResp.Hotels))
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	username, password := r.URL.Query().Get("username"), r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "Please specify username and password", http.StatusBadRequest)
		return
	}

	// Check username and password
	recResp, err := s.userClient.CheckUser(ctx, &user.Request{
		Username: username,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := "Login successfully!"
	if recResp.Correct == false {
		str = "Failed. Please check your username and password. "
	}

	res := map[string]interface{}{
		"message": str,
	}

	json.NewEncoder(w).Encode(res)
}

func (s *Server) reservationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx := r.Context()

	inDate, outDate := r.URL.Query().Get("inDate"), r.URL.Query().Get("outDate")
	if inDate == "" || outDate == "" {
		http.Error(w, "Please specify inDate/outDate params", http.StatusBadRequest)
		return
	}

	if !checkDataFormat(inDate) || !checkDataFormat(outDate) {
		http.Error(w, "Please check inDate/outDate format (YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	hotelId := r.URL.Query().Get("hotelId")
	if hotelId == "" {
		http.Error(w, "Please specify hotelId params", http.StatusBadRequest)
		return
	}

	customerName := r.URL.Query().Get("customerName")
	if customerName == "" {
		http.Error(w, "Please specify customerName params", http.StatusBadRequest)
		return
	}

	username, password := r.URL.Query().Get("username"), r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "Please specify username and password", http.StatusBadRequest)
		return
	}

	numberOfRoom := 0
	num := r.URL.Query().Get("number")
	if num != "" {
		numberOfRoom, _ = strconv.Atoi(num)
	}

	// Check username and password
	recResp, err := s.userClient.CheckUser(ctx, &user.Request{
		Username: username,
		Password: password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	str := "Reserve successfully!"
	if recResp.Correct == false {
		str = "Failed. Please check your username and password. "
	}

	// Make reservation
	resResp, err := s.reservationClient.MakeReservation(ctx, &reservation.Request{
		CustomerName: customerName,
		HotelId:      []string{hotelId},
		InDate:       inDate,
		OutDate:      outDate,
		RoomNumber:   int32(numberOfRoom),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(resResp.HotelId) == 0 {
		str = "Failed. Already reserved. "
	}

	res := map[string]interface{}{
		"message": str,
	}

	json.NewEncoder(w).Encode(res)
}

// return a geoJSON response that allows google map to plot points directly on map
// https://developers.google.com/maps/documentation/javascript/datalayer#sample_geojson
func geoJSONResponse(hs []*profile.Hotel) map[string]interface{} {
	fs := []interface{}{}

	for _, h := range hs {
		fs = append(fs, map[string]interface{}{
			"type": "Feature",
			"id":   h.Id,
			"properties": map[string]string{
				"name":         h.Name,
				"phone_number": h.PhoneNumber,
			},
			"geometry": map[string]interface{}{
				"type": "Point",
				"coordinates": []float32{
					h.Address.Lon,
					h.Address.Lat,
				},
			},
		})
	}

	return map[string]interface{}{
		"type":     "FeatureCollection",
		"features": fs,
	}
}

func checkDataFormat(date string) bool {
	if len(date) != 10 {
		return false
	}
	for i := 0; i < 10; i++ {
		if i == 4 || i == 7 {
			if date[i] != '-' {
				return false
			}
		} else {
			if date[i] < '0' || date[i] > '9' {
				return false
			}
		}
	}
	return true
}
