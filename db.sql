create table "user"
(
    id            serial
        constraint user_pk
            primary key,
    name          varchar       not null,
    balance_value int default 0 not null
);

create unique index user_name_uindex
    on "user" (name);


create table service
(
    id    serial
        constraint service_pk
            primary key,
    name  varchar           not null,
    prive integer default 0 not null
);


create table "order"
(
    id         serial
        constraint order_pk
            primary key,
    user_id    integer               not null
        constraint order_user_id_fk
            references "user"
            on update cascade on delete cascade,
    service_id integer               not null
        constraint order_service_id_fk
            references service
            on update cascade on delete cascade,
    created_at timestamp             not null,
    is_paid    boolean default false not null
);

create index order_service_id_index
    on "order" (service_id);

create index order_user_id_index
    on "order" (user_id);


create table reserve
(
    id         serial
        constraint reserve_pk
            primary key,
    user_id    integer               not null
        constraint reserve_user_id_fk
            references "user"
            on update cascade on delete cascade,
    order_id   integer               not null
        constraint reserve_order_id_fk
            references "order"
            on update cascade on delete cascade,
    value      integer               not null,
    created_at timestamp             not null,
    is_debited boolean default false not null
);

create index reserve_order_id_index
    on reserve (order_id);

create index reserve_user_id_index
    on reserve (user_id);

