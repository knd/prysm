package spectest

import (
	"io/ioutil"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/blocks"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params/spectest"
	"gopkg.in/d4l3k/messagediff.v1"
)

func runAttesterSlashingTest(t *testing.T, filename string) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("Could not load file %v", err)
	}

	test := &BlockOperationTest{}
	if err := yaml.Unmarshal(file, test); err != nil {
		t.Fatalf("Failed to Unmarshal: %v", err)
	}

	if err := spectest.SetConfig(test.Config); err != nil {
		t.Fatal(err)
	}

	for _, tt := range test.TestCases {
		t.Run(tt.Description, func(t *testing.T) {
			helpers.ClearAllCaches()

			body := &pb.BeaconBlockBody{AttesterSlashings: []*pb.AttesterSlashing{tt.AttesterSlashing}}

			postState, err := blocks.ProcessAttesterSlashings(tt.Pre, body, true)
			// Note: This doesn't test anything worthwhile. It essentially tests
			// that *any* error has occurred, not any specific error.
			if tt.Post == nil {
				if err == nil {
					t.Fatal("Did not fail when expected")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !proto.Equal(postState, tt.Post) {
				diff, _ := messagediff.PrettyDiff(postState, tt.Post)
				t.Log(diff)
				t.Fatal("Post state does not match expected")
			}
		})
	}
}

var attesterSlashingPrefix = "eth2_spec_tests/tests/operations/attester_slashing/"

func TestAttesterSlashingMinimal(t *testing.T) {
	filepath, err := bazel.Runfile(attesterSlashingPrefix + "attester_slashing_minimal.yaml")
	if err != nil {
		t.Fatal(err)
	}
	runAttesterSlashingTest(t, filepath)
}

func TestAttesterSlashingMainnet(t *testing.T) {
	filepath, err := bazel.Runfile(attesterSlashingPrefix + "attester_slashing_mainnet.yaml")
	if err != nil {
		t.Fatal(err)
	}
	runAttesterSlashingTest(t, filepath)
}
