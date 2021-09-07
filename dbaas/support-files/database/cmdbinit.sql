-- 正在导出表  public.modelobject 的数据：0 rows
/*!40000 ALTER TABLE "modelobject" DISABLE KEYS */;
INSERT INTO "public"."modelobject" ("id", "name", "code", "icon_id", "is_pre", "is_physical", "module", "init_pro_script") VALUES (5, 'K8sMySQLCluster', 'K8sMySQLCluster', 74, 't', 'f', '00000000000000000010', NULL);
INSERT INTO "public"."modelobject" ("id", "name", "code", "icon_id", "is_pre", "is_physical", "module", "init_pro_script") VALUES (6, 'K8sMySQLPod', 'K8sMySQLPod', 74, 't', 'f', '00000000000000000010', NULL);

-- 正在导出表  public.modelobjectrelation 的数据：0 rows
INSERT INTO "public"."modelobjectrelation" ("id", "tgt_object_id", "relation_id", "src_object_id", "quantity_relation_id", "desc", "collect_attribution") VALUES (1, 5, 2, 6, 3, 'K8sMySQLPod属于K8sMySQLCluster', 'attr_to_dest');


-- 正在导出表  public.modelobjectdisplay 的数据：0 rows
/*!40000 ALTER TABLE "modelobjectdisplay" DISABLE KEYS */;
INSERT INTO "public"."modelobjectdisplay" ("id", "model_id", "size", "offset_y", "position_x", "position_y") VALUES (2, 5, 55, 35, '619', '217');
INSERT INTO "public"."modelobjectdisplay" ("id", "model_id", "size", "offset_y", "position_x", "position_y") VALUES (3, 6, 55, 35, '789', '216');
/*!40000 ALTER TABLE "modelobjectdisplay" ENABLE KEYS */;


-- 正在导出表  public.module 的数据：0 rows
/*!40000 ALTER TABLE "module" DISABLE KEYS */;
INSERT INTO "public"."module" ("id", "name", "seq") VALUES (1, 'All', 1);
INSERT INTO "public"."module" ("id", "name", "seq") VALUES (2, 'DBaaS', 2);
/*!40000 ALTER TABLE "module" ENABLE KEYS */;

-- 正在导出表  public.standardfield 的数据：0 rows
/*!40000 ALTER TABLE "standardfield" DISABLE KEYS */;
DELETE FROM standardfield WHERE "id"<=4 and "id"!=1;
INSERT INTO "public"."standardfield" ("id", "name", "desc") VALUES (5, 'k8s_performance', '性能数据');
-- INSERT INTO "public"."standardfield"("name", "desc") VALUES ('k8s_pod_performance', 'pod性能数据');

-- 正在导出表  public.modelobjectattr 的数据：89 rows
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (1, 'cpuUsage', 'cpu使用率', 7, 'ft_1', 6, NULL, 'collect', 'f', 'f', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (2, 'memUsage', '内存使用率', 7, 'ft_2', 6, NULL, 'collect', 'f', 'f', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (3, 'UPTIME', 'UPTIME', 7, 'ft_3', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (4, 'InnoDBCashe', 'InnoDB缓存', 7, 'ft_5', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (5, 'QPS', 'QPS', 7, 'ft_4', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (6, 'session', '会话数', 7, 'ft_6', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (7, 'MySQL_NetworkFlow(received)', 'MySQL网络流量（接收）', 7, 'ft_11', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (8, 'MySQL_NetworkFlow(sent)', 'MySQL网络流量（发送）', 7, 'ft_11', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (9, 'Binlog_Size', 'binlog size', 7, 'ft_9', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (10, 'MySQL_Slow_Queries', 'MySQL Slow Queries', 7, 'ft_16', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (11, 'MySQL_Table_Locks', 'MySQL Table  Locks', 7, 'ft_18', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (12, 'Largest_Tables_by_Size', 'Largest Tables by Size', 7, 'ft_24', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (13, 'Largest_Tables_by_Row_Count', 'Largest Tables by Row Count', 7, 'ft_25', 6, NULL, 'collect', 'f', 't', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (14, 'colUser', '采集用户', 6, 'fs_1', 6, NULL, 'default', 't', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (15, 'version', '版本', 5, 'ft_26', 6, NULL, 'default', 'f', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (16, 'remarks', '备注', 5, 'ft_27', 6, NULL, 'default', 'f', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (17, 'volumeInfo', '存储卷信息', 5, 'ft_28', 6, NULL, 'default', 'f', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (18, 'baseInfo', '基本信息', 5, 'ft_29', 6, NULL, 'default', 'f', 'f', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (19, 'initContainerInfo', '初始化容器信息', 5, 'ft_30', 6, NULL, 'default', 'f', 'f', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (20, 'container_info', '容器信息', 5, 'ft_31', 6, NULL, 'default', 'f', 'f', 'f');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (21, 'connectString', '连接串', 6, 'fs_1', 5, NULL, 'default', 'f', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (22, 'limitMem', '最大内存', 6, 'fs_2', 5, NULL, 'default', 'f', 'f', 't');
INSERT INTO "public"."modelobjectattr" ("id", "name", "name_zh", "field_type_id", "save_field", "model_object_id", "unit_id", "collection_index", "is_necessary", "is_inoper", "is_table_show") VALUES (23, 'limitCpu', '最大CPU', 6, 'fs_3', 5, NULL, 'default', 'f', 'f', 't');




-- 正在导出表  public.modelobjectattrchart 的数据：89 rows
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (1, 1, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (2, 2, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (3, 3, 'Text', NULL, NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (4, 4, 'Text', NULL, NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (5, 5, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (6, 6, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (7, 7, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (8, 8, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (9, 9, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (10, 10, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (11, 11, 'Line', '/api/cmdb/chart/line', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (12, 12, 'Table', '/api/cmdb/chart/table', NULL);
INSERT INTO "public"."modelobjectattrchart" ("id", "modelobjectattr_id", "chart_type", "chart_url", "layout") VALUES (13, 13, 'Table', '/api/cmdb/chart/table', NULL);


-- 正在导出表  public.modelobjectattrstandardfieldrel 的数据：89 rows
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (1, 6, 1, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (2, 6, 2, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (3, 6, 3, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (4, 6, 4, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (5, 6, 5, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (6, 6, 6, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (7, 6, 7, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (8, 6, 8, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (9, 6, 9, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (10, 6, 10, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (11, 6, 11, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (12, 6, 12, 5);
INSERT INTO "public"."modelobjectattrstandardfieldrel" ("id", "model_id", "attr_id", "standardfield_id") VALUES (13, 6, 13, 5);



INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (1, 1, 'time', 'cpu', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (2, 2, 'time', 'mem', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (3, 5, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (4, 6, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (5, 7, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (6, 8, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (7, 9, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (8, 10, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (9, 11, 'time', 'value', '', NULL, NULL, NULL, '', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (10, 12, NULL, NULL, NULL, NULL, NULL, 'table_name,total_size', 'Metric,Curent', NULL);
INSERT INTO "public"."chartdisplayoption" ("id", "attr_id", "x_field", "y_field", "group_field", "unit", "measurement", "table_field", "alias", "desc") VALUES (11, 13, NULL, NULL, NULL, NULL, NULL, 'table_name,table_rows', 'Metric,Curent', NULL);
