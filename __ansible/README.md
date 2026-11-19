# Ansible Collection: devops.vmman4

Ansible collection and modules wrapping the **`vmman4`** virtual machine management CLI.

---

## Modules Included

* **`vmman`**: Manage virtual machine lifecycles (create from spec, provision from template, start, stop, reset, adjust vCPUs/memory, and remove).
* **`vmman_info`**: Read-only module to gather structured facts about a virtual machine from `vmman vm info`.

---

## Directory Structure

```text
__ansible/
├── Makefile                 # Build, clean, lint, and test targets
├── galaxy.yaml               # Ansible Galaxy collection metadata
├── plugins/
│   └── modules/
│       ├── vmman.py         # Main VM lifecycle module
│       └── vmman_info.py    # Facts/info gathering module
├── tests/
│   └── test_local.yaml       # Local test playbook
└── README.md                # Documentation and guide
```

---

## Testing Methods

### 1. Direct CLI Invocation (Fast Debugging)
Ansible modules can be executed directly as standalone Python scripts without needing an Ansible playbook. Pipe JSON arguments into the script:

```bash
python3 plugins/modules/vmman.py <<EOF
{
  "ANSIBLE_MODULE_ARGS": {
    "name": "my-vm",
    "state": "running",
    "executable": "/usr/bin/vmman"
  }
}
EOF
```

### 2. Local Playbook via `ANSIBLE_LIBRARY`
You can test the modules directly from this directory without installing the collection by setting `ANSIBLE_LIBRARY`:

```bash
ANSIBLE_LIBRARY=plugins/modules ansible-playbook tests/test_local.yaml
```

Or run a syntax check:
```bash
make test
```

### 3. Collection Build & Local Install
Build the `.tar.gz` collection package and install it locally:

```bash
# Build devops-vmman4-*.tar.gz
make build

# Install into ~/.ansible/collections
make install
```

Once installed, use the fully-qualified collection name (`devops.vmman4.vmman`) in your playbooks:

```yaml
- hosts: localhost
  tasks:
    - name: Ensure test VM is running
      devops.vmman4.vmman:
        name: test-vm
        state: running
```

---

## Distribution & Dependency Configuration (`requirements.yaml`)

Depending on your deployment environment, configure your Ansible playbook repository (`requirements.yaml`) using one of the following methods:

### Method A: Nexus Galaxy Repository *(Official / Target Method)*
> **Note**: Support for an Ansible Galaxy repository hosted in Nexus will eventually be the official distribution repository for this collection.
> Official Repository URL: **`https://nexus.famillegratton.net:1808/repository/ansible/`**

1. In your `ansible.cfg` (e.g. `/localrepos/devops/ansible/ansible.cfg`), define the Nexus Galaxy server:
   ```ini
   [galaxy]
   server_list = nexus_repo

   [galaxy_server.nexus_repo]
   url = https://nexus.famillegratton.net:1808/repository/ansible/
   # If authentication is required:
   # token = <your-api-token>
   ```

2. In your playbook's `requirements.yaml`:
   ```yaml
   collections:
     - name: devops.vmman4
       version: ">=1.0.0"
   ```

3. Install:
   ```bash
   ansible-galaxy collection install -r requirements.yaml
   ```

---

### Method B: Direct Git Repository
Install directly from the Git repository:

```yaml
collections:
  - name: https://git.famillegratton.net:3000/devops/vmman4.git#/__ansible
    type: git
    version: main
```

---

### Method C: Local Tarball
Install from a locally built collection archive:

```yaml
collections:
  - name: /localrepos/devops/vmman4/__ansible/devops-vmman4-1.2.5.tar.gz
    type: file
```

---

## Module Parameters (`devops.vmman4.vmman`)

| Parameter | Type | Default | Description |
|---|---|---|---|
| `name` | string | *required* | Domain name of the virtual machine. |
| `state` | string | `present` | Desired state: `present`, `absent`, `running`, `started`, `stopped`, `restarted`. |
| `template` | string | `None` | Template VM name for cloning (used with `state: present`). |
| `ip_address` | string | `None` | Static IP for the provisioned VM (used with `template`). |
| `env_file` | path | `None` | Path to custom `env.json` passed to `vmman vm provision -E`. |
| `spec_file` | path | `None` | Path to JSON spec file passed to `vmman vm create`. |
| `keep_storage` | bool | `false` | Preserve disk storage when deleting a VM (`vm rm -k`). |
| `memory_mb` | int | `None` | VM memory in MiB (`vm setvmem`). |
| `vcpus` | int | `None` | vCPU count (`vm setvcpus`). |
| `connection_uri`| string | `None` | Hypervisor connection URI passed to `-C`. |
| `connection_file`| path | `None` | Hypervisor connection configuration file passed to `-c`. |
| `executable` | string | `vmman` | Path or binary name for `vmman`. Checks `$PATH` first, falling back to `/opt/bin/vmman`. |
