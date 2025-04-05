default:
  @just --list

build: templ
  go build -C ./cmd -o marketServer

run:
  ./cmd/marketServer

templ:
  #!/usr/bin/env sh
  go tool templ generate
  for directory in `find web/templates -type d`; do
    copy_to=${directory/web\/templates/internal/views/templates};
    if ! `find ${copy_to} -type d | read v`; then
      mkdir ${copy_to};
    fi;
    mv ${directory}/*.go ${copy_to}
  done

TEST_ENV_FILE := "./build/test/.env.test"
TEST_DOCKER_COMPOSE := "./build/test/docker-compose.test.yml"
TEST_CMD := "go test ./internal/... -count=1"
BENCH_CMD := "go test ./internal/... -bench=. -benchmem -count=5 -run=^#"

tests:
  #!/usr/bin/env sh
  source {{TEST_ENV_FILE}}
  COMPOSE_BAKE=true docker compose -f {{TEST_DOCKER_COMPOSE}} --env-file {{TEST_ENV_FILE}} build
  for directory in `find internal -type d`; do
    if `find ${directory} -maxdepth 1 -name "*_test.go" | read v`; then
      test_dir="market-test-${directory////-}";
      test_cmd="go test ./${directory} -count=1";
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} rm -f -s -v;
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} run --rm app ${test_cmd};
      docker compose -f {{TEST_DOCKER_COMPOSE}} -p ${test_dir} --env-file {{TEST_ENV_FILE}} stop db;
    fi;
  done

