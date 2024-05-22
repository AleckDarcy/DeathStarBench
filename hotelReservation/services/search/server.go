package search

import (
	// "encoding/json"
	"fmt"
	"github.com/AleckDarcy/ContextBus"
	cb "github.com/AleckDarcy/ContextBus/proto"

	// F"io/ioutil"
	"net"
	// "os"
	"time"

	"github.com/delimitrou/DeathStarBench/hotelreservation/dialer"
	"github.com/delimitrou/DeathStarBench/hotelreservation/registry"
	geo "github.com/delimitrou/DeathStarBench/hotelreservation/services/geo/proto"
	rate "github.com/delimitrou/DeathStarBench/hotelreservation/services/rate/proto"
	pb "github.com/delimitrou/DeathStarBench/hotelreservation/services/search/proto"
	"github.com/delimitrou/DeathStarBench/hotelreservation/tls"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	"github.com/delimitrou/DeathStarBench/hotelreservation/services/context_bus"
)

const name = "srv-search"

// Server implments the search service
type Server struct {
	geoClient  geo.GeoClient
	rateClient rate.RateClient

	Tracer     opentracing.Tracer
	Port       int
	IpAddr     string
	KnativeDns string
	Registry   *registry.Client
	uuid       string

	CBConfig *cb_configure.ServerConfigure
}

// Run starts the server
func (s *Server) Run() error {
	if s.Port == 0 {
		return fmt.Errorf("server port must be set")
	}

	s.uuid = uuid.New().String()

	opts := []grpc.ServerOption{
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: 120 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(s.Tracer),
		),
	}

	// ContextBus initialization
	context_bus.Set(s.CBConfig, context_bus.SetConfigureForTesting)
	// ContextBus disable opentracing
	if context_bus.CONTEXTBUS_ON {
		fmt.Println("ContextBus is on, disable opentracing interceptor")
		opts = opts[0 : len(opts)-1]
	}

	if tlsopt := tls.GetServerOpt(); tlsopt != nil {
		opts = append(opts, tlsopt)
	}

	srv := grpc.NewServer(opts...)
	pb.RegisterSearchServer(srv, s)

	// init grpc clients
	if err := s.initGeoClient("srv-geo"); err != nil {
		return err
	}
	if err := s.initRateClient("srv-rate"); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

	// register with consul
	// jsonFile, err := os.Open("config.json")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// defer jsonFile.Close()

	// byteValue, _ := ioutil.ReadAll(jsonFile)

	// var result map[string]string
	// json.Unmarshal([]byte(byteValue), &result)

	err = s.Registry.Register(name, s.uuid, s.IpAddr, s.Port)
	if err != nil {
		return fmt.Errorf("failed register: %v", err)
	}
	log.Info().Msg("Successfully registered in consul")

	return srv.Serve(lis)
}

// Shutdown cleans up any processes
func (s *Server) Shutdown() {
	s.Registry.Deregister(s.uuid)
}

func (s *Server) initGeoClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.geoClient = geo.NewGeoClient(conn)
	return nil
}

func (s *Server) initRateClient(name string) error {
	conn, err := s.getGprcConn(name)
	if err != nil {
		return fmt.Errorf("dialer error: %v", err)
	}
	s.rateClient = rate.NewRateClient(conn)
	return nil
}

func (s *Server) getGprcConn(name string) (*grpc.ClientConn, error) {
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

var perfMetric *cb.PerfMetric

func init() {
	if cb.PERF_METRIC {
		perfMetric = cb.NewPerfMetric(context_bus.MetricSize, cb.Metric_Search_NearBy_Observation_1, cb.Metric_Search_NearBy_Logic_4)
	}
}

func (s *Server) GetMetric(ctx context.Context, req *pb.NearbyRequest) (*cb.PerfMetric, error) {
	log.Info().Msg("get perf metric and reset")
	if !cb.PERF_METRIC {
		return nil, nil
	}

	if perfMetric != nil {
		res := perfMetric.Calculate()

		return res, nil
	}

	return nil, nil
}

func (s *Server) ResetDB(ctx context.Context, req *pb.NearbyRequest) (*pb.SearchResult, error) {
	log.Info().Msg("reset databases")
	s.geoClient.ResetDB(ctx, new(geo.Request))
	s.rateClient.ResetDB(ctx, new(rate.Request))

	return new(pb.SearchResult), nil
}

// Nearby returns ids of nearby hotels ordered by ranking algo
func (s *Server) Nearby(ctx context.Context, req *pb.NearbyRequest) (*pb.SearchResult, error) {
	var s1, e1, s2, e2, s3, e3, e4 time.Time
	if cb.PERF_METRIC {
		s1 = time.Now()
	}

	// Context Bus
	cbCtx, cbOK := ContextBus.FromContext(ctx)

	// find nearby hotels
	// Context Bus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "search.Nearby.1",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "in Search Nearby",
			Paths:   nil,
		})
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "search.Nearby.2",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: fmt.Sprintf("nearby lat = %f", req.Lat),
			Paths:   nil,
		})
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "search.Nearby.3",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: fmt.Sprintf("nearby lon = %f", req.Lon),
			Paths:   nil,
		})
	} else {
		log.Info().Msg("in Search Nearby")

		log.Info().Msgf("nearby lat = %f", req.Lat)
		log.Info().Msgf("nearby lon = %f", req.Lon)
	}

	if cb.PERF_METRIC {
		e1 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Observation_1, float64(e1.UnixNano()-s1.UnixNano()))
	}

	nearby, err := s.geoClient.Nearby(ctx, &geo.Request{
		Lat:       req.Lat,
		Lon:       req.Lon,
		CBPayload: cbCtx.Payload(), // set ContextBus payload
	})
	if err != nil {
		return nil, err
	}

	if cb.PERF_METRIC {
		s2 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Logic_2, float64(s2.UnixNano()-e1.UnixNano()))
	}

	// Context Bus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "search.Nearby.4",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: fmt.Sprintf("get Nearby hotelId = %v", nearby.HotelIds),
			Paths:   nil,
		})
	} else {
		log.Info().Msgf("get Nearby hotelId = %v", nearby.HotelIds)
	}

	if cb.PERF_METRIC {
		e2 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Observation_2, float64(e2.UnixNano()-s2.UnixNano()))
	}

	// find rates for hotels
	rates, err := s.rateClient.GetRates(ctx, &rate.Request{
		HotelIds:  nearby.HotelIds,
		InDate:    req.InDate,
		OutDate:   req.OutDate,
		CBPayload: cbCtx.Payload(), // set ContextBus payload
	})

	if err != nil {
		return nil, err
	}

	if cb.PERF_METRIC {
		s3 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Logic_3, float64(s3.UnixNano()-e2.UnixNano()))
	}

	// Context Bus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "search.Nearby.5",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "",
			Paths:   nil,
		})
	}

	if cb.PERF_METRIC {
		e3 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Observation_3, float64(e3.UnixNano()-s3.UnixNano()))
	}

	// TODO(hw): add simple ranking algo to order hotel ids:
	// * geo distance
	// * price (best discount?)
	// * reviews

	// build the response
	res := new(pb.SearchResult)
	for _, ratePlan := range rates.RatePlans {
		//// Context Bus
		//if cbOK {
		//	ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
		//		Type: cb.EventRecorderType_EventRecorderServiceHandler,
		//		Name: "search.Nearby.6",
		//	}, &cb.EventMessage{
		//		Attrs:   nil,
		//		Message: fmt.Sprintf("get RatePlan HotelId = %s, Code = %s", ratePlan.HotelId, ratePlan.Code),
		//		Paths:   nil,
		//	})
		//} else {
		//	log.Trace().Msgf("get RatePlan HotelId = %s, Code = %s", ratePlan.HotelId, ratePlan.Code)
		//}
		res.HotelIds = append(res.HotelIds, ratePlan.HotelId)
	}

	if cb.PERF_METRIC {
		e4 = time.Now()
		perfMetric.AddLatency(cb.Metric_Search_NearBy_Logic_4, float64(e4.UnixNano()-e3.UnixNano()))
	}

	return res, nil
}
