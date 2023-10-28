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
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration1$ \
                    -binary-path=cmd/shortener/shortener)
if [ $? -eq 0 ]; then
  echo "==> Test INC_1 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_1 ... [FAIL]"
  exit 128
fi
# Increment 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration2$ -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_2 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_2 ... [FAIL]"
  exit 128
fi
# Increment 3
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration3$ -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_3 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_3 ... [FAIL]"
  exit 128
fi
# Increment 4
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration4$ \
                    -binary-path=cmd/shortener/shortener \
                    -server-port=8080)
if [ $? -eq 0 ]; then
  echo "==> Test INC_4 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_4 ... [FAIL]"
  exit 128
fi

# Increment 5
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration5$ \
                    -binary-path=cmd/shortener/shortener \
                    -server-port=8080)
if [ $? -eq 0 ]; then
  echo "==> Test INC_5 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_5 ... [FAIL]"
  exit 128
fi

# Increment 6
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration6$ \
                    -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_6 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_6 ... [FAIL]"
  exit 128
fi

# Increment 7
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration7$ \
                    -binary-path=cmd/shortener/shortener \
                    -source-path=.)
if [ $? -eq 0 ]; then
  echo "==> Test INC_7 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_7 ... [FAIL]"
  exit 128
fi