delete from role_component where id = 10001 and component_id in (select c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.parameter'));
delete from authority where id in (select c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.parameter'));
delete from component where id in (select c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.parameter'));
delete from menu_component where menu_id in (select menu.id from menu where name = 'dbaas.parameter');
delete from menu where name = 'dbaas.parameter';

INSERT INTO menu select 102, 'dbaas.system', 'system', '系统', id, 't', 'System', NULL, NULL from menu where name = 'dbaas';
INSERT INTO menu select 103, 'dbaas.system.parameter', 'parameter', '系统参数', id, 't', 'Parameter', NULL, NULL from menu where name = 'dbaas.system';
INSERT INTO menu select 104, 'dbaas.system.combo', 'combo', '套餐管理', id, 't', 'Package Management', NULL, NULL from menu where name = 'dbaas.system';

INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.parameter.list', 1, '系统参数列表', 't', 'Parameter List');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.parameter.update', 1, '系统参数修改', 't', 'Update Parameter');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.parameter.reset', 1, '还原默认值', 't', 'Reset Parameter');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.parameter.list', 1, '/api/dbaas/parameter/list', 'GET', '系统参数列表', 't', NULL, NULL, 'Parameter List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.parameter.update', 1, '/api/dbaas/parameter/update', 'POST', '系统参数修改', 't', NULL, NULL, 'Update Parameter');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.parameter.reset', 1, '/api/dbaas/parameter/reset', 'POST', '还原默认值', 't', NULL, NULL, 'Reset Parameter');
INSERT INTO menu_component (menu_id, component_id) select (select id from menu where name = 'dbaas.system.parameter'), c.id from component c where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset');
INSERT INTO component_authority (component_id, authority_id) select c.id, c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.parameter');
-- insert into role_component (role_id, component_id) select 10001,c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.parameter');

INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.list', 1, '套餐列表', 't', 'Package List');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.add', 3, '套餐-新增', 't', 'Add');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.edit', 3, '套餐-修改', 't', 'Update');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.delete', 3, '套餐-删除', 't', 'Delete');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.user.list', 1, '租户列表', 't', 'User List');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.user', 3, '分配租户', 't', 'Config User');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.list', 1, '/api/dbaas/combo/list', 'GET', '套餐列表', 't', NULL, NULL, 'Package List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.add', 3, '/api/dbaas/combo/add', 'POST', '套餐-新增', 't', NULL, NULL, 'Add');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.edit', 3, '/api/dbaas/combo/edit', 'POST', '套餐-修改', 't', NULL, NULL, 'Update');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.delete', 1, '/api/dbaas/combo/delete', 'POST', '套餐-删除', 't', NULL, NULL, 'Delete');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.user.list', 1, '/api/dbaas/user/list', 'GET', '租户列表', 't', NULL, NULL, 'User List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.user', 3, '/api/dbaas/combo/user', 'POST', '分配租户', 't', NULL, NULL, 'Config User');

INSERT INTO menu_component (menu_id, component_id) select (select id from menu where name = 'dbaas.system.combo'), c.id from component c where c.name in ('dbaas.combo.list', 'dbaas.combo.add', 'dbaas.combo.edit', 'dbaas.combo.delete', 'dbaas.user.list', 'dbaas.combo.user');
INSERT INTO component_authority (component_id, authority_id) select c.id, c.id from component c, menu_component mc where c.name in ('dbaas.combo.list', 'dbaas.combo.add', 'dbaas.combo.edit', 'dbaas.combo.delete', 'dbaas.user.list', 'dbaas.combo.user') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.combo');
-- insert into role_component (role_id, component_id) select 10001, c.id from component c, menu_component mc where c.name in ('dbaas.combo.list', 'dbaas.combo.add', 'dbaas.combo.edit', 'dbaas.combo.delete', 'dbaas.user.list', 'dbaas.combo.user') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.combo');

INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.tag.list', 1, '特性列表', 't', 'Tags List');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.tag.add', 1, '特性-新增', 't', 'Tags Add');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.tag.delete', 1, '特性-删除"', 't', 'Tags Delete');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.tag.list', 1, '/api/dbaas/combo/tag/list', 'GET', '特性列表', 't', NULL, NULL, 'Tags List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.tag.add', 1, '/api/dbaas/combo/tag/add', 'POST', '特性-新增', 't', NULL, NULL, 'Tags Add');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.tag.delete', 1, '/api/dbaas/combo/tag/delete', 'POST', '特性-删除', 't', NULL, NULL, 'Tags Delete');
INSERT INTO menu_component (menu_id, component_id) select (select id from menu where name = 'dbaas.system.combo'), c.id from component c where c.name in ('dbaas.combo.tag.list', 'dbaas.combo.tag.add', 'dbaas.combo.tag.delete');
INSERT INTO component_authority (component_id, authority_id) select c.id, c.id from component c, menu_component mc where c.name in ('dbaas.combo.tag.list', 'dbaas.combo.tag.add', 'dbaas.combo.tag.delete') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.combo');
-- insert into role_component (role_id, component_id) select 10001,c.id from component c, menu_component mc where c.name in ('dbaas.parameter.list', 'dbaas.parameter.update', 'dbaas.parameter.reset') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.system.combo');

INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.cluster.nodeport.check', 1, 'nodeport 可用性检测"', 't', 'NodePort check');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.cluster.nodeport.check', 1, '/api/dbaas/cluster/nodeport/check', 'POST', 'nodeport 可用性检测', 't', NULL, NULL, 'NodePort check');
INSERT INTO menu_component (menu_id, component_id) select (select id from menu where name = 'dbaas.cluster'), c.id from component c where c.name = 'dbaas.cluster.nodeport.check';
INSERT INTO component_authority (component_id, authority_id) select c.id, c.id from component c, menu_component mc where c.name = 'dbaas.cluster.nodeport.check' and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.cluster');
-- insert into role_component (role_id, component_id) select 10001,c.id from component c, menu_component mc where c.name = 'dbaas.cluster.nodeport.check' and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.cluster');

INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.parameter.list', 1, '系统参数列表', 't', 'Parameter List');
INSERT INTO component VALUES ((select max(id)+1 from component), 'dbaas.combo.list', 1, '套餐列表', 't', 'Package List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.parameter.list', 1, '/api/dbaas/parameter/list', 'GET', '系统参数列表', 't', NULL, NULL, 'Parameter List');
INSERT INTO authority VALUES ((select max(id)+1 from authority), 'dbaas.combo.list', 1, '/api/dbaas/combo/list', 'GET', '套餐列表', 't', NULL, NULL, 'Package List');
INSERT INTO menu_component (menu_id, component_id) select (select id from menu where name = 'dbaas.cluster'), c.id from component c order by c.id desc limit 2;
INSERT INTO component_authority (component_id, authority_id) values ((select max(c.id) from component c where c.name = 'dbaas.parameter.list'),(select max(a.id) from authority a where a.name = 'dbaas.parameter.list'));
INSERT INTO component_authority (component_id, authority_id) values ((select max(c.id) from component c where c.name = 'dbaas.combo.list'),(select max(a.id) from authority a where a.name = 'dbaas.combo.list'));
insert into role_component (role_id, component_id) select 10001,c.id from component c, menu_component mc where c.name in ('dbaas.combo.list', 'dbaas.parameter.list') and mc.component_id = c.id and mc.menu_id = (select id from menu where name = 'dbaas.cluster');