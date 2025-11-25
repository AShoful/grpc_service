#!/bin/bash

go test ./pkg/... -coverpkg=./pkg/service,./pkg/handler -coverprofile=test_coverage.out

go tool cover -func=test_coverage.out

go tool cover -html=test_coverage.out -o coverage.html

echo "Tests are complete. The HTML report is available at coverage.html"

if [[ "$OSTYPE" == "darwin"* ]]; then
  open coverage.html
fi

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  xdg-open coverage.html
fi
