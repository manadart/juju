CREATE TABLE storage_resize_plan (
    request_uuid TEXT NOT NULL PRIMARY KEY,
    identity TEXT NOT NULL,
    backend_machine_uuid TEXT NOT NULL,
    node_machine_uuid TEXT NOT NULL,
    approved_size_mib INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_storage_resize_plan_request FOREIGN KEY (request_uuid)
        REFERENCES storage_resize_request (uuid)
);

CREATE TABLE storage_resize_stage (
    request_uuid TEXT NOT NULL,
    stage TEXT NOT NULL CHECK (stage IN ('backend', 'attachment')),
    state TEXT NOT NULL CHECK (state IN
        ('checking', 'checked', 'running', 'completed', 'failed', 'skipped')),
    message TEXT NOT NULL DEFAULT '',
    prerequisites TEXT NOT NULL DEFAULT '[]',
    follow_up TEXT NOT NULL DEFAULT '[]',
    observed_size_mib INT NOT NULL DEFAULT 0,
    filesystem_size_mib INT NOT NULL DEFAULT 0,
    PRIMARY KEY (request_uuid, stage),
    CONSTRAINT fk_storage_resize_stage_plan FOREIGN KEY (request_uuid)
        REFERENCES storage_resize_plan (request_uuid)
);

CREATE INDEX idx_storage_resize_plan_backend_machine
ON storage_resize_plan (backend_machine_uuid);
CREATE INDEX idx_storage_resize_plan_node_machine
ON storage_resize_plan (node_machine_uuid);
