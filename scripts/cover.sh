#!/usr/bin/env bash
go test ./... -coverprofile=covprofile
go tool cover -html=covprofile -o coverage.html
open coverage.html
rm covprofile
