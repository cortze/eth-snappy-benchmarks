package cmd

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"time"

	"github.com/attestantio/go-eth2-client/http"
	"github.com/attestantio/go-eth2-client/spec"
	"github.com/cortze/eth-snappy-benchmarks/csvs"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

const (
	jsonEncoding = "json"
	sszEncoding  = "ssz"

	queryTimeout = 10 * time.Second
)

var (
	ethNodeEndp      string = ""
	blockEncoding    string = jsonEncoding
	targetBlocksFile string = ""
	outputFolder     string = "raw-blocks"
)

var FetchBlocksCmd = &cli.Command{
	Name:        "fetch-blocks",
	Description: "",
	Action:      RunBlockFetcher,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "eth-node",
			DefaultText: "endpoint of the node that will provide the blocks",
			Required:    true,
			Destination: &ethNodeEndp,
		},
		&cli.StringFlag{
			Name:        "encoding",
			DefaultText: "type of encoding for the block download",
			Required:    true,
			Destination: &blockEncoding,
		},
		&cli.StringFlag{
			Name:        "block-list",
			DefaultText: "folder placing all the blocks to test",
			Required:    true,
			Destination: &targetBlocksFile,
		},
		&cli.StringFlag{
			Name:        "output-folder",
			DefaultText: "folder where the raw-blocks will be stored",
			Required:    true,
			Destination: &outputFolder,
		},
	},
}

func RunBlockFetcher(ctx *cli.Context) error {

	// get connection to beacon node
	beaconNode, err := http.New(
		ctx.Context,
		http.WithAddress(ethNodeEndp),
		http.WithLogLevel(zerolog.WarnLevel),
		http.WithTimeout(queryTimeout),
	)

	// read the number of blocks from the given csv file
	csvFile, err := csvs.NewCsvImporter(targetBlocksFile)
	if err != nil {
		return errors.Wrap(err, "unable to open csv file with target blocks")
	}
	blockList, err := csvFile.ReadTargetBlocks()
	if err != nil {
		return errors.Wrap(err, "unable to parse csv file to target blocks")
	}

	log.Info("Client ready, requesting blocks")
	for _, block := range blockList {
		// fetch the block
		signedBlock, err := fetchBlock(beaconNode.(*http.Service), ctx.Context, block.Slot)
		if err != nil {
			log.Error(errors.Wrap(err, "unable to retrieve block "+block.Slot))
			continue
		}
		if signedBlock == nil {
			log.Warnf("NIL block downloaded, was block %s proposed?", block.Slot)
			continue
		}

		var blockBytes []byte
		switch blockEncoding {
		case jsonEncoding:
			blockBytes, err = encodeToJson(signedBlock)
			if err != nil {
				log.Error(errors.Wrap(err, "unable to convert to JSON block "+block.Slot))
			}
		case sszEncoding:
			blockBytes, err = encodeToSSZ(signedBlock)
			if err != nil {
				log.Error(errors.Wrap(err, "unable to convert to SSZ block "+block.Slot))
			}
		default:
			log.Error("unable to recognize encoding og blocks " + blockEncoding)
		}
		err = persitBlockToDisk("block_"+block.Slot+"."+blockEncoding, outputFolder, blockBytes)
		if err != nil {
			log.Error(errors.Wrap(err, "unable to persist to disk block "+block.Slot))
		}
	}
	return nil
}

func fetchBlock(bn *http.Service, ctx context.Context, blockNumber string) (*spec.VersionedSignedBeaconBlock, error) {
	versionedBlock, err := bn.SignedBeaconBlock(ctx, blockNumber)
	if err != nil {
		// close the channel (to tell other routines to stop processing and end)
		return nil, errors.Wrap(err, "unable to retrieve Beacon Block at slot "+blockNumber)
	}
	log.Info("successfully got block " + blockNumber)
	return versionedBlock, err
}

func encodeToJson(block *spec.VersionedSignedBeaconBlock) (blockBytes []byte, err error) {
	// check with version are they using
	switch block.Version {
	case spec.DataVersionPhase0:
		blockBytes, err = block.Phase0.MarshalJSON()
	case spec.DataVersionAltair:
		blockBytes, err = block.Altair.MarshalJSON()
	case spec.DataVersionBellatrix:
		blockBytes, err = block.Bellatrix.MarshalJSON()
	case spec.DataVersionCapella:
		blockBytes, err = block.Capella.MarshalJSON()
	case spec.DataVersionDeneb:
		blockBytes, err = block.Deneb.MarshalJSON()
	default:
		err = errors.New("unable to get JSON data from block")
	}
	return
}

func encodeToSSZ(block *spec.VersionedSignedBeaconBlock) (blockBytes []byte, err error) {
	// check with version are they using
	switch block.Version {
	case spec.DataVersionPhase0:
		blockBytes, err = block.Phase0.MarshalSSZ()
	case spec.DataVersionAltair:
		blockBytes, err = block.Altair.MarshalSSZ()
	case spec.DataVersionBellatrix:
		blockBytes, err = block.Bellatrix.MarshalSSZ()
	case spec.DataVersionCapella:
		blockBytes, err = block.Capella.MarshalSSZ()
	case spec.DataVersionDeneb:
		blockBytes, err = block.Deneb.MarshalSSZ()
	default:
		err = errors.New("unable to get JSON data from block")
	}
	return
}

func persitBlockToDisk(blockName, outputFolder string, blockBytes []byte) error {
	path := outputFolder + "/" + blockName
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(blockBytes)
	if err != nil {
		return err
	}
	return f.Close()
}
