package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/delimitrou/DeathStarBench/hotelreservation/registry"
	pb "github.com/delimitrou/DeathStarBench/hotelreservation/services/rate/proto"
	"github.com/delimitrou/DeathStarBench/hotelreservation/tls"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/AleckDarcy/ContextBus"
	cb_configure "github.com/AleckDarcy/ContextBus/configure"
	cb "github.com/AleckDarcy/ContextBus/proto"
	"github.com/delimitrou/DeathStarBench/hotelreservation/services/context_bus"
)

const name = "srv-rate"

// Server implements the rate service
type Server struct {
	Tracer      opentracing.Tracer
	Port        int
	IpAddr      string
	MongoClient *mongo.Client
	Registry    *registry.Client
	MemcClient  *memcache.Client
	uuid        string

	CBConfig *cb_configure.ServerConfigure
}

// Run starts the server
func (s *Server) Run() error {
	opentracing.SetGlobalTracer(s.Tracer)

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

	pb.RegisterRateServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}

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

type RoomType struct {
	BookableRate       float64 `bson:"bookableRate"`
	Code               string  `bson:"code"`
	RoomDescription    string  `bson:"roomDescription"`
	TotalRate          float64 `bson:"totalRate"`
	TotalRateInclusive float64 `bson:"totalRateInclusive"`
}

type RatePlan struct {
	HotelId  string    `bson:"hotelId"`
	Code     string    `bson:"code"`
	InDate   string    `bson:"inDate"`
	OutDate  string    `bson:"outDate"`
	RoomType *RoomType `bson:"roomType"`
}

func (s *Server) ResetDataBases() {
	collection := s.MongoClient.Database("rate-db").Collection("inventory")
	if err := collection.Drop(context.Background()); err != nil {
		log.Error().Msgf("Drop Collection inventory from Database rate-db failed: %v\n", err)
	} else {
		log.Info().Msg("Drop Collection inventory from Database rate-db succeeded")
	}

	if err := s.MemcClient.DeleteAll(); err != nil {
		log.Error().Msgf("Clear memcached fail: %v\n", err)
	} else {
		log.Info().Msg("Clear memcached succeeded")
	}

	newRatePlans := []interface{}{
		RatePlan{
			"1",
			"RACK",
			"2015-04-09",
			"2015-04-10",
			&RoomType{
				109.00,
				"KNG",
				"King sized bed",
				109.00,
				123.17,
			},
		},
		RatePlan{
			"2",
			"RACK",
			"2015-04-09",
			"2015-04-10",
			&RoomType{
				139.00,
				"QN",
				"Queen sized bed",
				139.00,
				153.09,
			},
		},
		RatePlan{
			"3",
			"RACK",
			"2015-04-09",
			"2015-04-10",
			&RoomType{
				109.00,
				"KNG",
				"King sized bed",
				109.00,
				123.17,
			},
		},
	}

	for i := 7; i <= 80; i++ {
		if i%3 != 0 {
			continue
		}

		hotelID := strconv.Itoa(i)

		endDate := "2015-04-"
		if i%2 == 0 {
			endDate = fmt.Sprintf("%s17", endDate)
		} else {
			endDate = fmt.Sprintf("%s24", endDate)
		}

		rate := 109.00
		rateInc := 123.17
		if i%5 == 1 {
			rate = 120.00
			rateInc = 140.00
		} else if i%5 == 2 {
			rate = 124.00
			rateInc = 144.00
		} else if i%5 == 3 {
			rate = 132.00
			rateInc = 158.00
		} else if i%5 == 4 {
			rate = 232.00
			rateInc = 258.00
		}

		newRatePlans = append(
			newRatePlans,
			RatePlan{
				hotelID,
				"RACK",
				"2015-04-09",
				endDate,
				&RoomType{
					rate,
					"KNG",
					"King sized bed",
					rate,
					rateInc,
				},
			},
		)
	}

	_, err := collection.InsertMany(context.TODO(), newRatePlans)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("Successfully reset rate DB")
}

func (s *Server) ResetDB(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	log.Info().Msg("reset databases")
	s.ResetDataBases()

	return new(pb.Result), nil
}

// GetRates gets rates for hotels for specific date range.
func (s *Server) GetRates(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	// Context Bus
	cbCtx, cbOK := ContextBus.FromContext(ctx)

	res := new(pb.Result)

	ratePlans := make(RatePlans, 0)

	hotelIds := []string{}
	rateMap := make(map[string]struct{})
	for _, hotelID := range req.HotelIds {
		hotelIds = append(hotelIds, hotelID)
		rateMap[hotelID] = struct{}{}
	}

	// first check memcached(get-multi)
	var memSpan opentracing.Span

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "rate.GetRates.1",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "",
			Paths:   nil,
		})
	} else {
		memSpan, _ = opentracing.StartSpanFromContext(ctx, "memcached_get_multi_rate")
		memSpan.SetTag("span.kind", "client")
	}

	resMap, err := s.MemcClient.GetMulti(hotelIds)

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "rate.GetRates.2",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "",
			Paths:   nil,
		})
	} else {
		memSpan.Finish()
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex
	if err != nil && err != memcache.ErrCacheMiss {
		log.Panic().Msgf("Memmcached error while trying to get hotel [id: %v]= %s", hotelIds, err)
	} else {
		for hotelId, item := range resMap {
			rateStrs := strings.Split(string(item.Value), "\n")

			if cbOK {
				ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
					Type: cb.EventRecorderType_EventRecorderServiceHandler,
					Name: "rate.GetRates.3",
				}, &cb.EventMessage{
					Attrs:   nil,
					Message: fmt.Sprintf("memc hit, hotelId = %s,rate strings: %v", hotelId, rateStrs),
					Paths:   nil,
				})
			} else {
				log.Trace().Msgf("memc hit, hotelId = %s,rate strings: %v", hotelId, rateStrs)
			}

			for _, rateStr := range rateStrs {
				if len(rateStr) != 0 {
					rateP := new(pb.RatePlan)
					json.Unmarshal([]byte(rateStr), rateP)
					ratePlans = append(ratePlans, rateP)
				}
			}

			delete(rateMap, hotelId)
		}

		wg.Add(len(rateMap))
		for hotelId := range rateMap {
			go func(id string) {
				if cbOK {
					ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
						Type: cb.EventRecorderType_EventRecorderServiceHandler,
						Name: "rate.GetRates.4",
					}, &cb.EventMessage{
						Attrs:   nil,
						Message: fmt.Sprintf("memc miss, hotelId = %s", id),
						Paths:   nil,
					})

					ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
						Type: cb.EventRecorderType_EventRecorderServiceHandler,
						Name: "rate.GetRates.5",
					}, &cb.EventMessage{
						Attrs:   nil,
						Message: "memcached miss, set up mongo connection",
						Paths:   nil,
					})
				} else {
					log.Info().Msgf("memc miss, hotelId = %s", id)
					log.Info().Msg("memcached miss, set up mongo connection")
				}

				mongoSpan, _ := opentracing.StartSpanFromContext(ctx, "mongo_rate")
				mongoSpan.SetTag("span.kind", "client")

				// memcached miss, set up mongo connection
				collection := s.MongoClient.Database("rate-db").Collection("inventory")
				curr, err := collection.Find(context.TODO(), bson.D{})
				if err != nil {
					log.Error().Msgf("Failed get rate data: ", err)
				}

				tmpRatePlans := make(RatePlans, 0)
				curr.All(context.TODO(), &tmpRatePlans)
				if err != nil {
					log.Error().Msgf("Failed get rate data: ", err)
				}

				mongoSpan.Finish()

				memcStr := ""
				if err != nil {
					log.Panic().Msgf("Tried to find hotelId [%v], but got error", id, err.Error())
				} else {
					for _, r := range tmpRatePlans {
						mutex.Lock()
						ratePlans = append(ratePlans, r)
						mutex.Unlock()
						rateJson, err := json.Marshal(r)
						if err != nil {
							log.Error().Msgf("Failed to marshal plan [Code: %v] with error: %s", r.Code, err)
						}
						memcStr = memcStr + string(rateJson) + "\n"
					}
				}
				go s.MemcClient.Set(&memcache.Item{Key: id, Value: []byte(memcStr)})

				defer wg.Done()
			}(hotelId)
		}
	}
	wg.Wait()

	sort.Sort(ratePlans)
	res.RatePlans = ratePlans

	return res, nil
}

type RatePlans []*pb.RatePlan

func (r RatePlans) Len() int {
	return len(r)
}

func (r RatePlans) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r RatePlans) Less(i, j int) bool {
	return r[i].RoomType.TotalRate > r[j].RoomType.TotalRate
}
