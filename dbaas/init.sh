#!/bin/sh

echo "Database Password is $pg_passwd"

##connect to database
db_test=1
while true
do
  if [  $db_test == 4 ]
  then
    echo "Could not connect to database,exit!"
    exit 1
  fi
  psql -U postgres -h 127.0.0.1 -p 5432 password=$pg_passwd -l
  if [ $? != 0 ]
  then
    echo "Database connnect failed,retry!"
    let db_test=$db_test+1
    sleep 5
  else
    echo "Database connect sussced,contiuned!"
    break
  fi
done
psql -U postgres -h 127.0.0.1 -p 5432 password=$pg_passwd -l|grep dbaas
if [ $? == 1 ]
then
  python3 /opt/dbaas/support-files/scripts/initdb.py 127.0.0.1 5432 $pg_passwd
  db_test=1
  while true
  do
    if [  $db_test == 4 ]
    then
      echo "Could not find database capricorn,exit!"
      exit 1
    fi
    psql -U postgres -h 127.0.0.1 -p 5432 password=$pg_passwd -l|grep capricorn
    if [ $? == 1 ]
    then
      echo "Could not find database capricorn,retry!"
      let db_test=$db_test+1
      sleep 5
    else
      echo "Find database capricorn,contiuned!"
      break
    fi
  done
  table_test=1
  while true
  do
    if [  $table_test == 4 ]
    then
      echo "Could not found table menu_component,exit!"
      exit 1
    fi
    tab_exist=`psql -U postgres -h 127.0.0.1 -p 5432 -d capricorn -c "select count(*) from pg_class where relname = 'menu_component';"|sed -n '3p'|sed 's/^[ \t]*//g'`
    if [ $tab_exist == 0 ]
    then
      echo "Could not found table menu_component,retry!"
      let table_test=$table_test+1
      sleep 5
    else
      echo "Init menu data!"
      python3 /opt/dbaas/support-files/scripts/init-menu.py 127.0.0.1 5432 $pg_passwd
      python3 /opt/dbaas/support-files/scripts/init_role_component.py 127.0.0.1 5432 $pg_passwd
      break
    fi
  done

else
  echo "database dbaas already exists!"
fi
/opt/dbaas/version.sh
/usr/bin/supervisord
