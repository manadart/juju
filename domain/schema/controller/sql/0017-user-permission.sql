-- OpenFGA v1.22.1 logical storage, owned by Juju's controller schema.
-- This beta prototype requires a freshly bootstrapped controller.
CREATE TABLE openfga_store (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    deleted_at TEXT
);

CREATE TABLE openfga_authorization_model (
    store_id TEXT NOT NULL,
    model_id TEXT NOT NULL,
    model TEXT NOT NULL,
    PRIMARY KEY (store_id, model_id)
);

CREATE TABLE openfga_tuple (
    store_id TEXT NOT NULL,
    object_type TEXT NOT NULL,
    object_id TEXT NOT NULL,
    relation TEXT NOT NULL,
    user TEXT NOT NULL,
    condition TEXT NOT NULL DEFAULT '',
    permission_uuid TEXT,
    ulid TEXT NOT NULL DEFAULT (
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 45) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 40) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 35) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 30) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 25) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 20) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 15) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 10) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 5) & 31) + 1, 1) ||
        substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER) >> 0) & 31) + 1, 1) ||
        hex(randomblob(8))
    ),
    inserted_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f000000Z', 'now')),
    PRIMARY KEY (store_id, object_type, object_id, relation, user)
);
CREATE UNIQUE INDEX idx_openfga_tuple_ulid ON openfga_tuple (store_id, ulid);
CREATE UNIQUE INDEX idx_openfga_permission_uuid
ON openfga_tuple (permission_uuid) WHERE permission_uuid IS NOT NULL;
CREATE INDEX idx_openfga_tuple_user
ON openfga_tuple (store_id, user, object_type, relation, object_id);
-- Juju represents each explicit grant by its highest access level.
CREATE UNIQUE INDEX idx_openfga_juju_grant
ON openfga_tuple (object_type, object_id, user)
WHERE store_id = '01J00000000000000000000000'
AND relation IN ('read', 'write', 'consume', 'admin', 'login', 'add_model', 'superuser');

-- A transactional counter orders changes even within one millisecond.
CREATE TABLE openfga_sequence (
    id INT NOT NULL PRIMARY KEY CHECK (id = 0),
    millis INT NOT NULL,
    counter INT NOT NULL
);
INSERT INTO openfga_sequence VALUES (0, 0, 0);

CREATE TABLE openfga_changelog (
    store_id TEXT NOT NULL,
    ulid TEXT NOT NULL,
    object_type TEXT NOT NULL,
    object_id TEXT NOT NULL,
    relation TEXT NOT NULL,
    user TEXT NOT NULL,
    condition TEXT NOT NULL DEFAULT '',
    operation INT NOT NULL,
    inserted_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%f000000Z', 'now')),
    PRIMARY KEY (store_id, ulid)
);
CREATE INDEX idx_openfga_changelog_type
ON openfga_changelog (store_id, object_type, ulid);

-- Markers distinguish replay after an ambiguous commit from a new request.
CREATE TABLE openfga_operation (ulid TEXT NOT NULL PRIMARY KEY);

CREATE TABLE openfga_assertion (
    store_id TEXT NOT NULL,
    model_id TEXT NOT NULL,
    assertions TEXT NOT NULL,
    PRIMARY KEY (store_id, model_id)
);

INSERT INTO openfga_store (id, name, created_at, updated_at)
VALUES ('01J00000000000000000000000', 'juju', '2026-11-11T00:00:00.000000000Z', '2026-11-11T00:00:00.000000000Z');
INSERT INTO openfga_authorization_model (store_id, model_id, model)
VALUES ('01J00000000000000000000000', '01J00000000000000000000001', '{"id":"01J00000000000000000000001","schema_version":"1.1","type_definitions":[{"type":"user"},{"type":"external"},{"type":"controller","relations":{"login":{"this":{}},"can_login":{"union":{"child":[{"computedUserset":{"relation":"login"}},{"computedUserset":{"relation":"superuser"}}]}},"superuser":{"this":{}},"can_superuser":{"computedUserset":{"relation":"superuser"}}},"metadata":{"relations":{"login":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"superuser":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]}}}},{"type":"cloud","relations":{"add_model":{"this":{}},"can_add_model":{"union":{"child":[{"computedUserset":{"relation":"add_model"}},{"computedUserset":{"relation":"admin"}}]}},"admin":{"this":{}},"can_admin":{"computedUserset":{"relation":"admin"}}},"metadata":{"relations":{"add_model":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"admin":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]}}}},{"type":"model","relations":{"read":{"this":{}},"can_read":{"union":{"child":[{"computedUserset":{"relation":"read"}},{"computedUserset":{"relation":"write"}},{"computedUserset":{"relation":"admin"}}]}},"write":{"this":{}},"can_write":{"union":{"child":[{"computedUserset":{"relation":"write"}},{"computedUserset":{"relation":"admin"}}]}},"admin":{"this":{}},"can_admin":{"computedUserset":{"relation":"admin"}}},"metadata":{"relations":{"read":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"write":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"admin":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]}}}},{"type":"offer","relations":{"read":{"this":{}},"can_read":{"union":{"child":[{"computedUserset":{"relation":"read"}},{"computedUserset":{"relation":"consume"}},{"computedUserset":{"relation":"admin"}}]}},"consume":{"this":{}},"can_consume":{"union":{"child":[{"computedUserset":{"relation":"consume"}},{"computedUserset":{"relation":"admin"}}]}},"admin":{"this":{}},"can_admin":{"computedUserset":{"relation":"admin"}}},"metadata":{"relations":{"read":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"consume":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]},"admin":{"directly_related_user_types":[{"type":"user"},{"type":"external"},{"type":"external","wildcard":{}}]}}}}]}');

-- Every tuple writer, including domain transactions, records its changes.
-- The timestamp prefix and sequence suffix form a Crockford-base32 ULID.
-- Updates are represented as a deletion followed by an insertion.
CREATE TRIGGER trg_openfga_tuple_insert
AFTER INSERT ON openfga_tuple
BEGIN
    UPDATE openfga_sequence
    SET counter = counter + 1,
        millis = MAX(millis, CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER))
    WHERE id = 0;
    INSERT INTO openfga_changelog
        (ulid, store_id, object_type, object_id, relation, user, condition, operation)
    SELECT substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 45) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 40) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 35) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 30) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 25) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 20) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 15) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 10) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 5) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 0) & 31) + 1, 1) ||
           printf('%016X', counter),
           NEW.store_id, NEW.object_type, NEW.object_id,
           NEW.relation, NEW.user, NEW.condition, 0
    FROM openfga_sequence WHERE id = 0;
END;

CREATE TRIGGER trg_openfga_tuple_delete
AFTER DELETE ON openfga_tuple
BEGIN
    UPDATE openfga_sequence
    SET counter = counter + 1,
        millis = MAX(millis, CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER))
    WHERE id = 0;
    INSERT INTO openfga_changelog
        (ulid, store_id, object_type, object_id, relation, user, condition, operation)
    SELECT substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 45) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 40) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 35) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 30) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 25) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 20) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 15) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 10) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 5) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 0) & 31) + 1, 1) ||
           printf('%016X', counter),
           OLD.store_id, OLD.object_type, OLD.object_id,
           OLD.relation, OLD.user, '', 1
    FROM openfga_sequence WHERE id = 0;
END;

CREATE TRIGGER trg_openfga_tuple_update
AFTER UPDATE ON openfga_tuple
BEGIN
    UPDATE openfga_sequence
    SET counter = counter + 1,
        millis = MAX(millis, CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER))
    WHERE id = 0;
    INSERT INTO openfga_changelog
        (ulid, store_id, object_type, object_id, relation, user, condition, operation)
    SELECT substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 45) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 40) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 35) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 30) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 25) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 20) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 15) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 10) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 5) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 0) & 31) + 1, 1) ||
           printf('%016X', counter),
           OLD.store_id, OLD.object_type, OLD.object_id,
           OLD.relation, OLD.user, '', 1
    FROM openfga_sequence WHERE id = 0;
    UPDATE openfga_sequence
    SET counter = counter + 1,
        millis = MAX(millis, CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER))
    WHERE id = 0;
    INSERT INTO openfga_changelog
        (ulid, store_id, object_type, object_id, relation, user, condition, operation)
    SELECT substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 45) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 40) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 35) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 30) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 25) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 20) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 15) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 10) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 5) & 31) + 1, 1) ||
           substr('0123456789ABCDEFGHJKMNPQRSTVWXYZ', ((millis >> 0) & 31) + 1, 1) ||
           printf('%016X', counter),
           NEW.store_id, NEW.object_type, NEW.object_id,
           NEW.relation, NEW.user, NEW.condition, 0
    FROM openfga_sequence WHERE id = 0;
END;

-- Read-only projections preserve existing list/export-to-model APIs.
-- All mutations target openfga_tuple; these are not an authorisation engine.
CREATE VIEW permission_access_type AS
SELECT 0 AS id, 'read' AS type UNION ALL SELECT 1, 'write'
UNION ALL SELECT 2, 'consume' UNION ALL SELECT 3, 'admin'
UNION ALL SELECT 4, 'login' UNION ALL SELECT 5, 'add-model'
UNION ALL SELECT 6, 'superuser';
CREATE VIEW permission_object_type AS
SELECT 0 AS id, 'cloud' AS type UNION ALL SELECT 1, 'controller'
UNION ALL SELECT 2, 'model' UNION ALL SELECT 3, 'offer';
CREATE VIEW permission_object_access AS
SELECT 0 AS id, 3 AS access_type_id, 0 AS object_type_id
UNION ALL SELECT 1, 5, 0 UNION ALL SELECT 2, 4, 1
UNION ALL SELECT 3, 6, 1 UNION ALL SELECT 4, 0, 2
UNION ALL SELECT 5, 1, 2 UNION ALL SELECT 6, 3, 2
UNION ALL SELECT 7, 0, 3 UNION ALL SELECT 8, 2, 3
UNION ALL SELECT 9, 3, 3;
CREATE VIEW permission AS
SELECT COALESCE(t.permission_uuid, t.ulid) AS uuid,
       at.id AS access_type_id, ot.id AS object_type_id,
       t.object_id AS grant_on, u.uuid AS grant_to
FROM openfga_tuple AS t
JOIN permission_access_type AS at ON replace(at.type, '-', '_') = t.relation
JOIN permission_object_type AS ot ON ot.type = t.object_type
JOIN user AS u ON t.user = CASE
    WHEN u.name = 'everyone@external' THEN 'external:*'
    WHEN u.external THEN 'external:' || u.uuid
    ELSE 'user:' || u.uuid END
WHERE t.store_id = '01J00000000000000000000000';

-- All permissions
CREATE VIEW v_permission AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    at.type AS access_type,
    ot.type AS object_type
FROM permission AS p
JOIN permission_access_type AS at ON p.access_type_id = at.id
JOIN permission_object_type AS ot ON p.object_type_id = ot.id;

-- All model permissions, verifying the model does exist.
CREATE VIEW v_permission_model AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    p.access_type,
    p.object_type
FROM v_permission AS p
JOIN model ON p.grant_on = model.uuid
WHERE p.object_type = 'model';

-- All controller cloud, verifying the cloud does exist.
CREATE VIEW v_permission_cloud AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    p.access_type,
    p.object_type
FROM v_permission AS p
JOIN cloud ON p.grant_on = cloud.name
WHERE p.object_type = 'cloud';

-- All controller permissions, verifying the controller does exists.
CREATE VIEW v_permission_controller AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    p.access_type,
    p.object_type
FROM v_permission AS p
JOIN controller ON p.grant_on = controller.uuid
WHERE p.object_type = 'controller';

-- All offer permissions, NOT verifying the offer does exist.
CREATE VIEW v_permission_offer AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    p.access_type,
    p.object_type
FROM v_permission AS p
WHERE p.object_type = 'offer';

-- The permissions for the special user everyone@external.
CREATE VIEW v_everyone_external AS
SELECT
    p.uuid,
    p.grant_on,
    p.grant_to,
    p.access_type,
    p.object_type
FROM v_permission AS p
JOIN user AS u ON p.grant_to = u.uuid
WHERE u.name = 'everyone@external';
