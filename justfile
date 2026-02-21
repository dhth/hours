alias r := run
alias b := build
alias t := test
alias f := fmt
alias i := install
alias c := check
alias l := lint
alias up := upgrade
alias upt := upgradet
alias ti := tidy
alias v := vuln
alias us := update-snapshots

@default:
    just --choose

@all:
    just fmt
    just check
    just test

run:
    go run .

build:
    go build -ldflags "-w -s" .

test:
    go test ./...

fmt:
    gofumpt -l -w .

install:
    go install .

check:
    golangci-lint run

lint:
    golangci-lint run

upgrade:
    go get -u ./...

upgradet:
    go get -t -u ./...

tidy:
    go mod tidy

vuln:
    govulncheck ./...

update-snapshots $UPDATE_SNAPS='true':
    just test
