// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package forkid

import (
	"bytes"
	"math"
	"testing"

	"github.com/dominant-strategies/go-quai/common"
	"github.com/dominant-strategies/go-quai/params"
	"github.com/dominant-strategies/go-quai/rlp"
)

// TestCreation tests that different genesis and fork rule combinations result in
// the correct fork ID.
func TestCreation(t *testing.T) {
	type testcase struct {
		head uint64
		want ID
	}
	tests := []struct {
		config  *params.ChainConfig
		genesis common.Hash
		cases   []testcase
	}{
		// Mainnet test cases
		{
			params.MainnetChainConfig,
			params.MainnetGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xfc64ec04), Next: 1150000}},         // Unsynced
				{7280000, ID{Hash: checksumToBytes(0x668db0af), Next: 9069000}},   // First and last Constantinople, first Petersburg block
				{9068999, ID{Hash: checksumToBytes(0x668db0af), Next: 9069000}},   // Last Petersburg block
				{9069000, ID{Hash: checksumToBytes(0x879d6e30), Next: 9200000}},   // First Istanbul and first Muir Glacier block
				{9199999, ID{Hash: checksumToBytes(0x879d6e30), Next: 9200000}},   // Last Istanbul and first Muir Glacier block
				{9200000, ID{Hash: checksumToBytes(0xe029e991), Next: 12244000}},  // First Muir Glacier block
				{12243999, ID{Hash: checksumToBytes(0xe029e991), Next: 12244000}}, // Last Muir Glacier block
				{12244000, ID{Hash: checksumToBytes(0x0eb440f6), Next: 12965000}}, // First Berlin block
				{12964999, ID{Hash: checksumToBytes(0x0eb440f6), Next: 12965000}}, // Last Berlin block
				{12965000, ID{Hash: checksumToBytes(0xb715077d), Next: 0}},        // First London block
				{20000000, ID{Hash: checksumToBytes(0xb715077d), Next: 0}},        // Future London block
			},
		},
		// Ropsten test cases
		{
			params.RopstenChainConfig,
			params.RopstenGenesisHash,
			[]testcase{
				{4230000, ID{Hash: checksumToBytes(0x97b544f3), Next: 4939394}},   // First Constantinople block
				{4939393, ID{Hash: checksumToBytes(0x97b544f3), Next: 4939394}},   // Last Constantinople block
				{4939394, ID{Hash: checksumToBytes(0xd6e2149b), Next: 6485846}},   // First Petersburg block
				{6485845, ID{Hash: checksumToBytes(0xd6e2149b), Next: 6485846}},   // Last Petersburg block
				{6485846, ID{Hash: checksumToBytes(0x4bc66396), Next: 7117117}},   // First Istanbul block
				{7117116, ID{Hash: checksumToBytes(0x4bc66396), Next: 7117117}},   // Last Istanbul block
				{7117117, ID{Hash: checksumToBytes(0x6727ef90), Next: 9812189}},   // First Muir Glacier block
				{9812188, ID{Hash: checksumToBytes(0x6727ef90), Next: 9812189}},   // Last Muir Glacier block
				{9812189, ID{Hash: checksumToBytes(0xa157d377), Next: 10499401}},  // First Berlin block
				{10499400, ID{Hash: checksumToBytes(0xa157d377), Next: 10499401}}, // Last Berlin block
				{10499401, ID{Hash: checksumToBytes(0x7119b6b3), Next: 0}},        // First London block
				{11000000, ID{Hash: checksumToBytes(0x7119b6b3), Next: 0}},        // Future London block
			},
		},
		// Rinkeby test cases
		{
			params.RinkebyChainConfig,
			params.RinkebyGenesisHash,
			[]testcase{
				{3660663, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // First Constantinople block
				{4321233, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // Last Constantinople block
				{4321234, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // First Petersburg block
				{5435344, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // Last Petersburg block
				{5435345, ID{Hash: checksumToBytes(0xcbdb8838), Next: 8290928}}, // First Istanbul block
				{8290927, ID{Hash: checksumToBytes(0xcbdb8838), Next: 8290928}}, // Last Istanbul block
				{8290928, ID{Hash: checksumToBytes(0x6910c8bd), Next: 8897988}}, // First Berlin block
				{8897987, ID{Hash: checksumToBytes(0x6910c8bd), Next: 8897988}}, // Last Berlin block
				{8897988, ID{Hash: checksumToBytes(0x8E29F2F3), Next: 0}},       // First London block
				{10000000, ID{Hash: checksumToBytes(0x8E29F2F3), Next: 0}},      // Future London block
			},
		},
		// Goerli test cases
		{
			params.GoerliChainConfig,
			params.GoerliGenesisHash,
			[]testcase{
				{1561650, ID{Hash: checksumToBytes(0xa3f5ab08), Next: 1561651}}, // Last Petersburg block
				{1561651, ID{Hash: checksumToBytes(0xc25efa5c), Next: 4460644}}, // First Istanbul block
				{4460643, ID{Hash: checksumToBytes(0xc25efa5c), Next: 4460644}}, // Last Istanbul block
				{4460644, ID{Hash: checksumToBytes(0x757a1c47), Next: 5062605}}, // First Berlin block
				{5000000, ID{Hash: checksumToBytes(0x757a1c47), Next: 5062605}}, // Last Berlin block
				{5062605, ID{Hash: checksumToBytes(0xB8C6299D), Next: 0}},       // First London block
				{6000000, ID{Hash: checksumToBytes(0xB8C6299D), Next: 0}},       // Future London block
			},
		},
	}
	for i, tt := range tests {
		for j, ttt := range tt.cases {
			if have := NewID(tt.config, tt.genesis, ttt.head); have != ttt.want {
				t.Errorf("test %d, case %d: fork ID mismatch: have %x, want %x", i, j, have, ttt.want)
			}
		}
	}
}

// TestValidation tests that a local peer correctly validates and accepts a remote
// fork ID.
func TestValidation(t *testing.T) {
	tests := []struct {
		head uint64
		id   ID
		err  error
	}{
		// Local is mainnet Petersburg, remote announces the same. No future fork is announced.
		{7987396, ID{Hash: checksumToBytes(0x668db0af), Next: 0}, nil},

		// Local is mainnet Petersburg, remote announces the same. Remote also announces a next fork
		// at block 0xffffffff, but that is uncertain.
		{7987396, ID{Hash: checksumToBytes(0x668db0af), Next: math.MaxUint64}, nil},

		// Local is mainnet Petersburg, and isn't aware of more forks. Remote announces Petersburg +
		// 0xffffffff. Local needs software update, reject.
		{7987396, ID{Hash: checksumToBytes(0x5cddc0e1), Next: 0}, ErrLocalIncompatibleOrStale},

		// Local is mainnet Petersburg, remote is Rinkeby Petersburg.
		{7987396, ID{Hash: checksumToBytes(0xafec6b27), Next: 0}, ErrLocalIncompatibleOrStale},

		// Local is mainnet London, far in the future. Remote announces Gopherium (non existing fork)
		// at some future block 88888888, for itself, but past block for local. Local is incompatible.
		//
		// This case detects non-upgraded nodes with majority hash power (typical Ropsten mess).
		{88888888, ID{Hash: checksumToBytes(0xb715077d), Next: 88888888}, ErrLocalIncompatibleOrStale},
	}
	for i, tt := range tests {
		filter := newFilter(params.MainnetChainConfig, params.MainnetGenesisHash, func() uint64 { return tt.head })
		if err := filter(tt.id); err != tt.err {
			t.Errorf("test %d: validation error mismatch: have %v, want %v", i, err, tt.err)
		}
	}
}

// Tests that IDs are properly RLP encoded (specifically important because we
// use uint32 to store the hash, but we need to encode it as [4]byte).
func TestEncoding(t *testing.T) {
	tests := []struct {
		id   ID
		want []byte
	}{
		{ID{Hash: checksumToBytes(0), Next: 0}, common.Hex2Bytes("c6840000000080")},
		{ID{Hash: checksumToBytes(0xdeadbeef), Next: 0xBADDCAFE}, common.Hex2Bytes("ca84deadbeef84baddcafe,")},
		{ID{Hash: checksumToBytes(math.MaxUint32), Next: math.MaxUint64}, common.Hex2Bytes("ce84ffffffff88ffffffffffffffff")},
	}
	for i, tt := range tests {
		have, err := rlp.EncodeToBytes(tt.id)
		if err != nil {
			t.Errorf("test %d: failed to encode forkid: %v", i, err)
			continue
		}
		if !bytes.Equal(have, tt.want) {
			t.Errorf("test %d: RLP mismatch: have %x, want %x", i, have, tt.want)
		}
	}
}
