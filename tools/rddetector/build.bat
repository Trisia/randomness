@echo off

RD /S/Q target >nul 2>nul
MKDIR target

SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w" -o target/rddetector.exe
ECHO [+] Windows/amd64 rddetector.exe

SET GOOS=linux
SET GOARCH=amd64
go build -ldflags "-s -w" -o target/rddetector_linux_amd64
ECHO [+] Linux/amd64   rddetector_linux_amd64

SET GOOS=linux
SET GOARCH=arm64
go build -ldflags "-s -w" -o target/rddetector_linux_arm64
ECHO [+] Linux/arm64   rddetector_linux_arm64

@echo on