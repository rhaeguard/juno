DROP TABLE IF EXISTS tag_relations;
DROP TABLE IF EXISTS resource_relations;
DROP TABLE IF EXISTS applications;
DROP TABLE IF EXISTS resources;

create table resources(
    id          character varying not null,
    name        character varying not null,
    extension   character varying not null,
    size        integer not null,
    created_on  timestamp,
    PRIMARY KEY(id)
);
create table applications(
    id              character varying not null,
    description     character varying,
    password        character varying not null,
    PRIMARY KEY (id)
);
create table resource_relations(
    id                  serial,
    app_id              character varying not null,
    resource_id         character varying not null,
    saved_location      character varying not null,
    PRIMARY KEY (id),
    FOREIGN KEY (app_id) REFERENCES applications (id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources (id) ON DELETE CASCADE
);
create table tag_relations(
    id                  serial,
    resource_id         character varying not null,
    tag                 character varying not null,
    PRIMARY KEY (id),
    FOREIGN KEY (resource_id) REFERENCES resources (id) ON DELETE CASCADE
);