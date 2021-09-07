'''
@Description:
@version:
@Company: iwhalecloud
@Author: Zhu.Jie
@Date: 2020-04-26 14:55:27
@LastEditors: Zhu.Jie
@LastEditTime: 2020-05-26 13:57:09
'''
import sys
import os
import time


def init_db(host, port, password):
    filename = os.path.abspath(os.path.join(sys.path[0], os.path.pardir)) + "/database/init.sql"
    cmdbfilename = os.path.abspath(os.path.join(sys.path[0], os.path.pardir)) + "/database/cmdbinit.sql"
    collectfilename = os.path.abspath(os.path.join(sys.path[0], os.path.pardir)) + "/database/collectinit.sql"
    psql_password = "PGPASSWORD=\"" + password + "\""
    doccommands = "dos2unix  " + collectfilename
    os.system(doccommands)
    time.sleep(10)
    commands = []
    commands.append("psql -h " + host + " -p " + port + " password=" + password + " -U postgres -c \"CREATE DATABASE dbaas\"")
    commands.append(psql_password + " psql -h " + host + " -p " + port + " password=" + password + " -U postgres -d dbaas -f " + filename)
    commands.append(psql_password + " psql -h " + host + " -p " + port + " password=" + password + " -U postgres -d cmdb -f " + cmdbfilename)
    commands.append(psql_password + " psql -h " + host + " -p " + port + " password=" + password + " -U postgres -d collect_alarm -f " + collectfilename)
    for item in commands:
        time.sleep(5)
        os.system(item)


if __name__ == '__main__':
    if sys.argv[1] and sys.argv[2] and sys.argv[3]:
        init_db(sys.argv[1], sys.argv[2], sys.argv[3])
    else:
        print('Please input parameter: host, port, password')
