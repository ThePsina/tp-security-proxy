drop table IF EXISTS requests CASCADE;
drop table IF EXISTS headers;

create table requests
(
    id      bigserial primary key,
    host    text,
    request text
);

create table headers
(
    req_id bigint references requests (id) on delete cascade,
    key    text,
    val    text
);
