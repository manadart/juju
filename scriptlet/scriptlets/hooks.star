def init():
    juju.observe("config_changed", on_config_changed)
    juju.observe("relation_created", on_relation_created)
    juju.observe("relation_joined", on_relation_joined)
    juju.observe("relation_changed", on_relation_changed)
    juju.observe("relation_departed", on_relation_departed)
    juju.observe("relation_broken", on_relation_broken)

def on_config_changed(event):
    juju.status_set("active", message = "config changed")

def on_relation_created(event):
    juju.status_set("active", message = "relation created")

def on_relation_joined(event):
    juju.status_set("active", message = "related")
    juju.state_set("controller-uuid", event.controller_uuid)
    juju.state_set("model-name", event.model_name)

def on_relation_changed(event):
    juju.status_set("active", message = "relation changed")

def on_relation_departed(event):
    juju.status_set("waiting", message = "relation departed")

def on_relation_broken(event):
    juju.status_set("waiting", message = "relation broken")
