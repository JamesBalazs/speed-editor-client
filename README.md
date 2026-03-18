## Dev notes

My setup is weird (WSL remote via Zed) so some extra steps are required to pass the Speed Editor through to WSL

Installing [usbipd](https://github.com/dorssel/usbipd-win):
```
winget install usbipd
```

Listing devices:
```
usbipd list
```

Binding the Speed Editor (persists reboot, your BUSID will be different to mine):
```
sudo usbipd bind --busid=4-9
```

Attaching to WSL:
```
sudo usbipd attach --wsl --busid=4-9
```

To confirm w/ [lshid](https://github.com/FFY00/lshid) within WSL:
```
$HOME/go/bin/lshid
```
Should output something like `/dev/hidraw0: ID 1edb:da0e Blackmagic Design DaVinci Resolve Speed Editor`


### Deps

```
sudo dnf install systemd-devel
```
To get `libudev.h` on Fedora (required for lshid)

I then had permission issues reading from `/dev/hidraw0` so had to create a [udev](https://wiki.archlinux.org/title/Udev) rule:
```
KERNEL=="hidraw*", SUBSYSTEM=="hidraw", MODE="0660", GROUP="plugdev"
```
in `/etc/udev/rules.d/99-hidraw-permissions.rules`, then:
```
sudo groupadd plugdev
sudo usermod -a -G plugdev james
sudo udevadm control --reload
sudo udevadm trigger
```

After this `stat /dev/hidraw0` should list the new plugdev group.
