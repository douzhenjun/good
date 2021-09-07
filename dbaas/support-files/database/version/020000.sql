ALTER TABLE public."pod_log" ALTER COLUMN name TYPE varchar(100);
ALTER TABLE public."user" ADD backup_max integer default 0;

create table if not exists public.backup_job
(
    id             serial    not null,
    job_name       varchar(100) unique,
    pod_name       varchar(100),
    status         varchar(35),
    create_time    timestamp not null,
    duration       integer,
    backup_set     varchar(100),
    backup_task_id integer   not null,
    constraint backup_job_pkey
    primary key (id)
    );

alter table public.backup_job
    owner to postgres;

create unique index if not exists "UQE_backup_job_id"
    on public.backup_job (id);

create table if not exists public.backup_task
(
    id         serial      not null,
    type       varchar(10) not null,
    keep_copy  integer,
    crontab    varchar(20),
    storage_id integer     not null,
    cluster_id integer     not null,
    user_id    integer     not null,
    name       varchar(30),
    close      boolean default false,
    set_type   varchar(10),
    set_date   varchar(10),
    set_time   varchar(15),
    constraint backup_task_pkey
    primary key (id)
    );

alter table public.backup_task
    owner to postgres;

create unique index if not exists "UQE_backup_task_id"
    on public.backup_task (id);

create table if not exists public.backup_storage
(
    id         serial       not null,
    name       varchar(50)  not null,
    type       varchar(10)  not null,
    end_point  varchar(50)  not null,
    bucket     varchar(50)  not null,
    access_key varchar(50)  not null,
    secret_key varchar(100) not null,
    status     varchar(10),
    assign_all boolean default false,
    constraint backup_storage_pkey
    primary key (id)
    );

alter table public.backup_storage
    owner to postgres;

create unique index if not exists "UQE_backup_storage_id"
    on public.backup_storage (id);

create unique index if not exists "UQE_backup_storage_name"
    on public.backup_storage (name);

create table if not exists public.backup_storage_user
(
    id         serial  not null,
    storage_id integer not null,
    user_id    integer not null,
    constraint backup_storage_user_pkey
    primary key (id)
    );

alter table public.backup_storage_user
    owner to postgres;

create unique index if not exists "UQE_backup_storage_user_id"
    on public.backup_storage_user (id);

create table if not exists public.backup_storage_type
(
    id   serial      not null,
    type varchar(20) not null,
    constraint backup_storage_type_pkey
    primary key (id)
    );

alter table public.backup_storage_type
    owner to postgres;

create unique index if not exists "UQE_backup_storage_type_id"
    on public.backup_storage_type (id);

INSERT INTO public.sysparameter (id, param_key, param_value, default_value, is_modifiable) VALUES (11, 'operator_tz', 'Asia/Shanghai', 'Asia/Shanghai', true) ON CONFLICT(id) DO NOTHING;
INSERT INTO public.backup_storage_type (id, type) VALUES (1, 'Ceph') ON CONFLICT(id) DO NOTHING;

update defaultparameters set parameter_value='25' where parameter_name = 'innodb_buffer_pool_dump_pct';
delete from defaultparameters where parameter_name = 'innodb_buffer_pool_instances';
insert into sysparameter values ((select max(id)+1  from sysparameter), 'mp_innodb_buffer_pool_instances', 'min({m}/1073741824, 8)', 'min({m}/1073741824, 8)', 't');
insert into sysparameter values ((select max(id)+1  from sysparameter), 'mp_innodb_buffer_pool_size', '{m}*0.75', '{m}*0.75', 't');
insert into defaultparameters values ((select max(id)+1 from defaultparameters), 'innodb_disable_sort_file_cache', '"NO"', 4);
update defaultparameters set parameter_value='0' where parameter_name = 'innodb_flush_neighbors';
update defaultparameters set parameter_value='20000' where parameter_name='innodb_io_capacity';
update defaultparameters set parameter_value='40000' where parameter_name='innodb_io_capacity_max';
update defaultparameters set parameter_value='1024' where parameter_name='innodb_lru_scan_depth';
update defaultparameters set parameter_value='1' where parameter_name='innodb_page_cleaners';
insert into defaultparameters values ((select max(id)+1 from defaultparameters), 'innodb_purge_threads', '1', 4);
update defaultparameters set parameter_value='6' where parameter_name='innodb_spin_wait_delay';
update defaultparameters set parameter_value='20' where parameter_name='innodb_stats_persistent_sample_pages';
update defaultparameters set parameter_value='0' where parameter_name='innodb_undo_log_truncate';
update defaultparameters set parameter_value='0' where parameter_name='innodb_undo_tablespaces';

delete from defaultparameters where parameter_name='innodb_write_io_threads';
delete from defaultparameters where parameter_name='innodb_read_io_threads';
insert into sysparameter values ((select max(id)+1  from sysparameter), 'mp_innodb_write_io_threads', 'max(4, floor({m}/1073741824/64)*8)', 'max(4, floor({m}/1073741824/64)*8)', 't');
insert into sysparameter values ((select max(id)+1  from sysparameter), 'mp_innodb_read_io_threads', 'max(4, floor({m}/1073741824/64)*8)', 'max(4, floor({m}/1073741824/64)*8)', 't');

create table if not exists public.combo
(
    id         serial            not null,
    name       varchar(20),
    cpu        integer default 0 not null,
    mem        integer default 0 not null,
    storage    integer default 0 not null,
    tags       varchar(100),
    remark     varchar(255),
    read_bps   bigint            not null,
    write_bps  bigint            not null,
    read_iops  bigint            not null,
    write_iops bigint            not null,
    assign_all boolean default false,
    constraint combo_pkey
        primary key (id)
);

alter table public.combo
    owner to postgres;

create unique index if not exists "UQE_combo_name"
    on public.combo (name);

create unique index if not exists "UQE_combo_id"
    on public.combo (id);

create table if not exists public.combo_user
(
    id       serial  not null,
    combo_id integer not null,
    user_id  integer not null,
    constraint combo_user_pkey
        primary key (id)
);

alter table public.combo_user
    owner to postgres;

create unique index if not exists "UQE_combo_user_id"
    on public.combo_user (id);

create table if not exists public.qos
(
    id         serial  not null,
    cluster_id integer not null,
    read_bps   bigint  not null,
    write_bps  bigint  not null,
    read_iops  bigint  not null,
    write_iops bigint  not null,
    constraint qos_pkey
        primary key (id)
);

alter table public.qos
    owner to postgres;

create unique index if not exists "UQE_qos_id"
    on public.qos (id);
