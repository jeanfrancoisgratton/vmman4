- [x] `vmman vm provision $VM $IP fedora2` breaks (see note1)`<br> 
- [x] `vmman vm provision $VM $IP opensuse2` breaks (see note2)<br>
- [ ] inconsistent colored output in `vmmman vm provision`

___

Note 1 : 
```bash
[12:15:36|jfgratton@oslo:~]: vmman vm provision training 10.50.1.1 fedora2

ℹ️  Connected on hypervisor qemu:///system

⏳️ Cloning fedora2.qcow2 to training.qcow2. This might take a while...
✅️ Volume fedora2.qcow2 cloned  as training.qcow2
✅️ training.qcow2 cloned from fedora2.qcow2
✅️ training DEFINED
✅️ training started, waiting for QEMU guest agent
runProvisionScript: guest-side provisioning script failed: exit=1 output=/tmp/.vmman4-provision.sh: line 5: /etc/hostname: Permission denied

```



Note 2 :
```bash
[12:17:22|jfgratton@oslo:~]: vmman vm provision training 10.50.1.1 opensuse2

ℹ️  Connected on hypervisor qemu:///system

⏳️ Cloning opensuse2.qcow2 to training.qcow2. This might take a while...
✅️ Volume opensuse2.qcow2 cloned  as training.qcow2
✅️ training.qcow2 cloned from opensuse2.qcow2
✅️ training DEFINED
✅️ training started, waiting for QEMU guest agent
qgaCall: guest-exec failed: virError(Code=113, Domain=10, Message='guest agent command failed: unable to execute QEMU agent command 'guest-exec': Command guest-exec has been disabled: the command is not allowed')

```

