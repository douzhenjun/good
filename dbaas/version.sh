#!/bin/bash

readonly version=020100
readonly pg="psql -U postgres -h 127.0.0.1 -p 5432"

## connect to database
db_try=0
while true; do
  if $pg -l; then
    echo "connect to database successful"
    break
  fi
  if [ $db_try -eq 3 ]; then
    echo "connect to database fail, exit"
    exit 1
  fi
  db_try=$db_try+1
  sleep 3
  echo "connect to database try" $db_try
done

db_try=0
while true; do
  if $pg -l | grep dbaas; then
    echo "find database dbaas successful"
    break
  fi
  if [ $db_try -eq 3 ]; then
    echo "not found database dbaas, exit"
    exit 1
  fi
  db_try=$db_try+1
  sleep 3
  echo "find database dbaas try" $db_try
done

## check version
oldVersion=$($pg -t -d dbaas -A -c "select value from misc_config where key = 'version'")
echo ver: $version, old: "${oldVersion:=0}"
if [ "$oldVersion" -eq $version ]; then
  exit 0
fi

## exec sql
path="/opt/dbaas/support-files/database/version/"
files=$(ls $path)
for sql in $files; do
  v=${sql%.*}
  if [ "$v" -gt "$oldVersion" ] && [ "$v" -le $version ]; then
    $pg -d dbaas -f $path"$sql" && echo exec $path"$sql" successful
  fi
done

# exec menu sql
path="/opt/dbaas/support-files/database/menu/"
files=$(ls $path)
for sql in $files; do
  v=${sql%.*}
  if [ "$v" -gt "$oldVersion" ] && [ "$v" -le $version ]; then
    $pg -d capricorn -f $path"$sql" && echo exec $path"$sql" successful
  fi
done

## write version
if [ $oldVersion -eq 0 ]; then
  $pg -d dbaas -c "insert into misc_config (key, value) values ('version', '$version')" && echo version $version insert into 'misc_config' successful
else
  $pg -d dbaas -c "update misc_config set value='$version' where key = 'version'" && echo update 'misc_config' set version=$version successful
fi
