-- Pending intent is separate from both allocation and observed capacity.
CREATE TABLE storage_resize_intent (
    uuid TEXT NOT NULL PRIMARY KEY,
    storage_instance_uuid TEXT NOT NULL,
    requested_size_mib INT NOT NULL CHECK (requested_size_mib > 0),
    created_at DATETIME NOT NULL,
    CONSTRAINT fk_storage_resize_intent_storage_instance
    FOREIGN KEY (storage_instance_uuid)
    REFERENCES storage_instance (uuid) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_storage_resize_intent_storage_instance
ON storage_resize_intent (storage_instance_uuid);
