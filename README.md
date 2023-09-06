[![PkgGoDev](https://pkg.go.dev/badge/github.com/hashicorp/terraform-plugin-log)](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-testing)

# terraform-plugin-testing

terraform-plugin-testing is a helper module for testing Terraform providers. Terraform acceptance tests use real Terraform configurations to exercise the code in real plan, apply, refresh, and destroy life cycles. 
When run from the root of a Terraform Provider codebase, Terraform’s testing framework compiles the current provider in-memory and executes the provided configuration in developer defined steps, creating infrastructure along the way.

## Go Compatibility

This project follows the [support policy](https://golang.org/doc/devel/release.html#policy) of Go as its support policy. The two latest major releases of Go are supported by the project.

Currently, that means Go **1.20** or later must be used when including this project as a dependency.

## Contributing

See [`.github/CONTRIBUTING.md`](https://github.com/hashicorp/terraform-plugin-testing/blob/main/.github/CONTRIBUTING.md)

## License

[Mozilla Public License v2.0](https://github.com/hashicorp/terraform-plugin-testing/blob/main/LICENSE)
