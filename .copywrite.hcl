schema_version = 1

project {
  license        = "MPL-2.0"
  copyright_year = 2014

  header_ignore = [
    "META.d/**/*.yaml",

    # changie tooling configuration and CHANGELOG entries (prose)
    ".changes/unreleased/**",
    ".changie.yaml",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # GitHub Actions workflow-specific configurations
    ".github/labeler-*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",

    # Release Engineering tooling configuration
    ".release/*.hcl",
  ]
}
