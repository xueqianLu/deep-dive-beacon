package directlysync

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

func (s *DirectlyBlockScanner) getPhase0Attestations(blk *phase0.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getAltairAttestations(blk *altair.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getBellatrixAttestations(blk *bellatrix.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getCapellaAttestations(blk *capella.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getDenebAttestations(blk *deneb.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getElectraAttestations(blk *electra.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}

func (s *DirectlyBlockScanner) getFuluAttestations(blk *electra.SignedBeaconBlock) []*dbmodels.BeaconAttestation {
	slot := blk.Message.Slot
	var res = make([]*dbmodels.BeaconAttestation, 0)
	for i, att := range blk.Message.Body.Attestations {
		dbAtt := &dbmodels.BeaconAttestation{
			SlotNumber:      uint64(slot),
			AttestIndex:     i,
			AggregationBits: hex.EncodeToString(att.AggregationBits.Bytes()),
			BeaconBlockRoot: att.Data.BeaconBlockRoot.String(),
			CommitteeIndex:  uint64(att.Data.Index),
			SourceEpoch:     uint64(att.Data.Source.Epoch),
			SourceRoot:      att.Data.Source.Root.String(),
			TargetEpoch:     uint64(att.Data.Target.Epoch),
			TargetRoot:      att.Data.Target.Root.String(),
			Signature:       att.Signature.String(),
		}
		res = append(res, dbAtt)
	}
	return res
}
