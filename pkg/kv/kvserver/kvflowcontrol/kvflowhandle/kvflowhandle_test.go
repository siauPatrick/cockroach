// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package kvflowhandle_test

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvflowcontrol"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvflowcontrol/kvflowcontroller"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvflowcontrol/kvflowcontrolpb"
	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/kvflowcontrol/kvflowhandle"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings/cluster"
	"github.com/cockroachdb/cockroach/pkg/util/admission/admissionpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/metric"
	"github.com/stretchr/testify/require"
)

// TestHandleAdmit tests the blocking behavior of Handle.Admit():
// - we block until there are flow tokens available;
// - we unblock when streams without flow tokens are disconnected;
// - we unblock when the handle is closed;
// - we unblock when the handle is reset.
func TestHandleAdmit(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	stream := kvflowcontrol.Stream{TenantID: roachpb.MustMakeTenantID(42), StoreID: roachpb.StoreID(42)}
	pos := func(d uint64) kvflowcontrolpb.RaftLogPosition {
		return kvflowcontrolpb.RaftLogPosition{Term: 1, Index: d}
	}

	for _, tc := range []struct {
		name      string
		unblockFn func(context.Context, kvflowcontrol.Handle)
	}{
		{
			name: "blocks-for-tokens",
			unblockFn: func(ctx context.Context, handle kvflowcontrol.Handle) {
				// Return tokens tied to pos=1 (16MiB worth); the call to
				// .Admit() should unblock.
				handle.ReturnTokensUpto(ctx, admissionpb.NormalPri, pos(1), stream)
			},
		},
		{
			name: "unblocked-when-stream-disconnects",
			unblockFn: func(ctx context.Context, handle kvflowcontrol.Handle) {
				// Disconnect the stream; the call to .Admit() should unblock.
				handle.DisconnectStream(ctx, stream)
			},
		},
		{
			name: "unblocked-when-closed",
			unblockFn: func(ctx context.Context, handle kvflowcontrol.Handle) {
				// Close the handle; the call to .Admit() should unblock.
				handle.Close(ctx)
			},
		},
		{
			name: "unblocked-when-reset",
			unblockFn: func(ctx context.Context, handle kvflowcontrol.Handle) {
				// Reset all streams on the handle; the call to .Admit() should
				// unblock.
				handle.ResetStreams(ctx)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			registry := metric.NewRegistry()
			clock := hlc.NewClockForTesting(nil)
			st := cluster.MakeTestingClusterSettings()
			kvflowcontrol.Enabled.Override(ctx, &st.SV, true)
			kvflowcontrol.Mode.Override(ctx, &st.SV, int64(kvflowcontrol.ApplyToAll))

			controller := kvflowcontroller.New(registry, st, clock)
			handle := kvflowhandle.New(
				controller,
				kvflowhandle.NewMetrics(registry),
				clock,
				roachpb.RangeID(1),
				roachpb.SystemTenantID,
			)

			// Connect a single stream at pos=0 and deplete all 16MiB of regular
			// tokens at pos=1.
			handle.ConnectStream(ctx, pos(0), stream)
			handle.DeductTokensFor(ctx, admissionpb.NormalPri, pos(1), kvflowcontrol.Tokens(16<<20 /* 16MiB */))

			// Invoke .Admit() in a separate goroutine, and test below whether
			// the goroutine is blocked.
			admitCh := make(chan struct{})
			go func() {
				require.NoError(t, handle.Admit(ctx, admissionpb.NormalPri, time.Time{}))
				close(admitCh)
			}()

			select {
			case <-admitCh:
				t.Fatalf("unexpectedly admitted")
			case <-time.After(10 * time.Millisecond):
			}

			tc.unblockFn(ctx, handle)

			select {
			case <-admitCh:
			case <-time.After(5 * time.Second):
				t.Fatalf("didn't get admitted")
			}
		})
	}
}

func TestFlowControlMode(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	ctx := context.Background()
	stream := kvflowcontrol.Stream{
		TenantID: roachpb.MustMakeTenantID(42),
		StoreID:  roachpb.StoreID(42),
	}
	pos := func(d uint64) kvflowcontrolpb.RaftLogPosition {
		return kvflowcontrolpb.RaftLogPosition{Term: 1, Index: d}
	}

	for _, tc := range []struct {
		mode            kvflowcontrol.ModeT
		blocks, ignores []admissionpb.WorkClass
	}{
		{
			mode: kvflowcontrol.ApplyToElastic,
			blocks: []admissionpb.WorkClass{
				admissionpb.ElasticWorkClass,
			},
			ignores: []admissionpb.WorkClass{
				admissionpb.RegularWorkClass,
			},
		},
		{
			mode: kvflowcontrol.ApplyToAll,
			blocks: []admissionpb.WorkClass{
				admissionpb.ElasticWorkClass, admissionpb.RegularWorkClass,
			},
			ignores: []admissionpb.WorkClass{},
		},
	} {
		t.Run(tc.mode.String(), func(t *testing.T) {
			registry := metric.NewRegistry()
			clock := hlc.NewClockForTesting(nil)
			st := cluster.MakeTestingClusterSettings()
			kvflowcontrol.Enabled.Override(ctx, &st.SV, true)
			kvflowcontrol.Mode.Override(ctx, &st.SV, int64(tc.mode))

			controller := kvflowcontroller.New(registry, st, clock)
			handle := kvflowhandle.New(
				controller,
				kvflowhandle.NewMetrics(registry),
				clock,
				roachpb.RangeID(1),
				roachpb.SystemTenantID,
			)
			defer handle.Close(ctx)

			// Connect a single stream at pos=0 and deplete all 16MiB of regular
			// tokens at pos=1. It also puts elastic tokens in the -ve.
			handle.ConnectStream(ctx, pos(0), stream)
			handle.DeductTokensFor(ctx, admissionpb.NormalPri, pos(1), kvflowcontrol.Tokens(16<<20 /* 16MiB */))

			// Invoke .Admit() for {regular,elastic} work in a separate
			// goroutines, and test below whether the goroutines are blocked.
			regularAdmitCh := make(chan struct{})
			elasticAdmitCh := make(chan struct{})
			go func() {
				require.NoError(t, handle.Admit(ctx, admissionpb.NormalPri, time.Time{}))
				close(regularAdmitCh)
			}()
			go func() {
				require.NoError(t, handle.Admit(ctx, admissionpb.BulkNormalPri, time.Time{}))
				close(elasticAdmitCh)
			}()

			for _, ignoredClass := range tc.ignores { // work should not block
				classAdmitCh := regularAdmitCh
				if ignoredClass == admissionpb.ElasticWorkClass {
					classAdmitCh = elasticAdmitCh
				}

				select {
				case <-classAdmitCh:
				case <-time.After(5 * time.Second):
					t.Fatalf("%s work didn't get admitted", ignoredClass)
				}
			}

			for _, blockedClass := range tc.blocks { // work should get blocked
				classAdmitCh := regularAdmitCh
				if blockedClass == admissionpb.ElasticWorkClass {
					classAdmitCh = elasticAdmitCh
				}

				select {
				case <-classAdmitCh:
					t.Fatalf("unexpectedly admitted %s work", blockedClass)
				case <-time.After(10 * time.Millisecond):
				}
			}
		})
	}

}
