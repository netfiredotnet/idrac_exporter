# NetFire build of idrac_exporter

`master` mirrors upstream (mrlhansen/idrac_exporter). The `netfire` branch is
what nf-mon-stack ships: an upstream release tag plus the patches below,
released as `v<upstream>-nf.<n>` (for example `v2.6.3-nf.1`). Creating a
GitHub release for such a tag builds the binaries the agent pins in
`agent/release/dependencies.json`.

Patches on top of v2.6.3, each meant to go upstream:

- Export `chassis_health` and `chassis_intrusion` from the primary chassis.
- Accept volume sizes (`StripSizeBytes` and friends) sent as strings.

Only the Go binary workflow runs here; the Docker image and Helm chart
workflows are removed so a tag does not publish images under this org.
