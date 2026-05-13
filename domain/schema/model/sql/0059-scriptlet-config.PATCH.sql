-- Scriptlet configuration table (runtime, app)
CREATE TABLE scriptlet_config (
    charm_uuid TEXT NOT NULL PRIMARY KEY,
    runtime TEXT NOT NULL,
    app TEXT NOT NULL,
    CONSTRAINT fk_scriptlet_config_charm
    FOREIGN KEY (charm_uuid) REFERENCES charm (uuid)
);

CREATE INDEX idx_scriptlet_config_runtime
ON scriptlet_config (runtime);

-- Events list (unordered set)
CREATE TABLE scriptlet_event (
    charm_uuid TEXT NOT NULL,
    event_name TEXT NOT NULL,
    CONSTRAINT fk_scriptlet_event_charm
    FOREIGN KEY (charm_uuid) REFERENCES charm (uuid),
    PRIMARY KEY (charm_uuid, event_name)
);

CREATE INDEX idx_scriptlet_event_charm
ON scriptlet_event (charm_uuid);
