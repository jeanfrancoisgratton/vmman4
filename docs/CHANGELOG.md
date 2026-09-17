| Release  | Date               | Comments             |
|----------|--------------------|----------------------|
| 1.2.0 | 2026.09.17 | Stop/Start VM(s) now pause 3 seconds between operations |
| 1.1.0 | 2026.09.14 | New `cluster` command group: `ls`, `define`, `remove`, `start`, `stop`, `reset`, plus `cluster CLUSTER_NAME snaprev [SNAPSHOT_NAME]` to revert every node in a cluster to a snapshot<br>New `vol attach VM POOL NAME` command: attaches an existing volume to a VM as a new virtio disk, shutting the VM down first if it's running<br>Filled in missing `--help`/`Example` text across `src/cmd/` (notably the `vm` and `conn` command groups, which had none)<br>Fixed stale/copy-pasted flag and command descriptions in the `conn` group |
| 1.0.0 | 2026.09.14 | First stable release<br>Installed binary renamed from `vmman4` to `vmman` (package name stays `vmman4` on every format)<br>New `vol`/`volume` command group: `create`, `rm`, `list`/`ls`<br>`pool` command group extended with `create`, `start`, `stop`, `rm` (previously `list` only) |
| 0.30.0 | 2026.09.09 | Version numbering now SemVer-aligned<br>GO version bump: 1.27.1<br>Windows build support |



