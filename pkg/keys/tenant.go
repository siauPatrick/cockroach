// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package keys

import (
	"bytes"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/util/encoding"
	"github.com/pkg/errors"
)

// MakeTenantPrefix creates the key prefix associated with the specified tenant.
func MakeTenantPrefix(tenID roachpb.TenantID) roachpb.Key {
	if tenID == roachpb.SystemTenantID {
		return nil
	}
	return encoding.EncodeUvarintAscending(tenantPrefix, tenID.ToUint64())
}

// DecodeTenantPrefix determines the tenant ID from the key prefix, returning
// the remainder of the key (with the prefix removed) and the decoded tenant ID.
func DecodeTenantPrefix(key roachpb.Key) ([]byte, roachpb.TenantID, error) {
	if len(key) == 0 { // key.Equal(roachpb.RKeyMin)
		return nil, roachpb.SystemTenantID, nil
	}
	if key[0] != tenantPrefixByte {
		return key, roachpb.SystemTenantID, nil
	}
	rem, tenID, err := encoding.DecodeUvarintAscending(key[1:])
	if err != nil {
		return nil, roachpb.TenantID{}, err
	}
	return rem, roachpb.MakeTenantID(tenID), nil
}

// TenantIDKeyGen provides methods for generating SQL keys bound to a given
// tenant. The generator also provides methods for efficiently decoding keys
// previously generated by it. The generated keys are safe to use indefinitely
// and the generator is safe to use concurrently.
//
// The type is expressed as a pointer to a slice instead of a slice directly so
// that its zero value is not usable. Any attempt to use the methods on the zero
// value of a TenantIDKeyGen will panic.
type TenantIDKeyGen struct {
	buf *roachpb.Key
}

// SystemTenantKeyGen is a SQL key generator for the system tenant.
var SystemTenantKeyGen = MakeTenantIDKeyGen(roachpb.SystemTenantID)

// TODOTenantKeyGen is a SQL key generator. It is equivalent to
// SystemTenantKeyGen, but should be used when it is unclear which
// tenant should be referenced by the surrounding context.
var TODOTenantKeyGen = MakeTenantIDKeyGen(roachpb.SystemTenantID)

// MakeTenantIDKeyGen creates a new tenant ID key generator suitable for
// constructing various SQL keys.
func MakeTenantIDKeyGen(tenID roachpb.TenantID) TenantIDKeyGen {
	k := MakeTenantPrefix(tenID)
	k = k[:len(k):len(k)] // bound capacity, avoid aliasing
	return TenantIDKeyGen{&k}
}

// TenantPrefix returns the key prefix used for the tenants's data.
func (g TenantIDKeyGen) TenantPrefix() roachpb.Key {
	return *g.buf
}

// TablePrefix returns the key prefix used for the table's data.
func (g TenantIDKeyGen) TablePrefix(tableID uint32) roachpb.Key {
	k := g.TenantPrefix()
	return encoding.EncodeUvarintAscending(k, uint64(tableID))
}

// IndexPrefix returns the key prefix used for the index's data.
func (g TenantIDKeyGen) IndexPrefix(tableID, indexID uint32) roachpb.Key {
	k := g.TablePrefix(tableID)
	return encoding.EncodeUvarintAscending(k, uint64(indexID))
}

// DescMetadataPrefix returns the key prefix for all descriptors.
func (g TenantIDKeyGen) DescMetadataPrefix() roachpb.Key {
	return g.IndexPrefix(DescriptorTableID, DescriptorTablePrimaryKeyIndexID)
}

// DescMetadataKey returns the key for the descriptor.
func (g TenantIDKeyGen) DescMetadataKey(descID uint32) roachpb.Key {
	k := g.DescMetadataPrefix()
	k = encoding.EncodeUvarintAscending(k, uint64(descID))
	return MakeFamilyKey(k, DescriptorTableDescriptorColFamID)
}

// SequenceKey returns the key used to store the value of a sequence.
func (g TenantIDKeyGen) SequenceKey(tableID uint32) roachpb.Key {
	k := g.IndexPrefix(tableID, SequenceIndexID)
	k = encoding.EncodeUvarintAscending(k, 0)    // Primary key value
	k = MakeFamilyKey(k, SequenceColumnFamilyID) // Column family
	return k
}

// ZoneKeyPrefix returns the key prefix for id's row in the system.zones table.
func (g TenantIDKeyGen) ZoneKeyPrefix(id uint32) roachpb.Key {
	k := g.IndexPrefix(ZonesTableID, ZonesTablePrimaryIndexID)
	return encoding.EncodeUvarintAscending(k, uint64(id))
}

// ZoneKey returns the key for id's entry in the system.zones table.
func (g TenantIDKeyGen) ZoneKey(id uint32) roachpb.Key {
	k := g.ZoneKeyPrefix(id)
	return MakeFamilyKey(k, uint32(ZonesTableConfigColumnID))
}

// StripTenantPrefix validates that the given key has the proper tenant ID
// prefix, returning the remainder of the key with the prefix removed. The
// method returns an error if the key has a different tenant ID prefix than
// would be generated by the generator.
func (g TenantIDKeyGen) StripTenantPrefix(key roachpb.Key) ([]byte, error) {
	tenPrefix := g.TenantPrefix()
	if !bytes.HasPrefix(key, tenPrefix) {
		return nil, errors.Errorf("invalid tenant id prefix: %q", key)
	}
	return key[len(tenPrefix):], nil
}

// DecodeTablePrefix validates that the given key has a table prefix, returning
// the remainder of the key (with the prefix removed) and the decoded descriptor
// ID of the table.
func (g TenantIDKeyGen) DecodeTablePrefix(key roachpb.Key) ([]byte, uint32, error) {
	key, err := g.StripTenantPrefix(key)
	if err != nil {
		return nil, 0, err
	}
	if encoding.PeekType(key) != encoding.Int {
		return nil, 0, errors.Errorf("invalid key prefix: %q", key)
	}
	key, tableID, err := encoding.DecodeUvarintAscending(key)
	return key, uint32(tableID), err
}

// DecodeIndexPrefix validates that the given key has a table ID followed by an
// index ID, returning the remainder of the key (with the table and index prefix
// removed) and the decoded IDs of the table and index, respectively.
func (g TenantIDKeyGen) DecodeIndexPrefix(key roachpb.Key) ([]byte, uint32, uint32, error) {
	key, tableID, err := g.DecodeTablePrefix(key)
	if err != nil {
		return nil, 0, 0, err
	}
	if encoding.PeekType(key) != encoding.Int {
		return nil, 0, 0, errors.Errorf("invalid key prefix: %q", key)
	}
	key, indexID, err := encoding.DecodeUvarintAscending(key)
	return key, tableID, uint32(indexID), err
}

// DecodeDescMetadataID decodes a descriptor ID from a descriptor metadata key.
func (g TenantIDKeyGen) DecodeDescMetadataID(key roachpb.Key) (uint64, error) {
	// Extract table and index ID from key.
	remaining, tableID, _, err := g.DecodeIndexPrefix(key)
	if err != nil {
		return 0, err
	}
	if tableID != DescriptorTableID {
		return 0, errors.Errorf("key is not a descriptor table entry: %v", key)
	}
	// Extract the descriptor ID.
	_, id, err := encoding.DecodeUvarintAscending(remaining)
	if err != nil {
		return 0, err
	}
	return id, nil
}
