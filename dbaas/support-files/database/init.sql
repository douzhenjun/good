create table if not exists public.cluster_instance
(
    id                   serial      not null,
    name                 varchar(40) not null,
    k8s_name             varchar(40),
    secret_name          varchar(40),
    connect_string       varchar(40),
    status               varchar(40),
    user_id              integer,
    image_id             integer,
    limit_mem            integer,
    limit_cpu            integer,
    remark               text,
    pod_status           text,
    yaml_text            text,
    replicas             varchar(40),
    operator             varchar(40),
    org_tag              varchar(10),
    user_tag             varchar(10),
    storage              integer,
    sc_name              varchar(100),
    inner_connect_string varchar(40),
    console_port         varchar(40),
    actual_replicas      varchar(40),
    master               varchar(100),
    is_deploy            boolean default false,
    secret               varchar(50),
    deleted_at           timestamp(6),
    pv_id                integer,
    constraint cluster_instance_pkey
        primary key (id)
);

alter table public.cluster_instance
    owner to postgres;

create unique index if not exists "UQE_cluster_instance_id"
    on public.cluster_instance (id);

create index if not exists "UQE_cluster_instance_k8s_name"
    on public.cluster_instance (k8s_name);

create index if not exists "UQE_cluster_instance_name"
    on public.cluster_instance (name);

create table if not exists public.clusterparameters
(
    id              serial not null,
    parameter_name  varchar(100),
    parameter_value varchar(200),
    cluster_id      integer,
    constraint clusterparameters_pkey
        primary key (id)
);

alter table public.clusterparameters
    owner to postgres;

create unique index if not exists "UQE_clusterparameters_id"
    on public.clusterparameters (id);

create table if not exists public.image_type
(
    type     varchar(100),
    category varchar(100),
    id       integer not null,
    constraint image_type_pkey
        primary key (id)
);

alter table public.image_type
    owner to postgres;

create table if not exists public.defaultparameters
(
    id              serial not null,
    parameter_name  varchar(100),
    parameter_value varchar(200),
    image_type_id   integer,
    constraint pk_defaultparameters
        primary key (id),
    constraint fk_defaultp_reference_images
        foreign key (image_type_id) references public.image_type
            on update cascade on delete cascade
);

alter table public.defaultparameters
    owner to postgres;

create table if not exists public.images
(
    id            serial       not null,
    image_name    varchar(100) not null,
    version       varchar(100) not null,
    status        varchar(100) not null,
    description   varchar(300),
    image_type_id integer,
    constraint pk_images
        primary key (id)
);

alter table public.images
    owner to postgres;

create unique index if not exists "UQE_images_unique_name_version"
    on public.images (image_name, version);

create table if not exists public.instance
(
    id             serial not null,
    name           varchar(100),
    domain_name    varchar(40),
    cluster_id     integer,
    version        varchar(20),
    status         varchar(20),
    volume         text,
    role           varchar(20),
    base_info      text,
    init_container text,
    container_info text,
    deleted_at     timestamp(6),
    constraint instance_pkey
        primary key (id)
);

alter table public.instance
    owner to postgres;

create unique index if not exists "UQE_instance_id"
    on public.instance (id);

create index if not exists "UQE_instance_name"
    on public.instance (name);

create table if not exists public.mysql_operator
(
    id                serial not null,
    name              varchar(40),
    ready             varchar(20),
    status            varchar(20),
    deployment_status varchar(40),
    org_tag           varchar(10),
    user_tag          varchar(10),
    node_name         varchar(40),
    ip                varchar(40),
    service_ip        varchar(40),
    replicas          varchar(40),
    container_status  varchar(20),
    constraint pk_mysql_operator
        primary key (id)
);

alter table public.mysql_operator
    owner to postgres;

create unique index if not exists "UQE_mysql_operator_name"
    on public.mysql_operator (name);

create table if not exists public.initinfo
(
    id       serial      not null,
    name     varchar(20) not null,
    message  varchar     not null,
    isaccess varchar(20) not null,
    isdeploy varchar(20),
    constraint pk_initinfo
        primary key (id)
);

alter table public.initinfo
    owner to postgres;

create table if not exists public.node
(
    id        serial not null,
    node_name varchar(30),
    status    varchar(30),
    age       varchar(20),
    label     text,
    org_tag   varchar(10),
    user_tag  varchar(10),
    constraint pk_node
        primary key (id)
);

alter table public.node
    owner to postgres;

create unique index if not exists unq_host_name
    on public.node (node_name);

create table if not exists public.oper_log
(
    id          serial not null,
    oper_date   timestamp(6),
    oper_people varchar(40),
    content     text,
    type_level  varchar(40),
    log_source  varchar(100),
    constraint oper_log_pkey
        primary key (id)
);

alter table public.oper_log
    owner to postgres;

create unique index if not exists "UQE_oper_log_id"
    on public.oper_log (id);

create table if not exists public.pod_log
(
    id        serial       not null,
    oper_date timestamp(6) not null,
    message   text,
    "from"    varchar(60),
    name      varchar(60),
    reason    varchar(60),
    type      varchar(60)  not null,
    constraint pod_log_pkey
        primary key (id)
);

alter table public.pod_log
    owner to postgres;

create unique index if not exists "UQE_pod_log_id"
    on public.pod_log (id);

create unique index if not exists "UQE_pod_log_name"
    on public.pod_log (name);

create table if not exists public.regular_role
(
    id       serial       not null,
    role     varchar(16),
    res_type varchar(20)  not null,
    res_name varchar(100) not null,
    level    integer      not null,
    constraint pk_regular_role
        primary key (id)
);

alter table public.regular_role
    owner to postgres;

create table if not exists public.sc
(
    id             serial not null,
    name           varchar(40),
    describe       varchar(100),
    reclaim_policy varchar(200),
    user_tag       varchar(10),
    org_tag        varchar(10),
    sc_type        varchar(60),
    node_num       integer,
    constraint pk_sc
        primary key (id)
);

alter table public.sc
    owner to postgres;

create table if not exists public.persistent_volume
(
    id             serial not null,
    sc_id          integer,
    capacity       varchar(20),
    name           varchar(50),
    ip_addr        varchar(20),
    port           varchar(20),
    iqn            varchar(40),
    lun            integer,
    org_tag        varchar(10),
    user_tag       varchar(10),
    reclaim_policy varchar(10),
    status         varchar(10),
    pod_id         integer,
    deleted_at     timestamp(6),
    pvc_name       varchar(50),
    constraint pk_persistent_volume
        primary key (id),
    constraint fk_persiste_reference_sc
        foreign key (sc_id) references public.sc
            on update cascade on delete cascade
);

alter table public.persistent_volume
    owner to postgres;

create unique index if not exists "UQE_sc_name"
    on public.sc (name);

create table if not exists public.sc_user
(
    id      serial not null,
    sc_id   integer,
    user_id integer,
    remark  varchar(200),
    constraint sc_user_pkey
        primary key (id),
    constraint fk_sc_id
        foreign key (sc_id) references public.sc
            on update cascade on delete cascade
);

alter table public.sc_user
    owner to postgres;

create unique index if not exists "UQE_sc_user_id"
    on public.sc_user (id);

create table if not exists public.sysparameter
(
    id            serial      not null,
    param_key     varchar(64) not null,
    param_value   text        not null,
    default_value text        not null,
    is_modifiable boolean,
    constraint pk_sysparameter
        primary key (id)
);

alter table public.sysparameter
    owner to postgres;

create unique index if not exists "UQE_sysparameter_unq_name_key"
    on public.sysparameter (param_key);

create table if not exists public."user"
(
    id          serial not null,
    zdcp_id     integer,
    mem_all     bigint,
    remarks     text,
    cpu_all     integer,
    user_name   varchar(100),
    password    varchar(32),
    storage_all integer,
    user_tag    varchar(100),
    auto_create boolean default false,
    constraint pk_user
        primary key (id)
);

alter table public."user"
    owner to postgres;

create table if not exists public.user_regular
(
    id      serial not null,
    user_id integer,
    role_id integer,
    remake  varchar(100),
    constraint pk_user_regular
        primary key (id),
    constraint fk_user_reg_reference_regular_
        foreign key (role_id) references public.regular_role
            on update cascade on delete cascade,
    constraint fk_user_reg_reference_user
        foreign key (user_id) references public."user"
            on update cascade on delete cascade
);

alter table public.user_regular
    owner to postgres;

create table if not exists public.misc_config
(
    key   varchar(50) not null,
    value varchar(100),
    constraint misc_config_pkey
        primary key (key)
);

alter table public.misc_config
    owner to postgres;

create unique index if not exists "UQE_misc_config_key"
    on public.misc_config (key);

create table if not exists public.api_quota
(
    id      serial            not null,
    path    varchar(50)       not null,
    cpu     integer default 0 not null,
    memory  integer default 0 not null,
    storage integer default 0 not null,
    constraint api_quota_pkey
        primary key (id)
);

alter table public.api_quota
    owner to postgres;

create unique index if not exists "UQE_api_quota_id"
    on public.api_quota (id);

INSERT INTO public.image_type (type, category, id) VALUES ('Operator', 'Operator', 1) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.image_type (type, category, id) VALUES ('Operator', 'Orchestrator', 2) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.image_type (type, category, id) VALUES ('Operator', 'sidecar5.7', 3) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.image_type (type, category, id) VALUES ('Operator', 'sidecar8', 5) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.image_type (type, category, id) VALUES ('Mysql', '5.7', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.image_type (type, category, id) VALUES ('Mysql', '8', 6) ON CONFLICT(id) DO NOTHING;

INSERT INTO public.api_quota (id, path, cpu, memory, storage) VALUES (1, '/external/cluster/add', 80, 360, 6144) ON CONFLICT(id) DO NOTHING;

INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (1, 'disabled_storage_engines', '"MyISAM,BLACKHOLE,FEDERATED,ARCHIVE,MEMORY"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (2, 'plugin-load', 'semisync_master.so:semisync_slave.so', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (3, 'rpl-semi-sync-master-enabled', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (4, 'rpl-semi-sync-master-timeout', '1000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (5, 'rpl-semi-sync-master-wait-for-slave-count', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (6, 'rpl-semi-sync-slave-enabled', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (7, 'slave-compressed-protocol', '"OFF"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (8, 'default-storage-engine', 'INNODB', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (9, 'character-set-server', 'utf8mb4', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (10, 'collation-server', 'utf8mb4_bin', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (11, 'explicit_defaults_for_timestamp', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (12, 'log_timestamps', 'system', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (13, 'default_time_zone', 'SYSTEM', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (14, 'skip_name_resolve', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (15, 'lower_case_table_names', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (16, 'auto_increment_increment', '6', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (17, 'auto_increment_offset', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (18, 'sql_mode', 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (19, 'autocommit', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (20, 'transaction_isolation', 'READ-COMMITTED', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (21, 'max_allowed_packet', '16M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (22, 'event_scheduler', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (23, 'slave_pending_jobs_size_max', '16M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (24, 'show_compatibility_56', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (25, 'query_cache_size', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (26, 'query_cache_type', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (27, 'interactive_timeout', '1800', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (28, 'wait_timeout', '1800', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (29, 'lock_wait_timeout', '1800', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (30, 'max_connections', '3000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (31, 'max_connect_errors', '1000000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (32, 'read_buffer_size', '16M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (33, 'read_rnd_buffer_size', '32M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (34, 'sort_buffer_size', '8M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (35, 'tmp_table_size', '1024M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (36, 'max_heap_table_size', '1024M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (37, 'thread_cache_size', '64', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (38, 'binlog_format', 'row', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (39, 'binlog_rows_query_log_events', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (40, 'log_slave_updates', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (41, 'expire_logs_days', '7', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (42, 'binlog_cache_size', '65536', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (43, 'binlog_checksum', 'none', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (44, 'sync_binlog', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (45, 'log_bin_trust_function_creators', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (46, 'binlog_gtid_simple_recovery', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (47, 'max_binlog_cache_size', '4096M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (48, 'optimizer_switch', '''use_index_extensions=off''', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (49, 'general_log', '"OFF"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (50, 'slow_query_log', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (51, 'log_queries_not_using_indexes', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (52, 'long_query_time', '1.000000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (53, 'log_slow_admin_statements', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (54, 'log_slow_slave_statements', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (55, 'log_throttle_queries_not_using_indexes', '10', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (56, 'min_examined_row_limit', '100', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (57, 'gtid_mode', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (58, 'enforce_gtid_consistency', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (59, 'skip_slave_start', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (60, 'master_info_repository', 'table', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (61, 'relay_log_info_repository', 'table', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (62, 'slave_parallel_type', 'logical_clock', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (63, 'slave_parallel_workers', '16', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (64, 'relay_log_recovery', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (65, 'slave_rows_search_algorithms', '''INDEX_SCAN,HASH_SCAN''', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (66, 'slave_preserve_commit_order', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (67, 'slave_transaction_retries', '128', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (68, 'sync_relay_log', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (69, 'sync_relay_log_info', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (70, 'sync_master_info', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (71, 'binlog_group_commit_sync_delay', '20', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (72, 'binlog_group_commit_sync_no_delay_count', '5', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (73, 'default_storage_engine', 'innodb', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (74, 'default_tmp_storage_engine', 'innodb', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (75, 'innodb_data_file_path', 'ibdata1:1024M:autoextend', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (76, 'innodb_temp_data_file_path', 'ibtmp1:12M:autoextend', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (77, 'innodb_buffer_pool_filename', 'ib_buffer_pool', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (78, 'innodb_log_files_in_group', '3', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (79, 'innodb_log_file_size', '2048M', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (80, 'innodb_file_per_table', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (81, 'innodb_online_alter_log_max_size', '1G', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (82, 'innodb_open_files', '4096', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (83, 'innodb_page_size', '16k', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (84, 'innodb_thread_concurrency', '0', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (85, 'innodb_read_io_threads', '16', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (86, 'innodb_write_io_threads', '16', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (87, 'innodb_large_prefix', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (88, 'innodb_page_cleaners', '16', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (89, 'innodb_print_all_deadlocks', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (90, 'innodb_deadlock_detect', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (91, 'innodb_lock_wait_timeout', '5', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (92, 'innodb_spin_wait_delay', '128', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (93, 'innodb_autoinc_lock_mode', '2', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (94, 'innodb_io_capacity', '200', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (95, 'innodb_io_capacity_max', '2000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (96, 'innodb_lru_scan_depth', '4096', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (97, 'innodb_undo_logs', '128', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (98, 'innodb_undo_tablespaces', '3', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (99, 'innodb_log_buffer_size', '16777216', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (100, 'innodb_strict_mode', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (101, 'innodb_sort_buffer_size', '6.7108864e+07', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (102, 'innodb_undo_log_truncate', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (103, 'innodb_max_undo_log_size', '2G', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (104, 'innodb_purge_rseg_truncate_frequency', '128', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (105, 'innodb_stats_auto_recalc', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (106, 'innodb_stats_persistent', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (107, 'innodb_stats_persistent_sample_pages', '64', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (108, 'innodb_adaptive_hash_index', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (109, 'innodb_change_buffering', 'all', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (110, 'innodb_change_buffer_max_size', '24', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (111, 'innodb_flush_neighbors', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (112, 'innodb_flush_method', 'O_DIRECT', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (113, 'innodb_doublewrite', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (114, 'innodb_flush_log_at_timeout', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (115, 'innodb_flush_log_at_trx_commit', '1', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (116, 'innodb_buffer_pool_instances', '8', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (117, 'innodb_old_blocks_pct', '37', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (118, 'innodb_old_blocks_time', '1000', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (119, 'innodb_read_ahead_threshold', '56', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (120, 'innodb_random_read_ahead', '"OFF"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (121, 'innodb_buffer_pool_dump_pct', '40', 4) ON CONFLICT(id) DO NOTHING;;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (122, 'innodb_buffer_pool_dump_at_shutdown', '"ON"', 4) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.defaultparameters (id, parameter_name, parameter_value, image_type_id) VALUES (123, 'innodb_buffer_pool_load_at_startup', '"ON"', 4) ON CONFLICT(id) DO NOTHING;

INSERT INTO public.initinfo (id, name, message, isaccess, isdeploy) VALUES (1, 'skipStstus', 'skipStstus', 'False', null) ON CONFLICT(id) DO NOTHING;

INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (10, 'mp_join_buffer_size', 'min({m}/1048576*128,262144)', 'min({m}/1048576*128,262144)', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (5, 'zcm_sc_normal', 'disecsi-iscsi-rdma-ioqos', 'disecsi-iscsi-rdma-ioqos', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (6, 'zcm_sc_performance', 'unknown', 'unknown', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (7, 'mp_table_open_cache', 'min({m}/1073741824*256,2048)', 'min({m}/1073741824*256,2048)', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (1, 'kubernetes_master_address', '', '', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (8, 'mp_table_definition_cache', 'min({m}/1073741824*512,2048)', 'min({m}/1073741824*512,2048)', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (9, 'mp_table_open_cache_instances', 'min({m}/1073741824,16)', 'min({m}/1073741824,16)', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (2, 'kubernetes_config', '', '', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (3, 'kubernetes_namespace', 'default', 'default', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (4, 'harbor_address', '', '', true) ON CONFLICT(id) DO NOTHING;
