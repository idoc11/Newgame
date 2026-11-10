#!/bin/sh
# Собирает оба exe для Windows. Нужен только Go; работает с любой ОС.
set -e
export GOOS=windows GOARCH=amd64
go build -trimpath -ldflags "-H windowsgui -s -w" -o dist/NeonRift.exe .
go build -trimpath -ldflags "-H windowsgui -s -w -X main.startPage=fps.html -X 'main.title=Сектор Ноль' -X main.fixedAddr=127.0.0.1:47814" -o dist/SectorZero.exe .
