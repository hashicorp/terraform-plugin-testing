package resource

import tfjson "github.com/hashicorp/terraform-json"

// TODO: document
type PlanAssert interface {
	RunAssert(*tfjson.Plan) error
}
