#!/bin/bash
### Workflow statictest.yml ###
# shellcheck disable=SC2046
MSG=$(go vet -vettool=$(which statictest) ./...)
if [ $? -eq 0 ]; then
  echo "==> Statictest ..... [OK]"
else
  echo "$MSG"
  echo "==> Statictest ... [FAIL]"
  exit 128
fi
### Workflow   shortenertest.yml ###
#build binary
go build -o cmd/shortener/shortener cmd/shortener/main.go
if [ $? -ne 0 ]; then
  echo "==> Build binary ... [FAIL]"
  exit 128
else
  echo "==> Build binary ... [OK]"
fi
# Increment 1
MSG=$(shortenertest -test.v -test.run=^TestIteration1$ \
                    -binary-path=cmd/shortener/shortener)
if [ $? -eq 0 ]; then
  echo "==> Test INC_1 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_1 ... [FAIL]"
  exit 128
fi
# Increment 2
MSG=$(shortenertest -test.v -test.run=^TestIteration2$ -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_2 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_2 ... [FAIL]"
  exit 128
fi
# Increment 3
MSG=$(shortenertest -test.v -test.run=^TestIteration3$ -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_3 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_3 ... [FAIL]"
  exit 128
fi