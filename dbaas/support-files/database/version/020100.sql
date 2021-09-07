alter table combo add copy integer default 1;
alter table sc_user drop column remark;
alter table cluster_instance add combo_id integer default 0;

create table if not exists public.combo_tag
(
    id     serial,
    name   varchar(20)           not null,
    preset boolean default false not null,
    constraint combo_tag_pkey
        primary key (id)
);

alter table public.combo_tag
    owner to postgres;

create unique index if not exists "UQE_combo_tag_id"
    on public.combo_tag (id);

alter table combo_tag
    add constraint unique_name unique (name);

INSERT INTO public.combo_tag (id, name, preset) VALUES (1, '高性能', true);
INSERT INTO public.combo_tag (id, name, preset) VALUES (2, '高可靠性', true);
INSERT INTO public.combo_tag (id, name, preset) VALUES (3, '个人/实验环境', true);
INSERT INTO public.combo_tag (id, name, preset) VALUES (4, '中小企业用户', true);
INSERT INTO public.combo_tag (id, name, preset) VALUES (5, '大型企业用户', true);
select setval('combo_tag_id_seq', 100)

alter table sc
    add assign_all bool default false not null;
