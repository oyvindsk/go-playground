CREATE TABLE Posts (
    id          integer NOT NULL     PRIMARY KEY,
    title       text    NOT NULL                ,
    body        text    NOT NULL                , 
    created_at  integer NOT NULL                ,
    updated_at  integer NOT NULL
);