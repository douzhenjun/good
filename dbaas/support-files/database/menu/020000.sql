INSERT INTO menu select 101, 'dbaas.backup','backup','备份恢复',id,'t','Backup',NULL,NULL from menu where name = 'dbaas';

INSERT INTO component VALUES (1001, 'dbaas.backup.list', 1, '备份列表', 't', 'Backup List');
INSERT INTO component VALUES (1003, 'dbaas.backup.cycle.delete', 1, '备份周期删除', 't', 'Delete Backup Cycle');
INSERT INTO component VALUES (1004, 'dbaas.backup.add', 1, '备份新增', 't', 'Add Backup');
INSERT INTO component VALUES (1005, 'dbaas.backup.storage.type', 1, '备份存储类型列表', 't', 'Backup Storage Type List');
INSERT INTO component VALUES (1006, 'dbaas.backup.recovery', 1, '恢复实例', 't', 'Recovery Backup');
INSERT INTO component VALUES (1007, 'dbaas.backup.delete', 1, '备份删除', 't', 'Delete Backup');
INSERT INTO component VALUES (1008, 'dbaas.backup.log', 1, '备份日志', 't', 'Backup Log');
INSERT INTO component VALUES (1009, 'dbaas.backup.event', 1, '备份事件', 't', 'Backup Event');
INSERT INTO component VALUES (1010, 'dbaas.backup.list', 1, '备份列表', 't', 'Backup List');
INSERT INTO component VALUES (1011, 'dbaas.backup.add', 1, '备份新增', 't', 'Add Backup');
INSERT INTO component VALUES (1012, 'dbaas.backup.recovery', 1, '恢复实例', 't', 'Recovery Backup');
INSERT INTO component VALUES (1013, 'dbaas.backup.delete', 1, '备份删除', 't', 'Delete Backup');
INSERT INTO component VALUES (1014, 'dbaas.backup.log', 1, '备份日志', 't', 'Backup Log');
INSERT INTO component VALUES (1015, 'dbaas.backup.event', 1, '备份事件', 't', 'Backup Event');
INSERT INTO component VALUES (1016, 'dbaas.backup.storage.list', 1, '备份存储列表', 't', 'Backup Storage List');
INSERT INTO component VALUES (1017, 'dbaas.backup.storage.type', 1, '备份存储类型列表', 't', 'Backup Storage Type List');
INSERT INTO component VALUES (1018, 'dbaas.backup.storage.add', 1, '备份存储新增', 't', 'Add Backup Storage');
INSERT INTO component VALUES (1019, 'dbaas.backup.storage.delete', 1, '备份存储删除', 't', 'Delete Backup Storage');
INSERT INTO component VALUES (1020, 'dbaas.backup.storage.user', 3, '备份存储-分配租户', 't', 'Assign User');
INSERT INTO component VALUES (1022, 'dbaas.backup.cycle.delete', 1, '备份周期删除', 't', 'Delete Backup Cycle');

INSERT INTO authority VALUES (1001, 'dbaas.backup.list', 1, '/api/dbaas/backup/list', 'GET', '备份列表', 't', NULL, NULL, 'Backup List');
INSERT INTO authority VALUES (1003, 'dbaas.backup.cycle.delete', 1, '/api/dbaas/backup/cycle/delete', 'POST', '备份周期删除', 't', NULL, NULL, 'Delete Backup Cycle');
INSERT INTO authority VALUES (1004, 'dbaas.backup.add', 1, '/api/dbaas/backup/create', 'POST', '备份新增', 't', NULL, NULL, 'Add Backup');
INSERT INTO authority VALUES (1005, 'dbaas.backup.storage.type', 1, '/api/dbaas/backup/storage/type', 'GET', '备份存储类型列表', 't', NULL, NULL, 'Backup Storage Type List');
INSERT INTO authority VALUES (1006, 'dbaas.backup.recovery', 1, '/api/dbaas/backup/recovery', 'POST', '恢复实例', 't', NULL, NULL, 'Recovery Backup');
INSERT INTO authority VALUES (1007, 'dbaas.backup.delete', 1, '/api/dbaas/backup/delete', 'POST', '备份删除', 't', NULL, NULL, 'Delete Backup');
INSERT INTO authority VALUES (1008, 'dbaas.backup.log', 1, '/api/dbaas/backup/log', 'GET', '备份日志', 't', NULL, NULL, 'Backup Log');
INSERT INTO authority VALUES (1009, 'dbaas.backup.event', 1, '/api/dbaas/backup/event', 'GET', '备份事件', 't', NULL, NULL, 'Backup Event');
INSERT INTO authority VALUES (1010, 'dbaas.backup.list', 1, '/api/dbaas/backup/list', 'GET', '备份列表', 't', NULL, NULL, 'Backup List');
INSERT INTO authority VALUES (1011, 'dbaas.backup.add', 1, '/api/dbaas/backup/create', 'POST', '备份新增', 't', NULL, NULL, 'Add Backup');
INSERT INTO authority VALUES (1012, 'dbaas.backup.recovery', 1, '/api/dbaas/backup/recovery', 'POST', '恢复实例', 't', NULL, NULL, 'Recovery Backup');
INSERT INTO authority VALUES (1013, 'dbaas.backup.delete', 1, '/api/dbaas/backup/delete', 'POST', '备份删除', 't', NULL, NULL, 'Delete Backup');
INSERT INTO authority VALUES (1014, 'dbaas.backup.log', 1, '/api/dbaas/backup/log', 'GET', '备份日志', 't', NULL, NULL, 'Backup Log');
INSERT INTO authority VALUES (1015, 'dbaas.backup.event', 1, '/api/dbaas/backup/event', 'GET', '备份事件', 't', NULL, NULL, 'Backup Event');
INSERT INTO authority VALUES (1016, 'dbaas.backup.storage.list', 1, '/api/dbaas/backup/storage/list', 'GET', '备份存储列表', 't', NULL, NULL, 'Backup Storage List');
INSERT INTO authority VALUES (1017, 'dbaas.backup.storage.type', 1, '/api/dbaas/backup/storage/type', 'GET', '备份存储类型列表', 't', NULL, NULL, 'Backup Storage Type List');
INSERT INTO authority VALUES (1018, 'dbaas.backup.storage.add', 1, '/api/dbaas/backup/storage/create', 'POST', '备份存储新增', 't', NULL, NULL, 'Add Backup Storage');
INSERT INTO authority VALUES (1019, 'dbaas.backup.storage.delete', 1, '/api/dbaas/backup/storage/delete', 'POST', '备份存储删除', 't', NULL, NULL, 'Delete Backup Storage');
INSERT INTO authority VALUES (1020, 'dbaas.backup.storage.user', 3, '/api/dbaas/backup/storage/user', 'POST', '备份存储-分配租户', 't', NULL, NULL, 'Assign User');
INSERT INTO authority VALUES (1022, 'dbaas.backup.cycle.delete', 1, '/api/dbaas/backup/cycle/delete', 'POST', '备份周期删除', 't', NULL, NULL, 'Delete Backup Cycle');

INSERT INTO menu_component (menu_id, component_id) select m.id, c.id from component c,menu m where c.id BETWEEN 1001 and 1009 and m.name ='dbaas.cluster' order by 1;
INSERT INTO menu_component (menu_id, component_id) select m.id, c.id from component c,menu m where c.id BETWEEN 1010 and 1022 and m.name ='dbaas.backup' order by 1;

INSERT INTO component_authority (component_id, authority_id) VALUES (1001, 1001);
INSERT INTO component_authority (component_id, authority_id) VALUES (1003, 1003);
INSERT INTO component_authority (component_id, authority_id) VALUES (1004, 1004);
INSERT INTO component_authority (component_id, authority_id) VALUES (1005, 1005);
INSERT INTO component_authority (component_id, authority_id) VALUES (1006, 1006);
INSERT INTO component_authority (component_id, authority_id) VALUES (1007, 1007);
INSERT INTO component_authority (component_id, authority_id) VALUES (1008, 1008);
INSERT INTO component_authority (component_id, authority_id) VALUES (1009, 1009);
INSERT INTO component_authority (component_id, authority_id) VALUES (1010, 1010);
INSERT INTO component_authority (component_id, authority_id) VALUES (1011, 1011);
INSERT INTO component_authority (component_id, authority_id) VALUES (1012, 1012);
INSERT INTO component_authority (component_id, authority_id) VALUES (1013, 1013);
INSERT INTO component_authority (component_id, authority_id) VALUES (1014, 1014);
INSERT INTO component_authority (component_id, authority_id) VALUES (1015, 1015);
INSERT INTO component_authority (component_id, authority_id) VALUES (1016, 1016);
INSERT INTO component_authority (component_id, authority_id) VALUES (1017, 1017);
INSERT INTO component_authority (component_id, authority_id) VALUES (1018, 1018);
INSERT INTO component_authority (component_id, authority_id) VALUES (1019, 1019);
INSERT INTO component_authority (component_id, authority_id) VALUES (1020, 1020);
INSERT INTO component_authority (component_id, authority_id) VALUES (1022, 1022);

insert into role_component (role_id, component_id) select 10001, id from component where id BETWEEN 1001 and 1022;

INSERT INTO component VALUES (1023, 'dbaas.cluster.param.list', 1, '实例详情页-参数列表', 't', 'Param List');
INSERT INTO component VALUES (1024, 'dbaas.cluster.parameter.edit', 3, '实例详情页-参数修改', 't', 'Modify Param');

INSERT INTO authority VALUES (1023, 'dbaas.cluster.param.list', 1, '/api/dbaas/cluster/param/list', 'GET', '实例详情页-参数列表', 't', NULL, NULL, 'Param List');
INSERT INTO authority VALUES (1024, 'dbaas.cluster.parameter.edit', 3, '/api/dbaas/cluster/parameter/edit', 'POST', '实例详情页-参数修改', 't', NULL, NULL, 'Modify Param');


INSERT INTO menu_component (menu_id, component_id) select m.id, c.id from component c,menu m where c.id in (1023,1024) and m.name ='dbaas.cluster';

INSERT INTO component_authority (component_id, authority_id) VALUES (1023, 1023);
INSERT INTO component_authority (component_id, authority_id) VALUES (1024, 1024);

insert into role_component (role_id, component_id) select 10001, id from component where id in (1023,1024);

INSERT INTO component VALUES (1025, 'dbaas.image.operator', 1, '应用Operator镜像', 't', 'Apply Operator Image');
INSERT INTO authority VALUES (1025, 'dbaas.image.operator', 1, '/api/dbaas/host/operator/image', 'POST', '应用Operator镜像', 't', NULL, NULL, 'Apply Operator Image');
INSERT INTO menu_component (menu_id, component_id) select m.id, c.id from component c,menu m where c.id in (1025) and m.name ='dbaas.image';
INSERT INTO component_authority (component_id, authority_id) VALUES (1025, 1025);
insert into role_component (role_id, component_id) select 10001, id from component where id in (1025);

INSERT INTO component VALUES (1026, 'dbaas.cluster.enable', 1, '实例启用', 't', 'Enable Cluster');
INSERT INTO component VALUES (1027, 'dbaas.cluster.disable', 1, '实例禁用', 't', 'Disable Cluster');
INSERT INTO authority VALUES (1026, 'dbaas.cluster.enable', 1, '/api/dbaas/cluster/enable', 'POST', '实例启用', 't', NULL, NULL, 'Enable Cluster');
INSERT INTO authority VALUES (1027, 'dbaas.cluster.disable', 1, '/api/dbaas/cluster/disable', 'POST', '实例禁用', 't', NULL, NULL, 'Disable Cluster');
INSERT INTO menu_component (menu_id, component_id) select m.id, c.id from component c,menu m where c.id in (1026, 1027) and m.name ='dbaas.image';
INSERT INTO component_authority (component_id, authority_id) VALUES (1026, 1026);
INSERT INTO component_authority (component_id, authority_id) VALUES (1027, 1027);
insert into role_component (role_id, component_id) select 10001, id from component where id in (1026, 1027);