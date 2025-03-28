// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"context"
	"fmt"
	"slices"
	"strings"

	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/model"
	modeltesting "github.com/juju/juju/core/model/testing"
	"github.com/juju/juju/core/user"
	usertesting "github.com/juju/juju/core/user/testing"
	accesserrors "github.com/juju/juju/domain/access/errors"
	accessstate "github.com/juju/juju/domain/access/state"
	"github.com/juju/juju/domain/keymanager"
	keyerrors "github.com/juju/juju/domain/keymanager/errors"
	modelerrors "github.com/juju/juju/domain/model/errors"
	modelstate "github.com/juju/juju/domain/model/state"
	statemodeltesting "github.com/juju/juju/domain/model/state/testing"
	schematesting "github.com/juju/juju/domain/schema/testing"
	"github.com/juju/juju/internal/ssh"
)

type stateSuite struct {
	schematesting.ControllerSuite

	userId   user.UUID
	userName user.Name
	modelId  model.UUID
}

var _ = gc.Suite(&stateSuite{})

var (
	testingPublicKeys = []string{
		// ecdsa testing public key
		"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBG00bYFLb/sxPcmVRMg8NXZK/ldefElAkC9wD41vABdHZiSRvp+2y9BMNVYzE/FnzKObHtSvGRX65YQgRn7k5p0= juju1@example.com",

		// ed25519 testing public key
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju2@example.com",

		// rsa testing public key
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDvplNOK3UBpULZKvZf/I5JHci/DufpSxj8yR4yKE2grescJxu6754jPT3xztSeLGD31/oJApJZGkMUAMRenvDqIaq+taRfOUo/l19AlGZc+Edv4bTlJzZ1Lzwex1vvL1doaLb/f76IIUHClGUgIXRceQH1ovHiIWj6nGltuLanG8YTWxlzzK33yhitmZt142DmpX1VUVF5c/Hct6Rav5lKmwej1TDed1KmHzXVoTHEsmWhKsOK27ue5yTuq0GX6LrAYDucF+2MqZCsuddXsPAW1tj5GNZSR7RrKW5q1CI0G7k9gSomuCsRMlCJ3BqID/vUSs/0qOWg4he0HUsYKQSrXIhckuZu+jYP8B80MoXT50ftRidoG/zh/PugBdXTk46FloVClQopG5A2fbqrphADcUUbRUxZ2lWQN+OVHKfEsfV2b8L2aSqZUGlryfW1cirB5JCTDvtv7rUy9/ny9iKA+8tAyKSDF0I901RDDqKc9dSkrHCg2bLnJZDoiRoWczE= juju3@example.com",
	}
)

func generatePublicKeys(c *gc.C, publicKeys []string) []keymanager.PublicKey {
	rval := make([]keymanager.PublicKey, 0, len(publicKeys))
	for _, pk := range publicKeys {
		parsedKey, err := ssh.ParsePublicKey(pk)
		c.Assert(err, jc.ErrorIsNil)

		rval = append(rval, keymanager.PublicKey{
			Comment:         parsedKey.Comment,
			FingerprintHash: keymanager.FingerprintHashAlgorithmSHA256,
			Fingerprint:     parsedKey.Fingerprint(),
			Key:             pk,
		})
	}

	return rval
}

func (s *stateSuite) SetUpTest(c *gc.C) {
	s.ControllerSuite.SetUpTest(c)
	s.SeedControllerUUID(c)

	s.modelId = statemodeltesting.CreateTestModel(c, s.TxnRunnerFactory(), "keys")

	model, err := modelstate.NewState(s.TxnRunnerFactory()).GetModel(
		context.Background(), s.modelId,
	)
	c.Assert(err, jc.ErrorIsNil)
	s.userId = model.Owner
	s.userName = model.OwnerName
}

// TestAddPublicKeyForUser is asserting the happy path of adding a public key
// for a user. Specifically we want to see that inserting the same key across
// multiple models doesn't result in constraint violations for the users public
// ssh keys.
func (s *stateSuite) TestAddPublicKeyForUser(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)

	// Create a second model to add keys onto
	modelId := statemodeltesting.CreateTestModel(c, s.TxnRunnerFactory(), "second-model")

	// Confirm that the users public ssh keys don't show up on the second model
	// yet
	keys, err = state.GetPublicKeysDataForUser(context.Background(), modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	c.Check(len(keys), gc.Equals, 0)

	// Add the users keys onto the second model. We want to see here that this
	// is a successful operation with no errors.
	err = state.AddPublicKeysForUser(context.Background(), modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	// Confirm the keys exists on the second model
	keys, err = state.GetPublicKeysDataForUser(context.Background(), modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)
}

// TestAddPublicKeysForUserAlreadyExists is asserting that if we try and add the
// same public key for a user more then once to a model we get back an error
// that satisfies [keyerrors.PublicKeyAlreadyExists].
func (s *stateSuite) TestAddPublicKeyForUserAlreadyExists(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)

	// Add the users keys onto the second model. We want to see here that this
	// is a successful operation with no errors.
	err = state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIs, keyerrors.PublicKeyAlreadyExists)

	// Confirm the key still exists on the model
	keys, err = state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)
}

// TestAddPublicKeyForUserNotFound is asserting that if we attempt to add a
// public key to a model for a user that doesn't exist we get back a
// [accesserrors.UserNotFound] error.
func (s *stateSuite) TestAddPublicKeyForUserNotFound(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	badUserId := usertesting.GenUserUUID(c)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, badUserId, keysToAdd)
	c.Check(err, jc.ErrorIs, accesserrors.UserNotFound)
}

// TestAddPublicKeyForUserOnNotFoundModel is asserting that if we attempt to add
// a public key for a user on a model that does not exist we get back a
// [modelerrors.NotFound] error.
func (s *stateSuite) TestAddPublicKeyForUserOnNotFoundModel(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	badModelId := modeltesting.GenModelUUID(c)

	err := state.AddPublicKeysForUser(context.Background(), badModelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIs, modelerrors.NotFound)
}

// TestEnsurePublicKeysForUser is asserting the happy path of
// [State.EnsurePublicKeysForUser].
func (s *stateSuite) TestEnsurePublicKeysForUser(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.EnsurePublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)

	// Run all of the operations again and confirm that there exists no errors.
	err = state.EnsurePublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	keys, err = state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)
}

// TestEnsurePublicKeysForUser is asserting the ensure user after keys have
// been stripped of the comments. This should ensure that we're checking against
// the fingerprint and not the public key.
// [State.EnsurePublicKeysForUser].
func (s *stateSuite) TestEnsurePublicKeysForUserForStrippedComments(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.EnsurePublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)

	// Run all of the operations again and confirm that there exists no errors.

	stripped := make([]keymanager.PublicKey, len(keysToAdd))
	for i, key := range keysToAdd {

		newKey := key.Key
		if parts := strings.Split(key.Key, " "); len(parts) > 2 {
			newKey = fmt.Sprintf("%s %s", parts[0], parts[1])
		}

		stripped[i] = keymanager.PublicKey{
			Comment:         key.Comment,
			FingerprintHash: key.FingerprintHash,
			Fingerprint:     key.Fingerprint,
			Key:             newKey,
		}
	}

	err = state.EnsurePublicKeysForUser(context.Background(), s.modelId, s.userId, stripped)
	c.Check(err, jc.ErrorIsNil)

	keys, err = state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(keys, jc.DeepEquals, testingPublicKeys)
}

// TestEnsurePublicKeyForUserNotFound is asserting that if we attempt to add a
// public key to a model for a user that doesn't exist we get back a
// [accesserrors.UserNotFound] error.
func (s *stateSuite) TestEnsurePublicKeyForUserNotFound(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	badUserId := usertesting.GenUserUUID(c)

	err := state.EnsurePublicKeysForUser(context.Background(), s.modelId, badUserId, keysToAdd)
	c.Check(err, jc.ErrorIs, accesserrors.UserNotFound)
}

// TestEnsurePublicKeyForUserOnNotFoundModel is asserting that if we attempt to
// add a public key for a user on a model that does not exist we get back a
// [modelerrors.NotFound] error.
func (s *stateSuite) TestEnsurePublicKeyForUserOnNotFoundModel(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	badModelId := modeltesting.GenModelUUID(c)

	err := state.EnsurePublicKeysForUser(context.Background(), badModelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIs, modelerrors.NotFound)
}

// TestDeletePublicKeysForNonExistentUser is asserting that if we try and
// delete public keys for a user that doesn't exist we get an
// [accesserrors.UserNotFound] error
func (s *stateSuite) TestDeletePublicKeysForNonExistentUser(c *gc.C) {
	userId := usertesting.GenUserUUID(c)
	state := NewState(s.TxnRunnerFactory())
	err := state.DeletePublicKeysForUser(context.Background(), s.modelId, userId, []string{"comment"})
	c.Check(err, jc.ErrorIs, accesserrors.UserNotFound)
}

// TestDeletePublicKeysForComment is testing that we can remove a users public
// keys via the comment string.
func (s *stateSuite) TestDeletePublicKeysForComment(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	err = state.DeletePublicKeysForUser(context.Background(), s.modelId, s.userId, []string{
		keysToAdd[0].Comment,
	})
	c.Assert(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(testingPublicKeys[1:], jc.DeepEquals, keys)
}

// TestDeletePublicKeysForComment is testing that we can remove a users public
// keys via the fingerprint.
func (s *stateSuite) TestDeletePublicKeysForFingerprint(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	err = state.DeletePublicKeysForUser(context.Background(), s.modelId, s.userId, []string{
		keysToAdd[0].Fingerprint,
	})
	c.Assert(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(testingPublicKeys[1:], jc.DeepEquals, keys)
}

// TestDeletePublicKeysForComment is testing that we can remove a users public
// keys via the keys data.
func (s *stateSuite) TestDeletePublicKeysForKeyData(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	err = state.DeletePublicKeysForUser(context.Background(), s.modelId, s.userId, []string{
		keysToAdd[0].Key,
	})
	c.Assert(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(testingPublicKeys[1:], jc.DeepEquals, keys)
}

// TestDeletePublicKeysForCombination is asserting that we can remove a users
// public keys via a combination of fingerprint and comment.
func (s *stateSuite) TestDeletePublicKeysForCombination(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	err = state.DeletePublicKeysForUser(context.Background(), s.modelId, s.userId, []string{
		keysToAdd[0].Comment,
		keysToAdd[1].Fingerprint,
	})
	c.Assert(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(testingPublicKeys[2:], jc.DeepEquals, keys)
}

// TestDeleteSamePublicKeyByTwoMethods is here to assert that if we call one
// delete operation with both a fingerprint and a comment for the same key only
// that key is removed and no other keys are removed and no other errors happen.
func (s *stateSuite) TestDeleteSamePublicKeyByTwoMethods(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(context.Background(), s.modelId, s.userId, keysToAdd)
	c.Check(err, jc.ErrorIsNil)

	err = state.DeletePublicKeysForUser(context.Background(), s.modelId, s.userId, []string{
		keysToAdd[0].Comment,
		keysToAdd[0].Fingerprint,
	})
	c.Assert(err, jc.ErrorIsNil)

	keys, err := state.GetPublicKeysDataForUser(context.Background(), s.modelId, s.userId)
	c.Assert(err, jc.ErrorIsNil)
	slices.Sort(keys)
	slices.Sort(testingPublicKeys)
	c.Check(testingPublicKeys[1:], jc.DeepEquals, keys)
}

// TestDeletePublicKeysForNonExistentModel is asserting the if we try and delete
// user keys off of a model that doesn't exist we get back a
// [modelerrors.NotFound] error.
func (s *stateSuite) TestDeletePublicKeysForNonExistentModel(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	badModelId := modeltesting.GenModelUUID(c)

	err := state.DeletePublicKeysForUser(context.Background(), badModelId, s.userId, []string{
		keysToAdd[0].Comment,
		keysToAdd[0].Fingerprint,
	})
	c.Check(err, jc.ErrorIs, modelerrors.NotFound)
}

// TestGetAllUsersPublicKeys is responsible for testing the happy path of
// getting all user keys in the model.
func (s *stateSuite) TestGetAllUsersPublicKeys(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	keysToAdd := generatePublicKeys(c, testingPublicKeys)

	err := state.AddPublicKeysForUser(
		context.Background(),
		s.modelId,
		s.userId,
		keysToAdd,
	)
	c.Assert(err, jc.ErrorIsNil)

	secondUserId := usertesting.GenUserUUID(c)
	secondUserName := usertesting.GenNewName(c, "tlm")
	userSt := accessstate.NewUserState(s.TxnRunnerFactory())
	err = userSt.AddUser(
		context.Background(),
		secondUserId,
		secondUserName,
		"tlm",
		false,
		s.userId,
	)
	c.Assert(err, jc.ErrorIsNil)

	err = state.AddPublicKeysForUser(
		context.Background(),
		s.modelId,
		secondUserId,
		keysToAdd,
	)
	c.Assert(err, jc.ErrorIsNil)

	allKeys, err := state.GetAllUsersPublicKeys(context.Background(), s.modelId)
	c.Check(err, jc.ErrorIsNil)

	for k := range allKeys {
		slices.Sort(allKeys[k])
	}
	expected := []string{
		"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBG00bYFLb/sxPcmVRMg8NXZK/ldefElAkC9wD41vABdHZiSRvp+2y9BMNVYzE/FnzKObHtSvGRX65YQgRn7k5p0= juju1@example.com",
		"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIN8h8XBpjS9aBUG5cdoSWubs7wT2Lc/BEZIUQCqoaOZR juju2@example.com",
		"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDvplNOK3UBpULZKvZf/I5JHci/DufpSxj8yR4yKE2grescJxu6754jPT3xztSeLGD31/oJApJZGkMUAMRenvDqIaq+taRfOUo/l19AlGZc+Edv4bTlJzZ1Lzwex1vvL1doaLb/f76IIUHClGUgIXRceQH1ovHiIWj6nGltuLanG8YTWxlzzK33yhitmZt142DmpX1VUVF5c/Hct6Rav5lKmwej1TDed1KmHzXVoTHEsmWhKsOK27ue5yTuq0GX6LrAYDucF+2MqZCsuddXsPAW1tj5GNZSR7RrKW5q1CI0G7k9gSomuCsRMlCJ3BqID/vUSs/0qOWg4he0HUsYKQSrXIhckuZu+jYP8B80MoXT50ftRidoG/zh/PugBdXTk46FloVClQopG5A2fbqrphADcUUbRUxZ2lWQN+OVHKfEsfV2b8L2aSqZUGlryfW1cirB5JCTDvtv7rUy9/ny9iKA+8tAyKSDF0I901RDDqKc9dSkrHCg2bLnJZDoiRoWczE= juju3@example.com",
	}
	slices.Sort(expected)

	c.Check(allKeys, jc.DeepEquals, map[user.Name][]string{
		s.userName:     expected,
		secondUserName: expected,
	})
}

// TestGetAllUserPublicKeysEmpty is asserting that if there exists no public
// keys for any user in the model and we call [State.GetAllUsersPublicKeys] we
// get back an empty map and no errors.
func (s *stateSuite) TestGetAllUserPublicKeysEmpty(c *gc.C) {
	state := NewState(s.TxnRunnerFactory())
	allKeys, err := state.GetAllUsersPublicKeys(context.Background(), s.modelId)
	c.Check(err, jc.ErrorIsNil)
	c.Check(len(allKeys), gc.Equals, 0)
}

// TestGetAllUserPublicKeysModelNotFound is asserting that is we ask for all the
// user public keys on a model that does not exist we get back a
// [modelerrors.NotFound] error.
func (s *stateSuite) TestGetAllUserPublicKeysModelNotFound(c *gc.C) {
	badModelUUID := modeltesting.GenModelUUID(c)
	_, err := NewState(s.TxnRunnerFactory()).GetAllUsersPublicKeys(
		context.Background(),
		badModelUUID,
	)
	c.Check(err, jc.ErrorIs, modelerrors.NotFound)
}
