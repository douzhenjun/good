INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (1, 'script', 'MySQLdb_UPTIME', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show global status like ''uptime'';
exit;
EOF
"
)
export UPTIME=$(echo $OUTPUT|grep Uptime|awk ''{print $NF}'')
echo "mysql_global_status_uptime value=${UPTIME}" ', NULL, NULL, NULL, 'f', 's', 15, 'mysql_global_status_uptime', '', 'value', 'value', 't', 6, 3, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (2, 'script', 'MySQLdb_InnoDBCashe', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show global variables like ''innodb_buffer_pool_chunk_size'';
exit;
EOF
"
)
export INNODB=$(echo $OUTPUT|grep buffer|awk ''{print $NF}'')
echo "mysql_global_variables_innodb_buffer_pool_chunk_size value=${INNODB}"', NULL, NULL, NULL, 'f', 'B', 15, 'mysql_global_variables_innodb_buffer_pool_chunk_size', '', 'value', 'value', 't', 6, 4, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (3, 'script', 'MySQLdb_QPS', 'standard', '#!/bin/bash
export IFS=''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show  global  status like ''Question%'';
select sleep(1);
show  global  status like ''Question%'';
exit;
EOF
"
)
export QPS=$(echo $OUTPUT|grep Questi|awk ''{print $NF}''|awk ''{if(NR==1){a[NR]=$NF} else{a[NR]=$NF; print a[NR]-a[NR-1]}}'')
echo "mysql_global_status_questions value=${QPS}"', NULL, NULL, NULL, 'f', 'number', 15, 'mysql_global_status_questions', '', 'value', 'value', 't', 6, 5, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (4, 'script', 'MySQLdb_session', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show status like ''Threads_connected'';
exit;
EOF
"
)
export SESSIONS=$(echo $OUTPUT|grep Threads|awk ''{print $NF}'')
echo "mysql_global_status_threads_connected value=${SESSIONS}"', NULL, NULL, NULL, 'f', 'number', 15, 'mysql_global_status_threads_connected', '', 'value', 'value', 't', 6, 6, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (5, 'script', 'NetworkFlow(rec)', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show status like ''Bytes%'';
exit;
EOF
"
)
count=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
export variable_name=$(echo $var|awk ''{print $(NF-1)}'')
export networkflow=$(echo $var|awk ''{print $NF}'')
echo "mysql_global_status_bytes_received value=${networkflow}"
else
count=1
fi
done', NULL, NULL, NULL, 'f', 'B', 15, 'mysql_global_status_bytes_received', '', 'value', 'value', 't', 6, 7, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (6, 'script', 'NetworkFlow(sent)', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show status like ''Bytes%'';
exit;
EOF
"
)
count=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
export variable_name=$(echo $var|awk ''{print $(NF-1)}'')
export networkflow=$(echo $var|awk ''{print $NF}'')
echo "mysql_global_status_bytes_sent value=${networkflow}"
else
count=1
fi
done', NULL, NULL, NULL, 'f', 'B', 15, 'mysql_global_status_bytes_sent', '', 'value', 'value', 't', 6, 8, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (7, 'script', 'MySQLdb_Binlog_Size', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
SHOW BINARY LOGS;
exit;
EOF
"
)
count=0
filesize=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
let filesize=$(echo $var|awk ''{print $NF}'')+filesize
else
count=1
fi
done
echo "mysql_global_variables_binlog_cache_size value=${filesize}"', NULL, NULL, NULL, 'f', 'B', 15, 'mysql_global_variables_binlog_cache_size', '', 'value', 'value', 't', 6, 9, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (8, 'script', 'MySQLdb_Slow_Queries', 'standard', '#!/bin/bash
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show global status like ''slow_queries'';
exit;
EOF
"
)
export SLOWQUERIES=$(echo $OUTPUT|grep Slow_queries|awk ''{print $NF}'')
echo "mysql_global_status_slow_queries value=${SLOWQUERIES}"', NULL, NULL, NULL, 'f', 'number', 15, 'mysql_global_status_slow_queries', '', 'value', 'value', 't', 6, 10, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (9, 'script', 'MySQLdb_Table_Locks', 'standard', '#!/bin/bash
export IFS=$''\n''

export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
show open tables;
exit;
EOF
"
)
count=0
rows=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
((rows++))
else
count=1
fi
done
echo "mysql_global_status_table_locks_waited value=${rows}"', NULL, NULL, NULL, 'f', 'number', 15, 'mysql_global_status_table_locks_waited', '', 'value', 'value', 't', 6, 11, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (10, 'script', 'MySQLdb_Tables_Size', 'standard', '#!/bin/bash
flag_time=$(date "+%Y%m%d%H%M%S")
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
SELECT TABLE_NAME,concat(round((DATA_LENGTH+INDEX_LENGTH)), ''B'') as total_size FROM information_schema.TABLES order by total_size desc limit 10;
exit;
EOF
"
)
count=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
export table_name=$(echo $var|awk ''{print $(NF-1)}'')
export total_size=$(echo $var|awk ''{print $NF}'')
echo "mysql_info_schema_table_size table=\"${table_name}\",value=${total_size%B},flag_time=\"${flag_time}\""
else
count=1
fi
done', NULL, NULL, NULL, 'f', 'B', 15, 'mysql_info_schema_table_size', '', 'table,value,flag_time', 'table,value,flag_time', 't', 6, 12, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (11, 'script', 'MySQLdb_Row_Count', 'standard', '#!/bin/bash
flag_time=$(date "+%Y%m%d%H%M%S")
export IFS=$''\n''
export OUTPUT=$(su - ${username} -c "mysql -uroot -p''${passwd}'' -S ${sockpath} 2>/dev/null << EOF
SELECT TABLE_NAME,TABLE_ROWS FROM information_schema.TABLES order by TABLE_ROWS desc limit 10;
exit;
EOF
"
)
count=0
for var in ${OUTPUT[@]}
do
if [ $count == 1 ]
then
export table_name=$(echo $var|awk ''{print $(NF-1)}'')
export table_rows=$(echo $var|awk ''{print $NF}'')
echo "mysql_info_schema_table_rows table=\"${table_name}\",value=${table_rows},flag_time=\"${flag_time}\""
else
count=1
fi
done', NULL, NULL, NULL, 'f', 'number', 15, 'mysql_info_schema_table_rows', '', 'table,value,flag_time', 'table,value,flag_time', 't', 6, 13, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (12, 'script', 'cpuUsage', 'standard', '#!/bin/bash
CPUMIDDLE=`top -b -n 1| grep Cpu`
CPUIDLE=`(echo $CPUMIDDLE)| awk ''{print $8}''`
CPUUSR=`(echo $CPUMIDDLE)| awk ''{print $2}''`
CPUWAIT=`(echo $CPUMIDDLE)| awk ''{print $10}''`
CPUSYS=`(echo $CPUMIDDLE)| awk ''{print $4}''`
echo "metrics_pod_cpu cpu=${CPUWAIT}"', NULL, NULL, NULL, 'f', '%', 60, 'metrics_pod_cpu', '', 'cpu', 'cpu', 't', 6, 1, 'f', '', 'ALL', 'AA', 'DBaaS');
INSERT INTO "public"."collectitem" ("id", "method", "name", "type", "script_content", "official_item_id", "http_url", "http_header", "has_params", "result_type", "interval", "measurement", "tag", "field", "display_fields", "status", "modelobject_id", "attr_id", "is_pre", "desc", "user_tag", "org_tag", "module_name") VALUES (13, 'script', 'memUsage', 'standard', '#!/bin/bash
mem_total=`(echo $mem_middle) | awk ''{print $2}''`
echo "metrics_pod_mem mem=${mem_total}"', NULL, NULL, NULL, 'f', '%', 60, 'metrics_pod_mem', '', 'mem', 'mem', 't', 6, 2, 'f', '', 'ALL', 'AA', 'DBaaS');

INSERT INTO "public"."alarmitem" ("id", "collect_id", "is_pre", "status", "desc", "period", "interval", "adjust_period", "send_other", "name", "alarm_aim", "decline_level_times", "key", "warn_period", "user_tag") VALUES (1, 4, 'f', 't', '说明:
重复告警逻辑:
 首次告警持续90s,value的值>=500推送一条严重告警消息,500>value的值>=300推送一条普通告警信息,如连续告警,则持续告警15s,推送一条告警消息.', 6, 15, 'f', NULL, 'alarm_session', 'value', NULL, NULL, NULL, 'ALL');
INSERT INTO "public"."alarmitem" ("id", "collect_id", "is_pre", "status", "desc", "period", "interval", "adjust_period", "send_other", "name", "alarm_aim", "decline_level_times", "key", "warn_period", "user_tag") VALUES (2, 12, 'f', 't', '说明:
重复告警逻辑:
 首次告警持续180s,cpu的值>=80推送一条严重告警消息,80>cpu的值>=60推送一条普通告警信息,如连续告警,则持续告警60s,推送一条告警消息.', 3, 60, 'f', NULL, 'alarm_cpu', 'cpu', NULL, NULL, NULL, 'ALL');
INSERT INTO "public"."alarmitem" ("id", "collect_id", "is_pre", "status", "desc", "period", "interval", "adjust_period", "send_other", "name", "alarm_aim", "decline_level_times", "key", "warn_period", "user_tag") VALUES (3, 13, 'f', 't', '说明:
重复告警逻辑:
 首次告警持续180s,mem的值>=80推送一条严重告警消息,80>mem的值>=60推送一条普通告警信息,如连续告警,则持续告警60s,推送一条告警消息.', 3, 60, 'f', NULL, 'alarm_mem', 'mem', NULL, NULL, NULL, 'ALL');


INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (1, 2, 2, '60', '>=', '80>x>=60');
INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (2, 2, 1, '80', '>=', 'x>=80');
INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (3, 3, 2, '60', '>=', '80>x>=60');
INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (4, 3, 1, '80', '>=', 'x>=80');
INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (5, 1, 2, '300', '>=', '500>x>=300');
INSERT INTO "public"."alarmattr" ("id", "alarm_id", "level_id", "value", "judge_condition", "formula") VALUES (6, 1, 1, '500', '>=', 'x>=500');
