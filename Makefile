build:
	go build

templ:
	templ generate
	mv views/templates/*.go views/

run:
	./ECommerce

buildRun: build run

TEST_CMD := "go test ./... -count=1"
tests:
	CMD=${TEST_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test rm -f -v
	CMD=${TEST_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test build
	CMD=${TEST_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test up --abort-on-container-exit

BENCH_CMD := "go test ./... -bench=. -benchmem -count=5 -run=^\#"
benchmarks:
	CMD=${BENCH_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test rm -f -v
	CMD=${BENCH_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test build
	CMD=${BENCH_CMD} docker compose -f ./build/docker-compose.test.yml --env-file ./build/.env.test up --abort-on-container-exit
