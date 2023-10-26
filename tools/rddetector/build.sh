#!/bin/bash

mkdir target/

GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o target/rddetector.exe
echo "[+] Windows/amd64 rddetector.exe"
GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w" -o target/rddetector_linux_amd64
echo "[+] Linux/amd64   rddetector_linux_amd64"
GOOS=linux   GOARCH=arm64 go build -ldflags "-s -w" -o target/rddetector_linux_arm64
echo "[+] Linux/arm64   rddetector_linux_arm64"

