gen_proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pkg/proto/conversation.proto

migrateup:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/conversations_db?sslmode=disable" --verbose up

migrateup1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/conversations_db?sslmode=disable" --verbose up 1

migratedown:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/conversations_db?sslmode=disable" --verbose down

migratedown1:
	 migrate -path migrations -database "postgresql://user:user@localhost:5440/conversations_db?sslmode=disable" --verbose down 1