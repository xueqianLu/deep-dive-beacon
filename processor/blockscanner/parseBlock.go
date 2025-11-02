package beaconscanner

import (
	"encoding/hex"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/electra"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/xueqianLu/deep-dive-beacon/models/dbmodels"
)

func (s *BeaconBlockScanner) phase0ToDBBlock(blk *phase0.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.Signature = blk.Signature.String()
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()

	return dbinfo, nil
}

func (s *BeaconBlockScanner) altairToDBBlock(blk *altair.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}

func (s *BeaconBlockScanner) bellatrixToDBBlock(blk *bellatrix.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}

func (s *BeaconBlockScanner) capellaToDBBlock(blk *capella.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}

func (s *BeaconBlockScanner) denebToDBBlock(blk *deneb.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}

func (s *BeaconBlockScanner) electraToDBBlock(blk *electra.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}

func (s *BeaconBlockScanner) fuluToDBBlock(blk *electra.SignedBeaconBlock, dbinfo *dbmodels.BeaconBlock) (*dbmodels.BeaconBlock, error) {
	dbinfo.StateRoot = blk.Message.StateRoot.String()
	dbinfo.ParentRoot = blk.Message.ParentRoot.String()
	dbinfo.ProposerIndex = uint64(blk.Message.ProposerIndex)
	dbinfo.Eth1BlockHash = hex.EncodeToString(blk.Message.Body.ETH1Data.BlockHash)
	dbinfo.Eth1DepositCount = uint64(blk.Message.Body.ETH1Data.DepositCount)
	dbinfo.Eth1DepositRoot = blk.Message.Body.ETH1Data.DepositRoot.String()
	dbinfo.Graffiti = hex.EncodeToString(blk.Message.Body.Graffiti[:])
	dbinfo.RandaoReveal = blk.Message.Body.RANDAOReveal.String()
	dbinfo.ProposerSlashed = uint(len(blk.Message.Body.ProposerSlashings))
	dbinfo.AttesterSlashed = uint(len(blk.Message.Body.AttesterSlashings))
	dbinfo.Signature = blk.Signature.String()
	return dbinfo, nil
}
