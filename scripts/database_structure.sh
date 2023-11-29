#!/usr/bin/env sh

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )/../" &> /dev/null && pwd )

DBMATE_PATH=$(which dbmate)

if [ "$?" != "0" ] ; then
  echo " - [ ERROR ] dbmate is not installed. Please install it and execute this command again. Follow the instructions here!: https://github.com/amacneil/dbmate#installation";

  exit 1;
fi

echo " - Running database migrations...";

export DATABASE_URL="postgres://username:password@localhost:5432/zenwallet?sslmode=disable";

dbmate wait && dbmate up && ([ -d "db/seeds" ] && dbmate --migrations-dir=db/seeds --migrations-table=seed_migrations up || true);
