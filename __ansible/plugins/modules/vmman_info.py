#!/usr/bin/python3
# -*- coding: utf-8 -*-

# Copyright: (c) 2026, J.F. Gratton <jean-francois@famillegratton.net>
# GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)

from __future__ import absolute_import, division, print_function
__metaclass__ = type

DOCUMENTATION = r'''
---
module: vmman_info
short_description: Gather facts about virtual machines using vmman4
description:
  - Query details of a specific virtual machine or list virtual machines using the vmman4 CLI.
version_added: "1.0.0"
author:
  - J.F. Gratton (@jeanfrancoisgratton)
options:
  name:
    description:
      - Name of the virtual machine to inspect.
    type: str
    required: true
  connection_uri:
    description:
      - Libvirt connection URI (e.g. C(qemu:///system)).
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
- name: Get VM details
  jfgratton.vmman4.vmman_info:
    name: web-node-01
  register: vm_details

- name: Display IP address
  ansible.builtin.debug:
    msg: "VM IP is {{ vm_details.vm.ip_address }}"
  when: vm_details.vm.exists and vm_details.vm.ip_address is defined
'''

RETURN = r'''
vm:
  description: Detailed information about the requested virtual machine.
  type: dict
  returned: always
  contains:
    exists:
      description: Whether the domain exists.
      type: bool
      returned: always
    name:
      description: Domain name.
      type: str
      returned: always
    state:
      description: Domain state (Running, Shutoff, etc.).
      type: str
      returned: when VM exists
    vcpus:
      description: Number of vCPUs.
      type: int
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

    if "not found" in clean_out.lower() or "not found" in clean_err.lower():
        return {'exists': False, 'name': vm_name}

    if rc != 0:
        if "unable to get domain" in clean_err.lower():
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


def main():
    module = AnsibleModule(
        argument_spec=dict(
            name=dict(type='str', required=True),
            connection_uri=dict(type='str', required=False),
            connection_file=dict(type='path', required=False),
            executable=dict(type='str', default='vmman'),
        ),
        supports_check_mode=True
    )

    name = module.params['name']
    connection_uri = module.params['connection_uri']
    connection_file = module.params['connection_file']
    exe_name = module.params['executable']

    executable = locate_executable(module, exe_name)
    base_cmd = build_base_cmd(executable, connection_uri, connection_file)

    info = get_vm_info(module, base_cmd, name)
    module.exit_json(changed=False, vm=info)


if __name__ == '__main__':
    main()
