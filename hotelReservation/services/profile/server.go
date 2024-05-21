package profile

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/delimitrou/DeathStarBench/hotelreservation/registry"
	pb "github.com/delimitrou/DeathStarBench/hotelreservation/services/profile/proto"
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

const name = "srv-profile"

// Server implements the profile service
type Server struct {
	Tracer      opentracing.Tracer
	uuid        string
	Port        int
	IpAddr      string
	MongoClient *mongo.Client
	Registry    *registry.Client
	MemcClient  *memcache.Client

	CBConfig *cb_configure.ServerConfigure
}

// Run starts the server
func (s *Server) Run() error {
	opentracing.SetGlobalTracer(s.Tracer)

	if s.Port == 0 {
		return fmt.Errorf("server port must be set")
	}

	s.uuid = uuid.New().String()

	log.Trace().Msgf("in run s.IpAddr = %s, port = %d", s.IpAddr, s.Port)

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

	pb.RegisterProfileServer(srv, s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Msgf("failed to configure listener: %v", err)
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

type Hotel struct {
	Id          string   `bson:"id"`
	Name        string   `bson:"name"`
	PhoneNumber string   `bson:"phoneNumber"`
	Description string   `bson:"description"`
	Address     *Address `bson:"address"`
}

type Address struct {
	StreetNumber string  `bson:"streetNumber"`
	StreetName   string  `bson:"streetName"`
	City         string  `bson:"city"`
	State        string  `bson:"state"`
	Country      string  `bson:"country"`
	PostalCode   string  `bson:"postalCode"`
	Lat          float32 `bson:"lat"`
	Lon          float32 `bson:"lon"`
}

func (s *Server) ResetDatabases() {
	collection := s.MongoClient.Database("profile-db").Collection("hotels")
	if err := collection.Drop(context.Background()); err != nil {
		log.Error().Msgf("Drop Collection hotels from Database profile-db failed: %v\n", err)
	} else {
		log.Info().Msg("Drop Collection hotels from Database profile-db succeeded")
	}

	if err := s.MemcClient.DeleteAll(); err != nil {
		log.Error().Msgf("Clear memcached fail: %v\n", err)
	} else {
		log.Info().Msg("Clear memcached succeeded")
	}

	newProfiles := []interface{}{
		Hotel{
			"1",
			"Clift Hotel",
			"(415) 775-4700",
			"A 6-minute walk from Union Square and 4 minutes from a Muni Metro station, this luxury hotel designed by Philippe Starck features an artsy furniture collection in the lobby, including work by Salvador Dali.",
			&Address{
				"495",
				"Geary St",
				"San Francisco",
				"CA",
				"United States",
				"94102",
				37.7867,
				-122.4112,
			},
		},
		Hotel{
			"2",
			"W San Francisco",
			"(415) 777-5300",
			"Less than a block from the Yerba Buena Center for the Arts, this trendy hotel is a 12-minute walk from Union Square.",
			&Address{
				"181",
				"3rd St",
				"San Francisco",
				"CA",
				"United States",
				"94103",
				37.7854,
				-122.4005,
			},
		},
		Hotel{
			"3",
			"Hotel Zetta",
			"(415) 543-8555",
			"A 3-minute walk from the Powell Street cable-car turnaround and BART rail station, this hip hotel 9 minutes from Union Square combines high-tech lodging with artsy touches.",
			&Address{
				"55",
				"5th St",
				"San Francisco",
				"CA",
				"United States",
				"94103",
				37.7834,
				-122.4071,
			},
		},
		Hotel{
			"4",
			"Hotel Vitale",
			"(415) 278-3700",
			"This waterfront hotel with Bay Bridge views is 3 blocks from the Financial District and a 4-minute walk from the Ferry Building.",
			&Address{
				"8",
				"Mission St",
				"San Francisco",
				"CA",
				"United States",
				"94105",
				37.7936,
				-122.3930,
			},
		},
		Hotel{
			"5",
			"Phoenix Hotel",
			"(415) 776-1380",
			"Located in the Tenderloin neighborhood, a 10-minute walk from a BART rail station, this retro motor lodge has hosted many rock musicians and other celebrities since the 1950s. It’s a 4-minute walk from the historic Great American Music Hall nightclub.",
			&Address{
				"601",
				"Eddy St",
				"San Francisco",
				"CA",
				"United States",
				"94109",
				37.7831,
				-122.4181,
			},
		},
		Hotel{
			"6",
			"St. Regis San Francisco",
			"(415) 284-4000",
			"St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
			&Address{
				"125",
				"3rd St",
				"San Francisco",
				"CA",
				"United States",
				"94109",
				37.7863,
				-122.4015,
			},
		},
	}

	for i := 7; i <= 80; i++ {
		hotelID := strconv.Itoa(i)
		phoneNumber := fmt.Sprintf("(415) 284-40%s", hotelID)

		lat := 37.7835 + float32(i)/500.0*3
		lon := -122.41 + float32(i)/500.0*4

		newProfiles = append(
			newProfiles,
			Hotel{
				hotelID,
				"St. Regis San Francisco",
				phoneNumber,
				"St. Regis Museum Tower is a 42-story, 484 ft skyscraper in the South of Market district of San Francisco, California, adjacent to Yerba Buena Gardens, Moscone Center, PacBell Building and the San Francisco Museum of Modern Art.",
				&Address{
					"125",
					"3rd St",
					"San Francisco",
					"CA",
					"United States",
					"94109",
					lat,
					lon,
				},
			},
		)
	}

	_, err := collection.InsertMany(context.TODO(), newProfiles)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("Successfully reset profile DB")
}

func (s *Server) ResetDB(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	log.Info().Msg("reset databases")
	s.ResetDatabases()

	return new(pb.Result), nil
}

// GetProfiles returns hotel profiles for requested IDs
func (s *Server) GetProfiles(ctx context.Context, req *pb.Request) (*pb.Result, error) {
	// ContextBus
	cbCtx, cbOK := ContextBus.FromContext(ctx)

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "profile.GetProfiles.1",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "In GetProfiles",
			Paths:   nil,
		})
	} else {
		log.Info().Msgf("In GetProfiles")
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// one hotel should only have one profile
	hotelIds := make([]string, 0)
	profileMap := make(map[string]struct{})
	for _, hotelId := range req.HotelIds {
		hotelIds = append(hotelIds, hotelId)
		profileMap[hotelId] = struct{}{}
	}

	var memSpan opentracing.Span

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "profile.GetProfiles.2",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "",
			Paths:   nil,
		})
	} else {
		memSpan, _ = opentracing.StartSpanFromContext(ctx, "memcached_get_profile")
		memSpan.SetTag("span.kind", "client")
	}

	resMap, err := s.MemcClient.GetMulti(hotelIds)

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "profile.GetProfiles.3",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "",
			Paths:   nil,
		})
	} else {
		memSpan.Finish()
	}

	res := new(pb.Result)
	hotels := make([]*pb.Hotel, 0)

	if err != nil && err != memcache.ErrCacheMiss {
		log.Panic().Msgf("Tried to get hotelIds [%v], but got memmcached error = %s", hotelIds, err)
	} else {
		for hotelId, item := range resMap {
			profileStr := string(item.Value)
			log.Trace().Msgf("memc hit with %v", profileStr)

			hotelProf := new(pb.Hotel)
			json.Unmarshal(item.Value, hotelProf)
			hotels = append(hotels, hotelProf)
			delete(profileMap, hotelId)
		}

		wg.Add(len(profileMap))
		for hotelId := range profileMap {
			go func(hotelId string) {
				var hotelProf *pb.Hotel

				collection := s.MongoClient.Database("profile-db").Collection("hotels")

				mongoSpan, _ := opentracing.StartSpanFromContext(ctx, "mongo_profile")
				mongoSpan.SetTag("span.kind", "client")
				err := collection.FindOne(context.TODO(), bson.D{{"id", hotelId}}).Decode(&hotelProf)
				mongoSpan.Finish()

				if err != nil {
					log.Error().Msgf("Failed get hotels data: ", err)
				}

				mutex.Lock()
				hotels = append(hotels, hotelProf)
				mutex.Unlock()

				profJson, err := json.Marshal(hotelProf)
				if err != nil {
					log.Error().Msgf("Failed to marshal hotel [id: %v] with err:", hotelProf.Id, err)
				}
				memcStr := string(profJson)

				// write to memcached
				go s.MemcClient.Set(&memcache.Item{Key: hotelId, Value: []byte(memcStr)})
				defer wg.Done()
			}(hotelId)
		}
	}
	wg.Wait()

	res.Hotels = hotels

	// ContextBus
	if cbOK {
		ContextBus.OnSubmission(cbCtx, &cb.EventWhere{}, &cb.EventRecorder{
			Type: cb.EventRecorderType_EventRecorderServiceHandler,
			Name: "profile.GetProfiles.4",
		}, &cb.EventMessage{
			Attrs:   nil,
			Message: "In GetProfiles after getting resp",
			Paths:   nil,
		})
	} else {
		log.Trace().Msgf("In GetProfiles after getting resp")
	}

	return res, nil
}
