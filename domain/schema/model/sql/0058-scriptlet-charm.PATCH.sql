CREATE TABLE scriptlet_charm (
    charm_uuid TEXT NOT NULL PRIMARY KEY,
    scriptlet TEXT NOT NULL,
    CONSTRAINT fk_scriptlet_charm_charm
    FOREIGN KEY (charm_uuid) REFERENCES charm (uuid)
);

DROP VIEW IF EXISTS v_application_origin;

CREATE VIEW v_application_origin AS
SELECT
    a.uuid,
    c.reference_name,
    c.source_id,
    c.revision,
    cdi.charmhub_identifier,
    ch.hash
FROM application AS a
JOIN charm AS c ON a.charm_uuid = c.uuid
LEFT JOIN charm_download_info AS cdi ON c.uuid = cdi.charm_uuid
LEFT JOIN charm_hash AS ch ON c.uuid = ch.charm_uuid;
