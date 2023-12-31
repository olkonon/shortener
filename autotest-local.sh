#!/bin/bash
### Workflow statictest.yml ###
# shellcheck disable=SC2046
MSG=$(go vet -vettool=$(which statictest) ./...)
# shellcheck disable=SC2181
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
# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "==> Build binary ... [FAIL]"
  exit 128
else
  echo "==> Build binary ... [OK]"
fi
# Increment 1
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration1$ \
                    -binary-path=cmd/shortener/shortener)
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_1 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_1 ... [FAIL]"
  exit 128
fi
# Increment 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration2$ -source-path=.)
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_2 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_2 ... [FAIL]"
  exit 128
fi
# Increment 3
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration3$ -source-path=.)
# shellcheck disable=SC2181
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
# shellcheck disable=SC2181
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
# shellcheck disable=SC2181
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
# shellcheck disable=SC2181
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
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_7 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_7 ... [FAIL]"
  exit 128
fi

# Increment 8
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration8$ \
                    -binary-path=cmd/shortener/shortener)
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_8 ..... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_8 ... [FAIL]"
  exit 128
fi

# Increment 9
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration9$ \
                    -binary-path=cmd/shortener/shortener \
                    -source-path=. \
                    -file-storage-path=65118D10-CB73-41FA-B4EE-9D5685AD310D.json)
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_9 ..... [OK]"
  if test -f "65118D10-CB73-41FA-B4EE-9D5685AD310D.json";
  then
    rm 65118D10-CB73-41FA-B4EE-9D5685AD310D.json
  fi
  if test -f "/tmp/short-url-db.json";
  then
    rm /tmp/short-url-db.json
  fi
else
  if test -f "65118D10-CB73-41FA-B4EE-9D5685AD310D.json";
  then
    rm 65118D10-CB73-41FA-B4EE-9D5685AD310D.json
  fi
  if test -f "/tmp/short-url-db.json";
  then
    rm /tmp/short-url-db.json
  fi
  echo "$MSG"
  echo "==> Test INC_9 ... [FAIL]"
  exit 128
fi

# Increment 10
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter10-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration10$ \
                    -binary-path=cmd/shortener/shortener \
                    -source-path=. \
                    -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_10 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_10 .. [FAIL]"
  exit 128
fi

# Increment 11
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter11-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration11$ \
                    -binary-path=cmd/shortener/shortener \
                    -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_11 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_11 .. [FAIL]"
  exit 128
fi


# Increment 12
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter12-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration12$ \
                    -binary-path=cmd/shortener/shortener \
                    -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_12 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_12 .. [FAIL]"
  exit 128
fi

# Increment 13
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter13-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration13$ \
                    -binary-path=cmd/shortener/shortener \
                    -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_13 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_13 .. [FAIL]"
  exit 128
fi

# Increment 14
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter14-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration14$ \
                    -binary-path=cmd/shortener/shortener \
                   -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_14 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_14 .. [FAIL]"
  exit 128
fi

# Increment 15
docker stop $(docker ps -a -q) > /dev/null
docker rm -f $(docker ps -a -q) > /dev/null
docker run --name iter15-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=praktikum -d arm64v8/postgres > /dev/null
sleep 2
MSG=$(shortenertestbeta -test.v -test.run=^TestIteration15$ \
                    -binary-path=cmd/shortener/shortener \
                   -database-dsn='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable')
# shellcheck disable=SC2181
if [ $? -eq 0 ]; then
  echo "==> Test INC_15 .... [OK]"
else
  echo "$MSG"
  echo "==> Test INC_15 .. [FAIL]"
  exit 128
fi