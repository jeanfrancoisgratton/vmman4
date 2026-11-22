#!/usr/bin/python3
# -*- coding: utf-8 -*-

# Copyright: (c) 2026, J.F. Gratton <jean-francois@famillegratton.net>
# GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)

from __future__ import absolute_import, division, print_function
__metaclass__ = type

DOCUMENTATION = r'''
---
module: vmman
short_description: Manage virtual machines using vmman4
description:
  - Create, provision, start, stop, restart, configure, and remove virtual machines via the vmman4 CLI.
version_added: "1.0.0"
author:
  - J.F. Gratton (@jeanfrancoisgratton)
options:
  name:
    description:
      - Name of the virtual machine.
    type: str
    required: true
  state:
    description:
      - Desired state of the virtual machine.
      - C(present) ensures the VM exists (creates or provisions if absent, updates memory/vCPUs if changed).
      - C(absent) ensures the VM is removed and undefined.
      - C(running) / C(started) ensures the VM is running.
      - C(stopped) ensures the VM is shut down.
      - C(restarted) resets/reboots the VM.
    type: str
    choices: [present, absent, running, started, stopped, restarted]
    default: present
  template:
    description:
      - Name of template VM to clone when provisioning a new VM with C(state=present).
    type: str
  ip_address:
    description:
      - IP address to assign to the provisioned VM (requires C(template)).
    type: str
  env_file:
    description:
      - Path to custom environment JSON file for provisioning (passed to C(-E)).
    type: path
  spec_file:
    description:
      - Path to a VMSpec JSON file to create a VM definition from.
    type: path
  keep_storage:
    description:
      - Whether to preserve the VM's disk volume when removing it (passed to C(vm rm -k)).
    type: bool
    default: false
  memory_mb:
    description:
      - Desired memory in MiB for the VM (sets memory via C(vm setvmem)).
    type: int
  vcpus:
    description:
      - Desired number of vCPUs for the VM (sets vCPUs via C(vm setvcpus)).
    type: int
  connection_uri:
    description:
      - Libvirt connection URI (e.g. C(qemu:///system) or C(qemu+ssh://host/system)).
      - Overrides C(VMMAN_CONNECTURI) / default local URI via C(-C).
    type: str
  connection_file:
    description:
      - Path to connection configuration file via C(-c).
    type: path
  executable:
    description:
      - Path or name of the vmman4 binary. Defaults to finding C(vmman) or C(vmman4) in PATH, falling back to C(/opt/bin/vmman).
    type: str
    default: vmman
'''

EXAMPLES = r'''
- name: Ensure a VM is provisioned from a template
  jfgratton.vmman4.vmman:
    name: web-node-01
    state: present
    template: debian12-template
    ip_address: 192.168.1.100

- name: Create a VM from a JSON specification file
  jfgratton.vmman4.vmman:
    name: db-node-01
    state: present
    spec_file: /etc/vmman/specs/db-node.json

- name: Ensure VM is running
  jfgratton.vmman4.vmman:
    name: web-node-01
    state: running

- name: Stop a running VM
  jfgratton.vmman4.vmman:
    name: web-node-01
    state: stopped

- name: Adjust VM resources
  jfgratton.vmman4.vmman:
    name: web-node-01
    memory_mb: 4096
    vcpus: 4

- name: Delete a VM but keep its storage disks
  jfgratton.vmman4.vmman:
    name: web-node-01
    state: absent
    keep_storage: true
'''

RETURN = r'''
vm:
  description: Facts and status of the virtual machine.
  type: dict
  returned: always
  contains:
    name:
      description: VM domain name.
      type: str
      returned: always
    exists:
      description: Whether the VM is currently defined.
      type: bool
      returned: always
    state:
      description: Current operational state (e.g. Running, Shutoff).
      type: str
      returned: when VM exists
    vcpus:
      description: Number of virtual CPUs.
      type: int
      returned: when VM exists
    uuid:
      description: Domain UUID.
      type: str
      returned: when VM exists
    ip_address:
      description: Current IP address.
      type: str
      returned: when available
    storage_pool:
      description: Storage pool where disks reside.
      type: str
      returned: when VM exists
'''

import os
import re
from ansible.module_utils.basic import AnsibleModule

ANSI_REGEX = re.compile(r'\x1b\[[0-9;]*[a-zA-Z]')


def strip_ansi(text):
    if not text:
        return ""
    return ANSI_REGEX.sub('', text)


def locate_executable(module, exe_name):
    # 1st option: find binary in $PATH
    path = module.get_bin_path(exe_name, required=False)
    if path:
        return path
    if exe_name == 'vmman':
        alt = module.get_bin_path('vmman4', required=False)
        if alt:
            return alt

    # Fallback preflight: check /opt/bin if not found in $PATH
    # (Packaging tools and src/build.sh build the binary into /opt/bin)
    path = module.get_bin_path(exe_name, required=False, opt_dirs=['/opt/bin'])
    if path:
        return path
    if exe_name == 'vmman':
        alt = module.get_bin_path('vmman4', required=False, opt_dirs=['/opt/bin'])
        if alt:
            return alt

    # Direct filesystem check for hardcoded /opt/bin fallback or absolute paths
    candidates = []
    if os.path.isabs(exe_name):
        candidates.append(exe_name)
    elif exe_name in ('vmman', 'vmman4'):
        candidates.extend(['/opt/bin/vmman', '/opt/bin/vmman4'])
    else:
        candidates.append(os.path.join('/opt/bin', exe_name))

    for candidate in candidates:
        if os.path.isfile(candidate) and not os.path.isdir(candidate) and os.access(candidate, os.X_OK):
            return candidate

    module.fail_json(msg="Could not find executable '{0}' (or 'vmman4') in PATH or /opt/bin.".format(exe_name))


def build_base_cmd(executable, conn_uri=None, conn_file=None):
    cmd = [executable]
    if conn_uri:
        cmd.extend(['-C', conn_uri])
    if conn_file:
        cmd.extend(['-c', conn_file])
    return cmd


def get_vm_info(module, base_cmd, vm_name):
    cmd = base_cmd + ['vm', 'info', vm_name]
    rc, stdout, stderr = module.run_command(cmd)

    clean_out = strip_ansi(stdout)
    clean_err = strip_ansi(stderr)

    # Check for domain not found
    if "not found" in clean_out.lower() or "not found" in clean_err.lower():
        return {'exists': False, 'name': vm_name}

    if rc != 0:
        # Some errors indicate missing domain
        if "unable to get domain" in clean_err.lower() or "not found" in clean_out.lower():
            return {'exists': False, 'name': vm_name}
        module.fail_json(msg="Failed to query VM '{0}': {1}".format(vm_name, clean_err or clean_out), rc=rc)

    info = {'exists': True, 'name': vm_name, 'disks': []}
    for raw_line in clean_out.splitlines():
        line = raw_line.strip()
        if not line or ':' not in line:
            continue
        parts = line.split(':', 1)
        key = parts[0].strip().lower()
        val = parts[1].strip()

        if key == 'vm name':
            info['name'] = val
        elif key == 'state':
            info['state'] = val
        elif key == 'domain id':
            try:
                info['domain_id'] = int(val)
            except ValueError:
                info['domain_id'] = val
        elif key == 'uuid':
            info['uuid'] = val
        elif key == 'persistent':
            info['persistent'] = (val.lower() == 'true')
        elif key == 'autostart':
            info['autostart'] = (val.lower() == 'true')
        elif 'min memory' in key:
            info['min_memory_gb'] = val.replace('GB', '').strip()
        elif 'max memory' in key:
            info['max_memory_gb'] = val.replace('GB', '').strip()
        elif key == 'vcpus':
            try:
                info['vcpus'] = int(val)
            except ValueError:
                info['vcpus'] = val
        elif key == 'snapshots':
            try:
                info['snapshots'] = int(val)
            except ValueError:
                info['snapshots'] = val
        elif key == 'current snapshot':
            info['current_snapshot'] = val
        elif key == 'storage pool':
            info['storage_pool'] = val
        elif key == 'interface':
            info['interface'] = val if val != 'n/a' else None
        elif key == 'ip address':
            info['ip_address'] = val if val != 'n/a' else None
        elif key.startswith('disk'):
            info['disks'].append(val)

    return info


def run_module():
    module_args = dict(
        name=dict(type='str', required=True),
        state=dict(type='str', choices=['present', 'absent', 'running', 'started', 'stopped', 'restarted'], default='present'),
        template=dict(type='str', required=False),
        ip_address=dict(type='str', required=False),
        env_file=dict(type='path', required=False),
        spec_file=dict(type='path', required=False),
        keep_storage=dict(type='bool', default=False),
        memory_mb=dict(type='int', required=False),
        vcpus=dict(type='int', required=False),
        connection_uri=dict(type='str', required=False),
        connection_file=dict(type='path', required=False),
        executable=dict(type='str', default='vmman'),
    )

    result = dict(
        changed=False,
        vm=dict(),
        msg="",
    )

    module = AnsibleModule(
        argument_spec=module_args,
        supports_check_mode=True
    )

    name = module.params['name']
    state = module.params['state']
    template = module.params['template']
    ip_address = module.params['ip_address']
    env_file = module.params['env_file']
    spec_file = module.params['spec_file']
    keep_storage = module.params['keep_storage']
    memory_mb = module.params['memory_mb']
    vcpus = module.params['vcpus']
    connection_uri = module.params['connection_uri']
    connection_file = module.params['connection_file']
    exe_name = module.params['executable']

    executable = locate_executable(module, exe_name)
    base_cmd = build_base_cmd(executable, connection_uri, connection_file)

    current_info = get_vm_info(module, base_cmd, name)
    exists = current_info.get('exists', False)
    current_state = current_info.get('state', '').lower()

    # --- STATE: ABSENT ---
    if state == 'absent':
        if not exists:
            result['changed'] = False
            result['vm'] = current_info
            result['msg'] = "VM '{0}' does not exist.".format(name)
            module.exit_json(**result)

        if module.check_mode:
            result['changed'] = True
            result['msg'] = "VM '{0}' would be removed.".format(name)
            module.exit_json(**result)

        rm_cmd = base_cmd + ['vm', 'rm', name]
        if keep_storage:
            rm_cmd.append('-k')
        rc, out, err = module.run_command(rm_cmd)
        if rc != 0:
            module.fail_json(msg="Failed to remove VM '{0}': {1}".format(name, err or out), rc=rc)

        result['changed'] = True
        result['vm'] = {'exists': False, 'name': name}
        result['msg'] = "VM '{0}' was removed.".format(name)
        module.exit_json(**result)

    # --- IF VM DOES NOT EXIST YET AND STATE REQUIRES CREATION ---
    if not exists:
        if state in ['running', 'started', 'stopped', 'restarted']:
            module.fail_json(msg="VM '{0}' does not exist; cannot set state to '{1}'.".format(name, state))

        # state == 'present'
        if template and ip_address:
            # Provision VM
            prov_cmd = base_cmd + ['vm', 'provision', name, ip_address, template]
            if env_file:
                prov_cmd.extend(['-E', env_file])

            if module.check_mode:
                result['changed'] = True
                result['msg'] = "VM '{0}' would be provisioned from template '{1}'.".format(name, template)
                module.exit_json(**result)

            rc, out, err = module.run_command(prov_cmd)
            if rc != 0:
                module.fail_json(msg="Failed to provision VM '{0}': {1}".format(name, err or out), rc=rc)

            result['changed'] = True
            result['vm'] = get_vm_info(module, base_cmd, name)
            result['msg'] = "VM '{0}' provisioned successfully.".format(name)
            module.exit_json(**result)

        elif spec_file:
            # Create from JSON spec
            create_cmd = base_cmd + ['vm', 'create', spec_file]
            if module.check_mode:
                result['changed'] = True
                result['msg'] = "VM '{0}' would be created from spec '{1}'.".format(name, spec_file)
                module.exit_json(**result)

            rc, out, err = module.run_command(create_cmd)
            if rc != 0:
                module.fail_json(msg="Failed to create VM '{0}' from spec: {1}".format(name, err or out), rc=rc)

            result['changed'] = True
            result['vm'] = get_vm_info(module, base_cmd, name)
            result['msg'] = "VM '{0}' created from spec successfully.".format(name)
            module.exit_json(**result)
        else:
            module.fail_json(msg="VM '{0}' does not exist. Provide 'template' and 'ip_address' to provision, or 'spec_file' to create.".format(name))

    # --- VM ALREADY EXISTS ---

    # Resource updates (memory, vcpus)
    if memory_mb is not None:
        mem_cmd = base_cmd + ['vm', 'setvmem', name, str(memory_mb)]
        if module.check_mode:
            result['changed'] = True
        else:
            rc, out, err = module.run_command(mem_cmd)
            if rc != 0:
                module.fail_json(msg="Failed to update memory for '{0}': {1}".format(name, err or out), rc=rc)
            result['changed'] = True

    if vcpus is not None:
        current_vcpus = current_info.get('vcpus')
        if current_vcpus != vcpus:
            vcpu_cmd = base_cmd + ['vm', 'setvcpus', name, str(vcpus)]
            if module.check_mode:
                result['changed'] = True
            else:
                rc, out, err = module.run_command(vcpu_cmd)
                if rc != 0:
                    module.fail_json(msg="Failed to set vCPUs for '{0}': {1}".format(name, err or out), rc=rc)
                result['changed'] = True

    # State transitions
    if state in ['running', 'started']:
        if current_state == 'running':
            result['vm'] = current_info
            result['msg'] = "VM '{0}' is already running.".format(name)
            module.exit_json(**result)

        if module.check_mode:
            result['changed'] = True
            result['msg'] = "VM '{0}' would be started.".format(name)
            module.exit_json(**result)

        start_cmd = base_cmd + ['vm', 'start', name]
        rc, out, err = module.run_command(start_cmd)
        if rc != 0:
            module.fail_json(msg="Failed to start VM '{0}': {1}".format(name, err or out), rc=rc)

        result['changed'] = True
        result['vm'] = get_vm_info(module, base_cmd, name)
        result['msg'] = "VM '{0}' started.".format(name)
        module.exit_json(**result)

    elif state == 'stopped':
        if current_state in ['shutoff', 'stopped', 'shut down']:
            result['vm'] = current_info
            result['msg'] = "VM '{0}' is already stopped.".format(name)
            module.exit_json(**result)

        if module.check_mode:
            result['changed'] = True
            result['msg'] = "VM '{0}' would be stopped.".format(name)
            module.exit_json(**result)

        stop_cmd = base_cmd + ['vm', 'stop', name]
        rc, out, err = module.run_command(stop_cmd)
        if rc != 0:
            module.fail_json(msg="Failed to stop VM '{0}': {1}".format(name, err or out), rc=rc)

        result['changed'] = True
        result['vm'] = get_vm_info(module, base_cmd, name)
        result['msg'] = "VM '{0}' stopped.".format(name)
        module.exit_json(**result)

    elif state == 'restarted':
        if module.check_mode:
            result['changed'] = True
            result['msg'] = "VM '{0}' would be reset/restarted.".format(name)
            module.exit_json(**result)

        reset_cmd = base_cmd + ['vm', 'reset', name]
        rc, out, err = module.run_command(reset_cmd)
        if rc != 0:
            module.fail_json(msg="Failed to reset VM '{0}': {1}".format(name, err or out), rc=rc)

        result['changed'] = True
        result['vm'] = get_vm_info(module, base_cmd, name)
        result['msg'] = "VM '{0}' restarted.".format(name)
        module.exit_json(**result)

    # If state == 'present', VM is confirmed present
    result['vm'] = get_vm_info(module, base_cmd, name)
    result['msg'] = "VM '{0}' is present.".format(name)
    module.exit_json(**result)


def main():
    run_module()


if __name__ == '__main__':
    main()
