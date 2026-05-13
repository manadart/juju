ALTER TABLE charm ADD COLUMN is_scriptlet BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE scriptlet_charm (
    charm_uuid TEXT NOT NULL PRIMARY KEY,
    scriptlet TEXT NOT NULL,
    CONSTRAINT fk_scriptlet_charm_charm
    FOREIGN KEY (charm_uuid) REFERENCES charm (uuid)
);
