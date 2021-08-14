DROP TABLE vm;

CREATE TABLE vm
(
    vm_id          INTEGER PRIMARY KEY,
    name           TEXT    NOT NULL,
    parent         TEXT    NOT NULL,
    overall_status TEXT    NOT NULL
);