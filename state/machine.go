// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/mgo/v3"
	"github.com/juju/mgo/v3/bson"
	"github.com/juju/mgo/v3/txn"
	"github.com/juju/names/v6"
	jujutxn "github.com/juju/txn/v3"
	"github.com/kr/pretty"

	"github.com/juju/juju/core/actions"
	"github.com/juju/juju/core/constraints"
	corecontainer "github.com/juju/juju/core/container"
	"github.com/juju/juju/core/instance"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/objectstore"
	"github.com/juju/juju/core/semversion"
	"github.com/juju/juju/core/status"
	internalpassword "github.com/juju/juju/internal/password"
	"github.com/juju/juju/internal/tools"
	stateerrors "github.com/juju/juju/state/errors"
)

// Machine represents the state of a machine.
type Machine struct {
	st  *State
	doc machineDoc
}

// machineDoc represents the internal state of a machine in MongoDB.
// Note the correspondence with MachineInfo in apiserver/juju.
type machineDoc struct {
	DocID          string `bson:"_id"`
	Id             string `bson:"machineid"`
	ModelUUID      string `bson:"model-uuid"`
	Base           Base   `bson:"base"`
	ContainerType  string
	Principals     []string
	Life           Life
	Tools          *tools.Tools `bson:",omitempty"`
	PasswordHash   string
	Clean          bool
	ForceDestroyed bool `bson:"force-destroyed"`

	// Volumes contains the names of volumes attached to the machine.
	Volumes []string `bson:"volumes,omitempty"`
	// Filesystems contains the names of filesystems attached to the machine.
	Filesystems []string `bson:"filesystems,omitempty"`

	// We store 2 different sets of addresses for the machine, obtained
	// from different sources.
	// Addresses is the set of addresses obtained by asking the provider.
	Addresses []address

	// MachineAddresses is the set of addresses obtained from the machine itself.
	MachineAddresses []address

	// PreferredPublicAddress is the preferred address to be used for
	// the machine when a public address is requested.
	PreferredPublicAddress address `bson:",omitempty"`

	// PreferredPrivateAddress is the preferred address to be used for
	// the machine when a private address is requested.
	PreferredPrivateAddress address `bson:",omitempty"`

	// The SupportedContainers attributes are used to advertise what containers this
	// machine is capable of hosting.
	SupportedContainersKnown bool
	SupportedContainers      []instance.ContainerType `bson:",omitempty"`
	// Placement is the placement directive that should be used when provisioning
	// an instance for the machine.
	Placement string `bson:",omitempty"`

	// AgentStartedAt records the time when the machine agent started.
	AgentStartedAt time.Time `bson:"agent-started-at,omitempty"`

	// Hostname records the machine's hostname as reported by the machine agent.
	Hostname string `bson:"hostname,omitempty"`
}

func newMachine(st *State, doc *machineDoc) *Machine {
	machine := &Machine{
		st:  st,
		doc: *doc,
	}
	return machine
}

// DocID returns the machine doc id.
// Deprecated: this is only used for migration to the domain model.
func (m *Machine) DocID() string {
	return m.doc.DocID
}

// Id returns the machine id.
func (m *Machine) Id() string {
	return m.doc.Id
}

// Principals returns the principals for the machine.
func (m *Machine) Principals() []string {
	return m.doc.Principals
}

// AddPrincipal adds a principal to the machine.
func (m *Machine) AddPrincipal(name string) {
	m.doc.Clean = false
	m.doc.Principals = append(m.doc.Principals, name)
	sort.Strings(m.doc.Principals)
}

// Base returns the os base running on the machine.
func (m *Machine) Base() Base {
	return m.doc.Base
}

// ContainerType returns the type of container hosting this machine.
func (m *Machine) ContainerType() instance.ContainerType {
	return instance.ContainerType(m.doc.ContainerType)
}

// ModelUUID returns the unique identifier
// for the model that this machine is in.
func (m *Machine) ModelUUID() string {
	return m.doc.ModelUUID
}

// FileSystems returns the names of the filesystems attached to the machine.
func (m *Machine) FileSystems() []string {
	return m.doc.Filesystems
}

// ForceDestroyed returns whether the destruction of a dying/dead
// machine was forced. It's always false for a machine that's alive.
func (m *Machine) ForceDestroyed() bool {
	return m.doc.ForceDestroyed
}

func (m *Machine) forceDestroyedOps() []txn.Op {
	return []txn.Op{{
		C:      machinesC,
		Id:     m.doc.DocID,
		Assert: txn.DocExists,
		Update: bson.D{{"$set", bson.D{{"force-destroyed", true}}}},
	}}
}

// machineGlobalKey returns the global database key for the identified machine.
func machineGlobalKey(id string) string {
	return machineGlobalKeyPrefix + id
}

// machineGlobalKeyPrefix is the kind string we use to denote machine kind.
const machineGlobalKeyPrefix = "m#"

// machineGlobalInstanceKey returns the global database key for the identified
// machine's instance.
func machineGlobalInstanceKey(id string) string {
	return machineGlobalKey(id) + "#instance"
}

// InstanceKind returns a human readable name identifying the machine instance
// kind.
func (m *Machine) InstanceKind() string {
	return m.Tag().Kind() + "-instance"
}

// globalInstanceKey returns the global database key for the machine's instance.
func (m *Machine) globalInstanceKey() string {
	return machineGlobalInstanceKey(m.doc.Id)
}

// machineGlobalModificationKey returns the global database key for the
// identified machine's modification changes.
func machineGlobalModificationKey(id string) string {
	return machineGlobalKey(id) + "#modification"
}

// ModificationKind returns the human readable kind string we use for when an
// lxd profile is applied on a machine..
func (m *Machine) ModificationKind() string {
	return m.Tag().Kind() + "-lxd-profile"
}

// globalModificationKey returns the global database key for the machine's
// modification changes.
func (m *Machine) globalModificationKey() string {
	return machineGlobalModificationKey(m.doc.Id)
}

// globalKey returns the global database key for the machine.
func (m *Machine) globalKey() string {
	return machineGlobalKey(m.doc.Id)
}

// Tag returns a tag identifying the machine. The String method provides a
// string representation that is safe to use as a file name. The returned name
// will be different from other Tag values returned by any other entities
// from the same state.
func (m *Machine) Tag() names.Tag {
	return m.MachineTag()
}

// Kind returns a human readable name identifying the machine kind.
func (m *Machine) Kind() string {
	return m.Tag().Kind()
}

// MachineTag returns the more specific MachineTag type as opposed
// to the more generic Tag type.
func (m *Machine) MachineTag() names.MachineTag {
	return names.NewMachineTag(m.Id())
}

// Life returns whether the machine is Alive, Dying or Dead.
func (m *Machine) Life() Life {
	return m.doc.Life
}

// AgentTools returns the tools that the agent is currently running.
// It returns an error that satisfies errors.IsNotFound if the tools
// have not yet been set.
func (m *Machine) AgentTools() (*tools.Tools, error) {
	if m.doc.Tools == nil {
		return nil, errors.NotFoundf("agent binaries for machine %v", m)
	}
	tools := *m.doc.Tools
	return &tools, nil
}

// checkVersionValidity checks whether the given version is suitable
// for passing to SetAgentVersion.
func checkVersionValidity(v semversion.Binary) error {
	if v.Release == "" || v.Arch == "" {
		return fmt.Errorf("empty series or arch")
	}
	return nil
}

// SetAgentVersion sets the version of juju that the agent is
// currently running.
func (m *Machine) SetAgentVersion(v semversion.Binary) (err error) {
	defer errors.DeferredAnnotatef(&err, "setting agent version for machine %v", m)
	ops, tools, err := m.setAgentVersionOps(v)
	if err != nil {
		return errors.Trace(err)
	}
	// A "raw" transaction is needed here because this function gets
	// called before database migrations have run so we don't
	// necessarily want the model UUID added to the id.
	if err := m.st.runRawTransaction(ops); err != nil {
		return onAbort(err, stateerrors.ErrDead)
	}
	m.doc.Tools = tools
	return nil
}

func (m *Machine) setAgentVersionOps(v semversion.Binary) ([]txn.Op, *tools.Tools, error) {
	if err := checkVersionValidity(v); err != nil {
		return nil, nil, err
	}
	tools := &tools.Tools{Version: v}
	ops := []txn.Op{{
		C:      machinesC,
		Id:     m.doc.DocID,
		Assert: notDeadDoc,
		Update: bson.D{{"$set", bson.D{{"tools", tools}}}},
	}}
	return ops, tools, nil
}

// PasswordValid returns whether the given password is valid
// for the given machine.
func (m *Machine) PasswordValid(password string) bool {
	agentHash := internalpassword.AgentPasswordHash(password)
	return agentHash == m.doc.PasswordHash
}

// Destroy sets the machine lifecycle to Dying if it is Alive. It does
// nothing otherwise. Destroy will fail if the machine has principal
// units assigned, or if the machine has JobManageModel.
// If the machine has assigned units, Destroy will return
// a HasAssignedUnitsError.  If the machine has containers, Destroy
// will return HasContainersError.
func (m *Machine) Destroy(_ objectstore.ObjectStore) error {
	return errors.Trace(m.advanceLifecycle(Dying, false, false, 0))
}

// DestroyWithContainers sets the machine lifecycle to Dying if it is Alive.
// It does nothing otherwise. DestroyWithContainers will fail if the machine
// has principal units assigned, or if the machine has JobManageModel. If the
// machine has assigned units, DestroyWithContainers will return a
// HasAssignedUnitsError.  The machine is allowed to have containers.  Use with
// caution.  Intended for model tear down.
func (m *Machine) DestroyWithContainers() error {
	return m.advanceLifecycle(Dying, false, true, 0)
}

// ForceDestroy queues the machine for complete removal, including the
// destruction of all units and containers on the machine.
func (m *Machine) ForceDestroy(maxWait time.Duration) error {
	ops, err := m.forceDestroyOps(maxWait)
	if err != nil {
		return errors.Trace(err)
	}
	if err := m.st.db().RunTransaction(ops); err != txn.ErrAborted {
		return errors.Annotatef(err, "failed to run transaction: %s", pretty.Sprint(ops))
	}
	return nil
}

func (m *Machine) forceDestroyOps(maxWait time.Duration) ([]txn.Op, error) {
	// Make sure the machine doesn't become a manager while we're destroying it
	return []txn.Op{newCleanupOp(cleanupForceDestroyedMachine, m.doc.Id, maxWait)}, nil
}

// EnsureDead sets the machine lifecycle to Dead if it is Alive or Dying.
// It does nothing otherwise. EnsureDead will fail if the machine has
// principal units assigned, or if the machine has JobManageModel.
// If the machine has assigned units, EnsureDead will return
// a HasAssignedUnitsError.
func (m *Machine) EnsureDead() error {
	return m.advanceLifecycle(Dead, false, false, 0)
}

// Containers returns the container ids belonging to a parent machine.
// TODO(wallyworld): move this method to a service
func (m *Machine) Containers() ([]string, error) {
	containerRefs, closer := m.st.db().GetCollection(containerRefsC)
	defer closer()

	var mc machineContainers
	err := containerRefs.FindId(m.doc.DocID).One(&mc)
	if err == nil {
		return mc.Children, nil
	}
	if err == mgo.ErrNotFound {
		return nil, errors.NotFoundf("container info for machine %v", m.Id())
	}
	return nil, err
}

// ParentId returns the Id of the host machine if this machine is a container.
func (m *Machine) ParentId() (string, bool) {
	parentId := corecontainer.ParentId(m.Id())
	return parentId, parentId != ""
}

// IsContainer returns true if the machine is a container.
func (m *Machine) IsContainer() bool {
	_, isContainer := m.ParentId()
	return isContainer
}

// advanceLifecycle ensures that the machine's lifecycle is no earlier
// than the supplied value. If the machine already has that lifecycle
// value, or a later one, no changes will be made to remote state. If
// the machine has any responsibilities that preclude a valid change in
// lifecycle, it will return an error. dyingAllowContainers indicates
// whether the machine can have containers when moving to the dying state.
// Not allowed for moving to dead.
func (original *Machine) advanceLifecycle(life Life, force, dyingAllowContainers bool, maxWait time.Duration) (err error) {
	logger.Debugf(context.TODO(), "%s.advanceLifecycle(%s, %t, %t)", original.Id(), life, force, dyingAllowContainers)

	if life == Dead && dyingAllowContainers {
		return errors.BadRequestf("life cannot be Dead if dyingAllowContainers true.")
	}

	// A machine can be set to dying with containers, but cannot have any when
	// advanced to dead.
	if !dyingAllowContainers && life == Dying {
		if err := original.advanceLifecycleIfNoContainers(); err != nil {
			return err
		}
	}

	m := original
	defer func() {
		if err == nil {
			// The machine's lifecycle is known to have advanced; it may be
			// known to have already advanced further than requested, in
			// which case we set the latest known valid value.
			if m == nil {
				life = Dead
			} else if m.doc.Life > life {
				life = m.doc.Life
			}
			original.doc.Life = life
		}
	}()

	ops := m.advanceLifecyleInitialOps(life)

	// multiple attempts: one with original data, one with refreshed data, and a final
	// one intended to determine the cause of failure of the preceding attempt.
	buildTxn := func(attempt int) ([]txn.Op, error) {
		var asserts bson.D
		// Grab a fresh copy of the machine data.
		// We don't write to original, because the expectation is that state-
		// changing methods only set the requested change on the receiver; a case
		// could perhaps be made that this is not a helpful convention in the
		// context of the new state API, but we maintain consistency in the
		// face of uncertainty.
		if m, err = m.st.Machine(m.doc.Id); errors.Is(err, errors.NotFound) {
			return nil, jujutxn.ErrNoOperations
		} else if err != nil {
			return nil, err
		}

		// Check that the life change is sane, and collect the assertions
		// necessary to determine that it remains so.
		switch life {
		case Dying:
			if m.doc.Life != Alive {
				return nil, jujutxn.ErrNoOperations
			}
			// Manager nodes are allowed to go to dying even when they have
			// the vote, as that is used as the signal that they should lose
			// their vote.
			asserts = append(asserts, isAliveDoc...)
		case Dead:
			if m.doc.Life == Dead {
				return nil, jujutxn.ErrNoOperations
			}
			asserts = append(asserts, notDeadDoc...)
			ops = append(ops, controllerAdvanceLifecyleVoteOp())
		default:
			panic(fmt.Errorf("cannot advance lifecycle to %v", life))
		}

		// Check that the machine does not have any responsibilities that
		// prevent a lifecycle change.
		// If there are no alive units left on the machine, or all the applications are dying,
		// then the machine may be soon destroyed by a cleanup worker.
		// In that case, we don't want to return any error about not being able to
		// destroy a machine with units as it will be a lie.
		if life == Dying {
			canDie := true

			var principalUnitNames []string
			for _, principalUnit := range m.doc.Principals {
				principalUnitNames = append(principalUnitNames, principalUnit)
				canDie, err = m.assessCanDieUnit(principalUnit)
				if err != nil {
					return nil, errors.Trace(err)
				}
				if !canDie {
					break
				}
			}

			if canDie && !dyingAllowContainers {
				err := m.advanceLifecycleIfNoContainers()
				if errors.Is(err, stateerrors.HasContainersError) {
					canDie = false
				} else if err != nil {
					return nil, err
				}
				ops = append(ops, m.noContainersOp())
			}

			cleanupOp := newCleanupOp(cleanupDyingMachine, m.doc.Id, force, maxWait)
			ops = append(ops, cleanupOp)

			if canDie {
				ops[0].Assert = append(asserts, advanceLifecycleUnitAsserts(principalUnitNames))
				txnLogger.Debugf(context.TODO(), "txn moving machine %q to %s", m.Id(), life)
				return ops, nil
			}
		}

		if len(m.doc.Principals) > 0 {
			return nil, stateerrors.NewHasAssignedUnitsError(m.doc.Id, m.doc.Principals)
		}
		asserts = append(asserts, noUnitAsserts())

		if life == Dead {
			// A machine may not become Dead until it has no
			// containers.
			if err := m.advanceLifecycleIfNoContainers(); err != nil {
				return nil, err
			}
			ops = append(ops, m.noContainersOp())
			// A machine may not become Dead until it has no more
			// attachments to detachable storage.
			storageAsserts, err := m.assertNoPersistentStorage()
			if err != nil {
				return nil, errors.Trace(err)
			}
			asserts = append(asserts, storageAsserts...)
		}

		// Add the additional asserts needed for this transaction.
		ops[0].Assert = asserts
		return ops, nil
	}

	if err = m.st.db().Run(buildTxn); err == jujutxn.ErrExcessiveContention {
		err = errors.Annotatef(err, "machine %s cannot advance lifecycle", m)
	}
	return err
}

func (m *Machine) advanceLifecyleInitialOps(life Life) []txn.Op {
	return []txn.Op{
		{
			C:      machinesC,
			Id:     m.doc.DocID,
			Update: bson.D{{"$set", bson.D{{"life", life}}}},
		},
	}
}

func controllerAdvanceLifecyleVoteOp() txn.Op {
	return txn.Op{
		C:  controllersC,
		Id: modelGlobalKey,
		Assert: bson.D{
			{"has-vote", bson.M{"$ne": true}},
			{"wants-vote", bson.M{"$ne": true}},
		},
	}
}

// noContainersOp returns an Op to assert that the machine
// has no containers.
func (m *Machine) noContainersOp() txn.Op {
	return txn.Op{
		C:  containerRefsC,
		Id: m.doc.DocID,
		Assert: bson.D{{"$or", []bson.D{
			{{"children", bson.D{{"$size", 0}}}},
			{{"children", bson.D{{"$exists", false}}}},
		}}},
	}
}

// assessCanDieUnit returns true if the machine can die, based on
// evaluating the provided unit.
func (m *Machine) assessCanDieUnit(principalUnit string) (bool, error) {
	canDie := true
	u, err := m.st.Unit(principalUnit)
	if err != nil {
		return false, errors.Annotatef(err, "reading machine %s principal unit %v", m, m.doc.Principals[0])
	}
	app, err := u.application()
	if err != nil {
		return false, errors.Annotatef(err, "reading machine %s principal unit application %v", m, u.doc.Application)
	}
	if u.life() == Alive && app.life() == Alive {
		canDie = false
	}
	return canDie, nil
}

// noUnitAsserts returns bson DocElem which assert that there are
// no units for the machine.
func noUnitAsserts() bson.DocElem {
	return bson.DocElem{
		Name: "$or", Value: []bson.D{
			{{"principals", bson.D{{"$size", 0}}}},
			{{"principals", bson.D{{"$exists", false}}}},
		},
	}
}

// advanceLifecycleUnitAsserts returns bson DocElem which assert that there are
// no units for the machine, or that the list of units has not changed.
func advanceLifecycleUnitAsserts(principalUnitNames []string) bson.DocElem {
	return bson.DocElem{
		Name: "$or", Value: []bson.D{
			{{Name: "principals", Value: principalUnitNames}},
			{{Name: "principals", Value: bson.D{{"$size", 0}}}},
			{{Name: "principals", Value: bson.D{{"$exists", false}}}},
		},
	}
}

// advanceLifecycleIfNoContainers determines if the machine has
// containers, if so, returns the appropriate error.
func (m *Machine) advanceLifecycleIfNoContainers() error {
	containers, err := m.Containers()
	if err != nil {
		return errors.Annotatef(err, "reading machine %s containers", m)
	}

	if len(containers) > 0 {
		return stateerrors.NewHasContainersError(m.doc.Id, containers)
	}
	return nil
}

// assertNoPersistentStorage ensures that there are no persistent volumes or
// filesystems attached to the machine, and returns any mgo/txn assertions
// required to ensure that remains true.
func (m *Machine) assertNoPersistentStorage() (bson.D, error) {
	attachments := names.NewSet()
	for _, v := range m.doc.Volumes {
		tag := names.NewVolumeTag(v)
		detachable, err := isDetachableVolumeTag(m.st.db(), tag)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if detachable {
			attachments.Add(tag)
		}
	}
	for _, f := range m.doc.Filesystems {
		tag := names.NewFilesystemTag(f)
		detachable, err := isDetachableFilesystemTag(m.st.db(), tag)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if detachable {
			attachments.Add(tag)
		}
	}
	if len(attachments) > 0 {
		return nil, stateerrors.NewHasAttachmentsError(m.doc.Id, attachments.SortedValues())
	}
	if m.doc.Life == Dying {
		return nil, nil
	}
	// A Dying machine cannot have attachments added to it,
	// but if we're advancing from Alive to Dead then we
	// must ensure no concurrent attachments are made.
	noNewVolumes := bson.DocElem{
		Name: "volumes", Value: bson.D{{
			"$not", bson.D{{
				"$elemMatch", bson.D{{
					"$nin", m.doc.Volumes,
				}},
			}},
		}},
		// There are no volumes that are not in
		// the set of volumes we previously knew
		// about => the current set of volumes
		// is a subset of the previously known set.
	}
	noNewFilesystems := bson.DocElem{
		Name: "filesystems", Value: bson.D{{
			"$not", bson.D{{
				"$elemMatch", bson.D{{
					"$nin", m.doc.Filesystems,
				}},
			}},
		}},
	}
	return bson.D{noNewVolumes, noNewFilesystems}, nil
}

func (m *Machine) removeOps() ([]txn.Op, error) {
	if m.doc.Life != Dead {
		return nil, fmt.Errorf("machine is not dead")
	}
	ops := []txn.Op{
		{
			C:      machinesC,
			Id:     m.doc.DocID,
			Assert: txn.DocExists,
			Remove: true,
		},
		{
			C:      machinesC,
			Id:     m.doc.DocID,
			Assert: isDeadDoc,
		},
		removeStatusOp(m.st, m.globalKey()),
		removeStatusOp(m.st, m.globalInstanceKey()),
		removeStatusOp(m.st, m.globalModificationKey()),
		removeConstraintsOp(m.globalKey()),
		removeModelMachineRefOp(m.st, m.Id()),
		removeSSHHostKeyOp(m.globalKey()),
	}
	linkLayerDevicesOps, err := m.removeAllLinkLayerDevicesOps()
	if err != nil {
		return nil, errors.Trace(err)
	}
	devicesAddressesOps, err := m.removeAllAddressesOps()
	if err != nil {
		return nil, errors.Trace(err)
	}

	sb, err := NewStorageBackend(m.st)
	if err != nil {
		return nil, errors.Trace(err)
	}
	filesystemOps, err := sb.removeMachineFilesystemsOps(m)
	if err != nil {
		return nil, errors.Trace(err)
	}
	volumeOps, err := sb.removeMachineVolumesOps(m)
	if err != nil {
		return nil, errors.Trace(err)
	}

	ops = append(ops, linkLayerDevicesOps...)
	ops = append(ops, devicesAddressesOps...)
	ops = append(ops, removeContainerRefOps(m.st, m.Id())...)
	ops = append(ops, filesystemOps...)
	ops = append(ops, volumeOps...)
	return ops, nil
}

// Remove removes the machine from state. It will fail if the machine
// is not Dead.
func (m *Machine) Remove() (err error) {
	defer errors.DeferredAnnotatef(&err, "cannot remove machine %s", m.doc.Id)
	logger.Tracef(context.TODO(), "removing machine %q", m.Id())
	// Local variable so we can re-get the machine without disrupting
	// the caller.
	machine := m
	buildTxn := func(attempt int) ([]txn.Op, error) {
		if attempt != 0 {
			machine, err = machine.st.Machine(machine.Id())
			if errors.Is(err, errors.NotFound) {
				// The machine's gone away, that's fine.
				return nil, jujutxn.ErrNoOperations
			}
			if err != nil {
				return nil, errors.Trace(err)
			}
		}
		ops, err := machine.removeOps()
		if err != nil {
			return nil, errors.Trace(err)
		}
		return ops, nil
	}
	return m.st.db().Run(buildTxn)
}

// Refresh refreshes the contents of the machine from the underlying
// state. It returns an error that satisfies errors.IsNotFound if the
// machine has been removed.
func (m *Machine) Refresh() error {
	mdoc, err := m.st.getMachineDoc(m.Id())
	if err != nil {
		if errors.Is(err, errors.NotFound) {
			return err
		}
		return errors.Annotatef(err, "cannot refresh machine %v", m)
	}
	m.doc = *mdoc
	return nil
}

// ApplicationNames returns the names of applications
// represented by units running on the machine.
func (m *Machine) ApplicationNames() ([]string, error) {
	units, err := m.units()
	if err != nil {
		return nil, errors.Trace(err)
	}
	apps := set.NewStrings()
	for _, unit := range units {
		apps.Add(unit.applicationName())
	}
	return apps.SortedValues(), nil
}

// units returns all the units that have been assigned to the machine.
func (m *Machine) units() (units []*Unit, err error) {
	defer errors.DeferredAnnotatef(&err, "cannot get units assigned to machine %v", m)
	unitsCollection, closer := m.st.db().GetCollection(unitsC)
	defer closer()

	pudocs := []unitDoc{}
	err = unitsCollection.Find(bson.D{{"machineid", m.doc.Id}}).All(&pudocs)
	if err != nil {
		return nil, err
	}
	model, err := m.st.Model()
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, pudoc := range pudocs {
		units = append(units, newUnit(m.st, model.Type(), &pudoc))
	}
	return units, nil
}

// SetInstanceInfo is used to provision a machine and in one step sets its
// instance ID, nonce, hardware characteristics, add link-layer devices and set
// their addresses as needed.  After, set charm profiles if needed.
func (m *Machine) SetInstanceInfo(
	id instance.Id, displayName string, nonce string, characteristics *instance.HardwareCharacteristics,
	devicesArgs []LinkLayerDeviceArgs, devicesAddrs []LinkLayerDeviceAddress,
	volumes map[names.VolumeTag]VolumeInfo,
	volumeAttachments map[names.VolumeTag]VolumeAttachmentInfo,
	charmProfiles []string,
) error {
	logger.Tracef(context.TODO(),
		"setting instance info: machine %v, deviceAddrs: %#v, devicesArgs: %#v",
		m.Id(), devicesAddrs, devicesArgs)

	sb, err := NewStorageBackend(m.st)
	if err != nil {
		return errors.Trace(err)
	}

	// Record volumes and volume attachments, and set the initial
	// status: attached or attaching.
	if err := setProvisionedVolumeInfo(sb, volumes); err != nil {
		return errors.Trace(err)
	}
	if err := setMachineVolumeAttachmentInfo(sb, m.Id(), volumeAttachments); err != nil {
		return errors.Trace(err)
	}
	volumeStatus := make(map[names.VolumeTag]status.Status)
	for tag := range volumes {
		volumeStatus[tag] = status.Attaching
	}
	for tag := range volumeAttachments {
		volumeStatus[tag] = status.Attached
	}
	for tag, volStatus := range volumeStatus {
		vol, err := sb.Volume(tag)
		if err != nil {
			return errors.Trace(err)
		}
		if err := vol.SetStatus(status.StatusInfo{
			Status: volStatus,
		}); err != nil {
			return errors.Annotatef(
				err, "setting status of %s", names.ReadableString(tag),
			)
		}
	}

	return nil
}

// Addresses returns any hostnames and ips associated with a machine,
// determined both by the machine itself, and by asking the provider.
//
// The addresses returned by the provider shadow any of the addresses
// that the machine reported with the same address value.
// Provider-reported addresses always come before machine-reported
// addresses. Duplicates are removed.
func (m *Machine) Addresses() (addresses network.SpaceAddresses) {
	return network.MergedAddresses(networkAddresses(m.doc.MachineAddresses), networkAddresses(m.doc.Addresses))
}

// PublicAddress returns a public address for the machine. If no address is
// available it returns an error that satisfies network.IsNoAddressError().
func (m *Machine) PublicAddress() (network.SpaceAddress, error) {
	publicAddress := m.doc.PreferredPublicAddress.networkAddress()
	var err error
	if publicAddress.Value == "" {
		err = network.NoAddressError("public")
	}
	return publicAddress, err
}

// PrivateAddress returns a private address for the machine. If no address is
// available it returns an error that satisfies network.IsNoAddressError().
func (m *Machine) PrivateAddress() (network.SpaceAddress, error) {
	privateAddress := m.doc.PreferredPrivateAddress.networkAddress()
	var err error
	if privateAddress.Value == "" {
		err = network.NoAddressError("private")
	}
	return privateAddress, err
}

// String returns a unique description of this machine.
func (m *Machine) String() string {
	return m.doc.Id
}

// Placement returns the machine's Placement structure that should be used when
// provisioning an instance for the machine.
func (m *Machine) Placement() string {
	return m.doc.Placement
}

// Constraints returns the exact constraints that should apply when provisioning
// an instance for the machine.
func (m *Machine) Constraints() (constraints.Value, error) {
	return readConstraints(m.st, m.globalKey())
}

// Clean returns true if the machine does not have any deployed units or containers.
func (m *Machine) Clean() bool {
	return m.doc.Clean
}

// VolumeAttachments returns the machine's volume attachments.
func (m *Machine) VolumeAttachments() ([]VolumeAttachment, error) {
	sb, err := NewStorageBackend(m.st)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return sb.MachineVolumeAttachments(m.MachineTag())
}

// PrepareActionPayload returns the payload to use in creating an action for this machine.
// Note that the use of spec.InsertDefaults mutates payload.
func (m *Machine) PrepareActionPayload(name string, payload map[string]interface{}, parallel *bool, executionGroup *string) (map[string]interface{}, bool, string, error) {
	if len(name) == 0 {
		return nil, false, "", errors.New("no action name given")
	}

	spec, ok := actions.PredefinedActionsSpec[name]
	if !ok {
		return nil, false, "", errors.Errorf("cannot add action %q to a machine; only predefined actions allowed", name)
	}

	// Reject bad payloads before attempting to insert defaults.
	err := spec.ValidateParams(payload)
	if err != nil {
		return nil, false, "", errors.Trace(err)
	}
	payloadWithDefaults, err := spec.InsertDefaults(payload)
	if err != nil {
		return nil, false, "", errors.Trace(err)
	}

	runParallel := spec.Parallel
	if parallel != nil {
		runParallel = *parallel
	}
	runExecutionGroup := spec.ExecutionGroup
	if executionGroup != nil {
		runExecutionGroup = *executionGroup
	}

	return payloadWithDefaults, runParallel, runExecutionGroup, nil
}

// CancelAction is part of the ActionReceiver interface.
func (m *Machine) CancelAction(action Action) (Action, error) {
	return action.Finish(ActionResults{Status: ActionCancelled})
}

// WatchActionNotifications is part of the ActionReceiver interface.
func (m *Machine) WatchActionNotifications() StringsWatcher {
	return m.st.watchActionNotificationsFilteredBy(m)
}

// WatchPendingActionNotifications is part of the ActionReceiver interface.
func (m *Machine) WatchPendingActionNotifications() StringsWatcher {
	return m.st.watchEnqueuedActionsFilteredBy(m)
}

// Actions is part of the ActionReceiver interface.
func (m *Machine) Actions() ([]Action, error) {
	return m.st.matchingActions(m)
}

// CompletedActions is part of the ActionReceiver interface.
func (m *Machine) CompletedActions() ([]Action, error) {
	return m.st.matchingActionsCompleted(m)
}

// PendingActions is part of the ActionReceiver interface.
func (m *Machine) PendingActions() ([]Action, error) {
	return m.st.matchingActionsPending(m)
}

// RunningActions is part of the ActionReceiver interface.
func (m *Machine) RunningActions() ([]Action, error) {
	return m.st.matchingActionsRunning(m)
}

// RecordAgentStartInformation updates the host name (if non-empty) reported
// by the machine agent and sets the agent start time to the current time.
func (m *Machine) RecordAgentStartInformation(hostname string) error {
	now := m.st.clock().Now()
	update := bson.D{
		{"agent-started-at", now},
	}

	if hostname != "" {
		update = append(update, bson.DocElem{"hostname", hostname})
	}

	ops := []txn.Op{{
		C:      machinesC,
		Id:     m.doc.DocID,
		Assert: notDeadDoc,
		Update: bson.D{{"$set", update}},
	}}

	if err := m.st.db().RunTransaction(ops); err != nil {
		// If instance doc doesn't exist, that's ok; there's nothing to keep,
		// but that's not an error we care about.
		return errors.Annotatef(onAbort(err, nil), "cannot update agent hostname/start time for machine %q", m)
	}
	m.doc.AgentStartedAt = now
	if hostname != "" {
		m.doc.Hostname = hostname
	}
	return nil
}

// AgentStartTime returns the last recorded timestamp when the machine agent
// was started.
func (m *Machine) AgentStartTime() time.Time {
	return m.doc.AgentStartedAt
}

// Hostname returns the hostname reported by the machine agent.
func (m *Machine) Hostname() string {
	return m.doc.Hostname
}
