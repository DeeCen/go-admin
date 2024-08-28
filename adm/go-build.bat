go env -w CGO_CFLAGS="-g -O2 -Wno-return-local-addr"
go build ./
pause