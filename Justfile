set dotenv-load

default:
  @just --list

run:
  ./tmp/bin/marketServer

build:
  go build -C ./cmd -o ../tmp/bin/marketServer

templ:
  go tool templ generate --path "web/templates"
  just mv-templ-files

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

live:  
  #!/usr/bin/env -S parallel --shebang --ungroup --jobs {{ num_cpus() }}
  just live-templ 
  just live-build
  just live-sync-web

live-templ:
  #!/usr/bin/env sh
  go tool templ generate --watch --proxy="http://localhost:8080" \
  --path "web/templates" \
  --open-browser=false \
  --cmd "just mv-templ-files" \
  -v

live-build:
  go run github.com/air-verse/air@latest \
  --build.cmd "just build" \
  --build.bin "tmp/bin/marketServer" \
  --build.delay "100" \
  --build.exclude_dir "build,web,tmp" \
  --build.include_ext "go" \
  --build.stop_on_error "false" \
  --misc.clean_on_exit true

live-sync-web:
  go run github.com/air-verse/air@latest \
  --build.cmd "go tool templ generate --notify-proxy" \
  --build.bin "true" \
  --build.delay "100" \
  --build.exclude_dir "" \
  --build.include_dir "web" \
  --build.include_ext "js, css, tmpl"

mv-templ-files:
  #!/usr/bin/env sh
  for directory in `find web/templates -type d`; do
    if `find "${directory}" -maxdepth 1 -name "*_templ.go" | read v`; then
      copy_to=${directory/web\/templates/internal/views/templates};
      if ! `find ${copy_to} -type d | read v`; then
        mkdir ${copy_to};
      fi;
      mv "${directory}"/*.go "${copy_to}"/;
    fi;
  done
