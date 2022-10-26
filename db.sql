create table "user"
(
    id            serial
        constraint user_pk
            primary key,
    name          varchar,
    balance_value double precision default 0 not null
);


create table service
(
    id    serial
        constraint service_pk
            primary key,
    name  varchar                    not null,
    price double precision default 0 not null
);


create table reserve
(
    id         serial
        constraint reserve_pk
            primary key,
    user_id    integer               not null
        constraint reserve_user_id_fk
            references "user"
            on update cascade on delete cascade,
    order_id   integer               not null,
    value      double precision      not null,
    created_at timestamp             not null,
    is_debited boolean default false not null,
    service_id integer               not null
        constraint reserve_service_id_fk
            references service
);

create index reserve_user_id_index
    on reserve (user_id);

create unique index reserve_order_id_uindex
    on reserve (order_id);

create index reserve_service_id_index
    on reserve (service_id);
