protoc services/geo/proto/geo.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/profile/proto/profile.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/rate/proto/rate.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/recommendation/proto/recommendation.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/reservation/proto/reservation.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/search/proto/search.proto --proto_path=services --go_out=plugins=grpc:services
protoc services/user/proto/user.proto --proto_path=services --go_out=plugins=grpc:services
