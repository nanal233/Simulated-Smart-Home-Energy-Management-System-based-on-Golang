CREATE DATABASE `project20240227` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci */ /*!80016 DEFAULT ENCRYPTION='N' */

create table client
(
    id         varchar(255)                              not null comment '客户端编号'
        primary key,
    name       varchar(255)                              not null comment '名称',
    type       int                                       not null comment '客户端类型',
    created_at timestamp(3) default CURRENT_TIMESTAMP(3) not null comment '加入时间',
    updated_at timestamp(3) default CURRENT_TIMESTAMP(3) not null on update CURRENT_TIMESTAMP(3) comment '上次更新时间'
)
    comment '客户端';

create index client_name_index
    on client (name);

create table client_activity
(
    id         bigint auto_increment comment '编号'
        primary key,
    client_id  varchar(255)                              not null comment '客户端ID',
    status     tinyint                                   not null comment '状态',
    created_at timestamp(3) default CURRENT_TIMESTAMP(3) not null comment '变动时间',
    constraint client_activity_client_id_fk
        foreign key (client_id) references client (id)
            on update cascade on delete cascade
)
    comment '客户端活跃历史';

create table client_command_execution
(
    id         bigint auto_increment comment '编号'
        primary key,
    client_id  varchar(255)                              not null comment '客户端ID',
    code       int                                       not null comment '命令代码',
    data       text                                      not null comment '命令内容',
    sent_at    timestamp(3)                              not null comment '发送时间',
    created_at timestamp(3) default CURRENT_TIMESTAMP(3) not null comment '创建时间',
    constraint client_command_execution_client_id_fk
        foreign key (client_id) references client (id)
            on update cascade on delete cascade
)
    comment '客户端命令执行历史';

create index client_command_execution_code_index
    on client_command_execution (code);

create table client_consumption
(
    id          bigint auto_increment comment '编号'
        primary key,
    client_id   varchar(255)                                not null comment '客户端编号',
    consumption float unsigned default '0'                  not null comment '功耗',
    recorded_at timestamp(3)                                not null comment '记录功耗时间',
    created_at  timestamp(3)   default CURRENT_TIMESTAMP(3) not null on update CURRENT_TIMESTAMP(3) comment '保存时间',
    constraint client_consumption_client_id_fk
        foreign key (client_id) references client (id)
            on update cascade on delete cascade
)
    comment '客户端功耗记录';

create index client_consumption_client_id_index
    on client_consumption (client_id);

create index client_consumption_recorded_at_index
    on client_consumption (recorded_at);

create table power_mode
(
    id         bigint auto_increment comment '编号'
        primary key,
    name       varchar(255)                              not null comment '名称',
    created_at timestamp(3) default CURRENT_TIMESTAMP(3) not null comment '创建时间',
    updated_at timestamp(3) default CURRENT_TIMESTAMP(3) not null on update CURRENT_TIMESTAMP(3) comment '最后更新时间',
    constraint power_mode_pk
        unique (name)
)
    comment '能耗模式';

create table client_prepared_command
(
    id            bigint auto_increment comment '编号'
        primary key,
    client_id     varchar(255)                              not null comment '客户端',
    power_mode_id bigint                                    not null comment '能耗模式ID',
    code          int                                       not null comment '命令代码',
    data          text                                      not null comment '命令内容',
    created_at    timestamp(3) default CURRENT_TIMESTAMP(3) null comment '创建时间',
    constraint client_prepared_command_pk
        unique (power_mode_id, client_id),
    constraint client_prepared_command_client_id_fk
        foreign key (client_id) references client (id)
            on update cascade on delete cascade,
    constraint client_prepared_command_power_mode_id_fk
        foreign key (power_mode_id) references power_mode (id)
            on update cascade
)
    comment '客户端预制命令';

create index client_prepared_command_client_id_index
    on client_prepared_command (client_id);

create table power_mode_execution
(
    id            bigint auto_increment comment '编号'
        primary key,
    power_mode_id bigint                                    not null comment '能耗模式编号',
    created_at    timestamp(3) default CURRENT_TIMESTAMP(3) not null comment '创建时间',
    constraint power_mode_execution_power_mode_id_fk
        foreign key (power_mode_id) references power_mode (id)
            on update cascade
)
    comment '能耗模式执行历史';


