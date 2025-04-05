default:
  @just --list

build: templ
  go build -C ./cmd -o marketServer

run:
  ./cmd/marketServer

templ:
	go tool templ generate
	mv web/templates/*.go internal/views/templates/

TEST_ENV_FILE := "./build/test/.env.test"
TEST_DOCKER_COMPOSE := "./build/test/docker-compose.test.yml"
TEST_CMD := "go test ./internal/... -count=1"
BENCH_CMD := "go test ./internal/... -bench=. -benchmem -count=5 -run=^#"

tests:
  #!/usr/bin/env sh
  source {{TEST_ENV_FILE}}
  COMPOSE_BAKE=true docker compose -f {{TEST_DOCKER_COMPOSE}} --env-file {{TEST_ENV_FILE}} build
  for directory in `find internal/ -maxdepth 2 -type d`; do
    if `find ${directory} -maxdepth 1 -name "*_test.go" | read v`; then
      test_dir="market_test_${directory////_}";
      test_cmd="go test ./${directory} -count=1";
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} rm -f -s -v;
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} run --rm app ${test_cmd};
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} stop db;
    fi;
  done

