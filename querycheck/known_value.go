package querycheck

import (
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type KnownValueCheck struct {
	Path       tfjsonpath.Path
	KnownValue knownvalue.Check
}
