// Code generated by triggergen. DO NOT EDIT.

package triggers

import (
	"fmt"

	"github.com/juju/juju/core/database/schema"
)


// ChangeLogTriggersForUserAuthentication generates the triggers for the
// user_authentication table.
func ChangeLogTriggersForUserAuthentication(columnName string, namespaceID int) func() schema.Patch {
	return func() schema.Patch {
		return schema.MakePatch(fmt.Sprintf(`
-- insert namespace for UserAuthentication
INSERT INTO change_log_namespace VALUES (%[2]d, 'user_authentication', 'UserAuthentication changes based on %[1]s');

-- insert trigger for UserAuthentication
CREATE TRIGGER trg_log_user_authentication_insert
AFTER INSERT ON user_authentication FOR EACH ROW
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (1, %[2]d, NEW.%[1]s, DATETIME('now'));
END;

-- update trigger for UserAuthentication
CREATE TRIGGER trg_log_user_authentication_update
AFTER UPDATE ON user_authentication FOR EACH ROW
WHEN 
	NEW.user_uuid != OLD.user_uuid OR
	NEW.disabled != OLD.disabled 
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (2, %[2]d, OLD.%[1]s, DATETIME('now'));
END;
-- delete trigger for UserAuthentication
CREATE TRIGGER trg_log_user_authentication_delete
AFTER DELETE ON user_authentication FOR EACH ROW
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (4, %[2]d, OLD.%[1]s, DATETIME('now'));
END;`, columnName, namespaceID))
	}
}
