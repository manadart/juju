-- Retain UUIDs for execution and add model-scoped user-facing request IDs.
ALTER TABLE storage_resize_request
ADD COLUMN request_id INTEGER NOT NULL DEFAULT 0 CHECK (request_id >= 0);

WITH numbered AS (
    SELECT uuid, ROW_NUMBER() OVER (ORDER BY created_at, uuid) AS request_id
    FROM storage_resize_request
)
UPDATE storage_resize_request AS r
SET request_id = numbered.request_id
FROM numbered
WHERE r.uuid = numbered.uuid;

CREATE UNIQUE INDEX idx_storage_resize_request_id
ON storage_resize_request (request_id);

-- NextValue increments this value, including in a fresh model with no requests.
INSERT INTO sequence (namespace, value)
SELECT 'storage_resize', COALESCE(MAX(request_id), 0)
FROM storage_resize_request;
