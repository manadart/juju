-- Bind execution to a machine and resource identities across agent restarts.
-- These are historical identities, so deletion must not cascade to requests.
ALTER TABLE storage_resize_request ADD COLUMN executor_machine_uuid TEXT;
ALTER TABLE storage_resize_request ADD COLUMN execution_identity TEXT;

CREATE INDEX idx_storage_resize_request_executor
ON storage_resize_request (executor_machine_uuid) WHERE status = 'pending';
