package beaconapi

import (
	"context"
	"fmt"
	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/api"
	apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/deneb"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	lru "github.com/hashicorp/golang-lru"
	"github.com/rs/zerolog"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

const (
	SLOTS_PER_EPOCH  = "SLOTS_PER_EPOCH"
	SECONDS_PER_SLOT = "SECONDS_PER_SLOT"
)

var (
	validatorListCacheKey = "validator_list"
)

type BeaconClient struct {
	endpoint string
	config   map[string]string
	service  eth2client.Service
	cache    *lru.Cache
}

func NewBeaconGwClient(endpoint string) *BeaconClient {
	cache, _ := lru.New(100)
	return &BeaconClient{
		endpoint: endpoint,
		config:   make(map[string]string),
		cache:    cache,
	}
}

func (b *BeaconClient) GetIntConfig(key string) (int, error) {
	config := b.GetBeaconConfig()
	if v, exist := config[key]; !exist {
		return 0, nil
	} else {
		return strconv.Atoi(v)
	}
}

func (b *BeaconClient) GetBeaconConfig() map[string]string {
	if len(b.config) == 0 {
		config, err := b.GetSpec()
		if err != nil {
			log.WithError(err).Error("get beacon spec failed")
			return nil
		}
		b.config = make(map[string]string)
		for key, v := range config {
			switch v.(type) {
			case time.Duration:
				b.config[key] = strconv.FormatFloat(v.(time.Duration).Seconds(), 'f', -1, 64)
			case time.Time:
				b.config[key] = strconv.FormatInt(v.(time.Time).Unix(), 10)
			case []uint8:
				b.config[key] = fmt.Sprintf("0x%#x", v.([]uint8))
			case int:
				b.config[key] = strconv.Itoa(v.(int))
			case uint64:
				b.config[key] = strconv.FormatUint(v.(uint64), 10)
			case int64:
				b.config[key] = strconv.FormatInt(v.(int64), 10)
			case float64:
				b.config[key] = strconv.FormatFloat(v.(float64), 'f', -1, 64)
			case string:
				b.config[key] = v.(string)
			case phase0.Version:
				b.config[key] = fmt.Sprintf("%#x", v.(phase0.Version))
			case phase0.DomainType:
				b.config[key] = fmt.Sprintf("%#x", v.(phase0.DomainType))
			default:
				log.Warnf("unknown beacon config key %s type %T", key, v)
			}
		}
	}
	return b.config
}

func (b *BeaconClient) getLatestBeaconHeader() (*apiv1.BeaconBlockHeader, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconBlockHeadersProvider).BeaconBlockHeader(context.Background(), &api.BeaconBlockHeaderOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: "head",
	})
	if err != nil {
		log.WithError(err).Error("get latest beacon header failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetValidatorsList() ([]*phase0.Validator, error) {
	if v, ok := b.cache.Get(validatorListCacheKey); ok {
		return v.([]*phase0.Validator), nil
	}
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconStateProvider).BeaconState(context.Background(), &api.BeaconStateOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		State: "head",
	})
	if err != nil {
		log.WithError(err).Error("get beacon state failed")
		return nil, err
	}
	vals, err := res.Data.Validators()
	if err != nil {
		log.WithError(err).Error("get validators failed")
		return nil, err
	}
	b.cache.Add(validatorListCacheKey, vals)

	return vals, nil
}

func (b *BeaconClient) GetLatestValidators() (*spec.VersionedBeaconState, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconStateProvider).BeaconState(context.Background(), &api.BeaconStateOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		State: "head",
	})
	if err != nil {
		log.WithError(err).Error("get beacon state failed")
		return nil, err
	}

	return res.Data, nil
}

// GetBeaconState
// slot: "head", "genesis", "finalized", "justified", <slot>, <hex encoded stateRoot with 0x prefix>.
func (b *BeaconClient) GetBeaconState(slot string) (*spec.VersionedBeaconState, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconStateProvider).BeaconState(context.Background(), &api.BeaconStateOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		State: slot,
	})
	if err != nil {
		log.WithError(err).Error("get beacon state failed")
		return nil, err
	}

	return res.Data, nil
}

func (b *BeaconClient) GetLatestBeaconHeader() (*apiv1.BeaconBlockHeader, error) {
	return b.getLatestBeaconHeader()
}

func (b *BeaconClient) getAllValReward(epoch int) (*apiv1.AttestationRewards, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.AttestationRewardsProvider).AttestationRewards(context.Background(), &api.AttestationRewardsOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 20,
		},
		Epoch: phase0.Epoch(epoch),
	})
	if err != nil {
		log.WithField("epoch", epoch).WithError(err).Error("get val reward failed")
		return nil, err
	}

	return res.Data, nil
}

func (b *BeaconClient) GetAllValReward(epoch int) (*apiv1.AttestationRewards, error) {
	info, err := b.getAllValReward(epoch)
	if err != nil {
		return nil, err
	}
	return info, err
}

func (b *BeaconClient) getProposerDuties(epoch int) ([]*apiv1.ProposerDuty, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.ProposerDutiesProvider).ProposerDuties(context.Background(), &api.ProposerDutiesOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Epoch: phase0.Epoch(epoch),
	})
	if err != nil {
		log.WithError(err).Error("get proposer duties failed")
		return nil, err
	}

	return res.Data, nil
}

// /eth/v1/validator/duties/proposer/:epoch
func (b *BeaconClient) GetProposerDuties(epoch int) ([]*apiv1.ProposerDuty, error) {
	return b.getProposerDuties(epoch)
}

func (b *BeaconClient) getAttesterDuties(epoch int, vals []int) ([]*apiv1.AttesterDuty, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	indices := make([]phase0.ValidatorIndex, len(vals))
	for _, val := range vals {
		indices = append(indices, phase0.ValidatorIndex(val))
	}
	if len(indices) == 0 {
		// get validators list
		valList, err := b.GetValidatorsList()
		if err != nil {
			log.WithError(err).Error("get validators failed")
			return nil, err
		}
		indices = make([]phase0.ValidatorIndex, len(valList))
		for i, _ := range valList {
			indices[i] = phase0.ValidatorIndex(i)
		}
	}
	res, err := service.(eth2client.AttesterDutiesProvider).AttesterDuties(context.Background(), &api.AttesterDutiesOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 20,
		},
		Epoch:   phase0.Epoch(epoch),
		Indices: indices,
	})
	if err != nil {
		log.WithError(err).Error("get attester duties failed")
		return nil, err
	}

	return res.Data, nil
}

func (b *BeaconClient) FetchBlockAttestation(slot int64) ([]*phase0.Attestation, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), &api.SignedBeaconBlockOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("get block attestation failed")
		return nil, err
	}
	blk := res.Data.Deneb
	return blk.Message.Body.Attestations, nil
}

func (b *BeaconClient) FetchBlocksAttestations(slots []int64) ([]*phase0.Attestation, error) {
	attestations := make([]*phase0.Attestation, 0)
	for _, slot := range slots {
		att, err := b.FetchBlockAttestation(slot)
		if err != nil {
			log.WithError(err).Errorf("fetch block attestation for slot %d failed", slot)
			continue
		}
		attestations = append(attestations, att...)
	}
	return attestations, nil
}

// POST /eth/v1/validator/duties/attester/:epoch
func (b *BeaconClient) GetAttesterDuties(epoch int, vals []int) ([]*apiv1.AttesterDuty, error) {
	return b.getAttesterDuties(epoch, vals)
}

func (b *BeaconClient) GetEpochProposerDuties(epoch int) ([]*apiv1.ProposerDuty, error) {
	return b.GetProposerDuties(epoch)
}

func (b *BeaconClient) getBlockReward(slot int) (*apiv1.BlockRewards, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BlockRewardsProvider).BlockRewards(context.Background(), &api.BlockRewardsOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithField("slot", slot).WithError(err).Error("get block reward failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetBlockReward(slot int) (*apiv1.BlockRewards, error) {
	return b.getBlockReward(slot)
}

func (b *BeaconClient) getSlotRoot(slot int64) (*phase0.Root, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.BeaconBlockRootProvider).BeaconBlockRoot(context.Background(), &api.BeaconBlockRootOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("getSlotRoot failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetSlotRoot(slot int64) (string, error) {
	root, err := b.getSlotRoot(slot)
	if err != nil {
		return "", err
	}
	if root == nil {
		return "0x", nil
	}
	return root.String(), nil
}

func (b *BeaconClient) getService() (eth2client.Service, error) {
	if b.service == nil {
		service, err := newClient(context.Background(), b.endpoint)
		if err != nil {
			log.WithField("endpoint", b.endpoint).WithError(err).Error("create eth2client failed")
			return nil, err
		}
		b.service = service
	}
	return b.service, nil
}

func (b *BeaconClient) GetBlockById(id string) (*spec.VersionedSignedBeaconBlock, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	opts := &api.SignedBeaconBlockOpts{
		Block: id,
	}
	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), opts)
	if err != nil {
		log.WithError(err).Error("get block failed")
		return &spec.VersionedSignedBeaconBlock{}, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetBlockHeaderById(id string) (*apiv1.BeaconBlockHeader, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	opts := &api.BeaconBlockHeaderOpts{
		Block: id,
	}
	res, err := service.(eth2client.BeaconBlockHeadersProvider).BeaconBlockHeader(context.Background(), opts)
	if err != nil {
		log.WithError(err).Error("get block header failed")
		return &apiv1.BeaconBlockHeader{}, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetDenebBlockBySlot(slot uint64) (*deneb.SignedBeaconBlock, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), &api.SignedBeaconBlockOpts{
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("get block failed")
		return nil, err
	}
	return res.Data.Deneb, nil
}

func (b *BeaconClient) GetCapellaBlockBySlot(slot uint64) (*capella.SignedBeaconBlock, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}

	res, err := service.(eth2client.SignedBeaconBlockProvider).SignedBeaconBlock(context.Background(), &api.SignedBeaconBlockOpts{
		Block: fmt.Sprintf("%d", slot),
	})
	if err != nil {
		log.WithError(err).Error("get block failed")
		return nil, err
	}
	return res.Data.Capella, nil
}

func (b *BeaconClient) GetSpec() (map[string]any, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.SpecProvider).Spec(context.Background(), &api.SpecOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
	})
	if err != nil {
		log.WithError(err).Error("get genesis failed")
		return nil, err
	}
	return res.Data, nil
}

func (b *BeaconClient) GetGenesis() (*apiv1.Genesis, error) {
	service, err := b.getService()
	if err != nil {
		log.WithError(err).Error("create eth2client failed")
		return nil, err
	}
	res, err := service.(eth2client.GenesisProvider).Genesis(context.Background(), &api.GenesisOpts{
		Common: api.CommonOpts{
			Timeout: time.Second * 10,
		},
	})
	if err != nil {
		log.WithError(err).Error("get genesis failed")
		return nil, err
	}
	return res.Data, nil
}

func newClient(ctx context.Context, endpoint string) (eth2client.Service, error) {
	return http.New(ctx,
		// WithAddress supplies the address of the beacon node, as a URL.
		http.WithAddress(endpoint),
		// LogLevel supplies the level of logging to carry out.
		http.WithLogLevel(zerolog.WarnLevel),
	)
}
