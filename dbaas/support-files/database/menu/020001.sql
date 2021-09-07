insert into component values ((select max(id)+1 from component), 'dbaas.cluster.apply.mysql', 1, '应用MySQL参数', 't', 'Apply MySQL Parameters');
insert into authority values ((select max(id)+1 from authority), 'dbaas.cluster.apply.mysql', 1, '/api/dbaas/cluster/configs/apply', 'POST', '应用MySQL参数', 't', null, null, 'Apply MySQL Parameters');
insert into menu_component (menu_id, component_id) values((select id from menu where name = 'dbaas.cluster'), (select id from component where name = 'dbaas.cluster.apply.mysql'));
insert into component_authority (component_id, authority_id) VALUES ((select max(id) from component), (select max(id) from authority));
insert into role_component (role_id, component_id) select 10001, id from component where name = 'dbaas.cluster.apply.mysql';