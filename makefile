build:
	go build -o bin/app

run:
	build
	./bin/app

test:
	go test -v ./... -count=1

seq?=mg
migration_create:
	/Users/labyrinth/.go/bin/migrate create -ext sql -dir ./api/migration -seq ${seq}

migration_up:
	/Users/labyrinth/.go/bin/migrate -path ./api/migration -database "postgres://localhost:5432/ecom?sslmode=disable" -verbose up ${v}

migration_down:
	/Users/labyrinth/.go/bin/migrate -path ./api/migration -database "postgres://localhost:5432/ecom?sslmode=disable" -verbose down ${v}

migration_fix:
	/Users/labyrinth/.go/bin/migrate -path ./api/migration -database "postgres://localhost:5432/ecom?sslmode=disable" force $(v)
# make migration_fix v=1
# v?=v_default # Assign default value if not provided

air_be:
	cd ./api && /Users/labyrinth/.go/bin/air

air_fe:
	cd ./web && /Users/labyrinth/.go/bin/air

start_db:
	/opt/homebrew/opt/postgresql@16/bin/postgres -D /opt/homebrew/var/postgresql@16 &
# Run inbackground (&)

stop_db:
	@-pkill -SIGTERM -f "postgres"
# lsof -i :3000
# lsof -ti :3000
# fuser 3000/tcp
# fuser 3000/udp


# STRIPE_SECRET=sk_test_mXWrR1RN6fjIJnDsLPq1mAGX
# STRIPE_KEY=pk_test_lOwqX0SiQCGm7wSkqNoBgMLc
# GOSTRIPE_PORT=4000
# API_PORT=4001

# ## build: builds all binaries
# build: clean build_front build_back
# 	@printf "All binaries built!\n"

# ## clean: cleans all binaries and runs go clean
# clean:
# 	@echo "Cleaning..."
# 	@- rm -f dist/*
# 	@go clean
# 	@echo "Cleaned!"

# ## build_front: builds the front end
# build_front:
# 	@echo "Building front end..."
# 	@go build -o dist/gostripe ./cmd/web
# 	@echo "Front end built!"

# ## build_back: builds the back end
# build_back:
# 	@echo "Building back end..."
# 	@go build -o dist/gostripe_api ./cmd/api
# 	@echo "Back end built!"

# ## start: starts front and back end
# start: start_front start_back
	
# ## start_front: starts the front end
# start_front: build_front
# 	@echo "Starting the front end..."
# 	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe -port=${GOSTRIPE_PORT} &
# 	@echo "Front end running!"

# ## start_back: starts the back end
# start_back: build_back
# 	@echo "Starting the back end..."
# 	@env STRIPE_KEY=${STRIPE_KEY} STRIPE_SECRET=${STRIPE_SECRET} ./dist/gostripe_api -port=${API_PORT} &
# 	@echo "Back end running!"

# ## stop: stops the front and back end
# stop: stop_front stop_back
# 	@echo "All applications stopped"

# ## stop_front: stops the front end
# stop_front:
# 	@echo "Stopping the front end..."
# 	@-pkill -SIGTERM -f "gostripe -port=${GOSTRIPE_PORT}"
# 	@echo "Stopped front end"

# ## stop_back: stops the back end
# stop_back:
# 	@echo "Stopping the back end..."
# 	@-pkill -SIGTERM -f "gostripe_api -port=${API_PORT}"
# 	@echo "Stopped back end"