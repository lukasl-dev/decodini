default:
    @just --list

test:
    go test ./...

bench:
    cd bench && go test -bench=. -benchmem
