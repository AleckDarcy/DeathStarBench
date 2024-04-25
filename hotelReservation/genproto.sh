protoc services/geo/proto/geo.proto --go_out=plugins=grpc:.
protoc services/profile/proto/profile.proto --proto_path=services/profile/proto --go_out=plugins=grpc:.
protoc services/rate/proto/rate.proto --go_out=plugins=grpc:.
protoc services/recommendation/proto/recommendation.proto --go_out=plugins=grpc:.
protoc services/reservation/proto/reservation.proto --go_out=plugins=grpc:.
protoc services/search/proto/search.proto --go_out=plugins=grpc:.
protoc services/user/proto/user.proto --go_out=plugins=grpc:.
