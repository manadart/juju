def init():
    juju.observe("config_changed", on_config_changed)
    juju.observe("relation_created", on_relation_created)
    juju.observe("relation_joined", on_relation_joined)
    juju.observe("relation_changed", on_relation_changed)
    juju.observe("relation_departed", on_relation_departed)
    juju.observe("relation_broken", on_relation_broken)

def on_config_changed(event):
    _stub_intent(event, "config changed")

def on_relation_created(event):
    _stub_intent(event, "relation created")

def on_relation_joined(event):
    _stub_intent(event, "relation joined")

def on_relation_changed(event):
    _stub_intent(event, "relation changed")

def on_relation_departed(event):
    _stub_intent(event, "relation departed")

def on_relation_broken(event):
    _stub_intent(event, "relation broken")

def _stub_intent(event, message):
    juju.status_set("active", message = message)
    juju.state_set("last_event", event.name)
    juju.state_set("last_message", message)
