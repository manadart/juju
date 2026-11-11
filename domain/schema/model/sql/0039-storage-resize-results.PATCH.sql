-- Requests retain historical identity even after storage is removed.
CREATE TABLE storage_resize_request (
    uuid TEXT NOT NULL PRIMARY KEY,
    storage_instance_uuid TEXT NOT NULL,
    storage_id TEXT NOT NULL,
    requested_size_mib INT NOT NULL CHECK (requested_size_mib > 0),
    created_at DATETIME NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'completed', 'failed')),
    result_message TEXT NOT NULL DEFAULT '',
    completed_at DATETIME,
    CONSTRAINT chk_storage_resize_request_completion
    CHECK (
        (status = 'pending' AND completed_at IS NULL)
        OR (status IN ('completed', 'failed') AND completed_at IS NOT NULL)
    )
);

INSERT INTO storage_resize_request (
    uuid, storage_instance_uuid, storage_id, requested_size_mib, created_at
)
SELECT sri.uuid, sri.storage_instance_uuid, si.storage_id,
       sri.requested_size_mib, sri.created_at
FROM storage_resize_intent AS sri
JOIN storage_instance AS si ON si.uuid = sri.storage_instance_uuid;

DROP TABLE storage_resize_intent;

CREATE UNIQUE INDEX idx_storage_resize_request_pending_instance
ON storage_resize_request (storage_instance_uuid) WHERE status = 'pending';

-- Removal must leave a terminal result for clients waiting on the request.
CREATE TRIGGER trg_storage_resize_request_storage_removed
AFTER DELETE ON storage_instance
FOR EACH ROW
BEGIN
    UPDATE storage_resize_request
    SET status = 'failed',
        result_message = 'Storage removed before resize completed',
        completed_at = CURRENT_TIMESTAMP
    WHERE storage_instance_uuid = OLD.uuid AND status = 'pending';
END;
