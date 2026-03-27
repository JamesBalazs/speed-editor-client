# Speed Editor Client

The aim of this project is to provide a fully featured, easy to use library for building open source software for the Davinci Resolve Speed Editor.

There are a few existing solutions in Python, but these are not really designed for consumption as a library and any projects eg keymap editors would
require the user to install Python, dependencies etc. I wanted a more polished solution in a compiled language, that anyone can use to exapnd the functionality
of this hardware.

## Examples

In the examples folder, there are a few projects to get you started:
- [volume wheel](https://github.com/JamesBalazs/speed-editor-client/blob/main/examples/volume_wheel/main.go) uses the Jog wheel as a volume controller for Windows, Mac and Linux
- [lightshow](https://github.com/JamesBalazs/speed-editor-client/blob/main/examples/lightshow/main.go) flashes all the LEDs in each column, then all the LEDs in each row, alternating
- [keypress](https://github.com/JamesBalazs/speed-editor-client/blob/main/examples/keypress/main.go) sets a custom keypress handler, which illuminates the LEDs of the last pressed key
- [reset](https://github.com/JamesBalazs/speed-editor-client/blob/main/examples/reset/main.go) switches off the LEDs

## Usage

To import as a dependency:

```
go get github.com/JamesBalazs/speed-editor-client
```

The project depends on the [go-hid](https://github.com/sstallion/go-hid) library.

Before creating a Speed Editor client, we need to initialize the HID library:

```
if err := hid.Init(); err != nil {
	log.Fatal(err)
}
defer hid.Exit()
```

Don't forget to defer a call to `Exit` to avoid memory leaks.

Next we can initialize the client:

```
client := speedEditor.NewClient()
```

This connects to the Speed Editor, requests the manufacturer info, and device info such as serial number, and sets up the default event handlers.

Device info is cached on initialize, since it will never change once the device is connected:=

```
deviceInfo := client.GetDeviceInfo()

fmt.Printf("Manufacturer: %s\nProduct: %s\nSerial: %s\n", deviceInfo.MfrStr, deviceInfo.ProductStr, deviceInfo.SerialNbr)
```

Which will output something like:
> Product: DaVinci Resolve Speed Editor
> Serial: 1234567890ABCDEFGHIJKLMNOPQRSTUV

The Speed Editor won't work without the correct auth handshake. Luckily for us [Sylvain Munaut reverse engineered and implemented the handshake here](https://github.com/smunaut/blackmagic-misc/blob/master/bmd.py#L133) all the way back in 2021, and published the code under an Apache 2.0 License for the benefit of others.

I re-implemented his authentication algorithm in Go, and exported the underlying functions for consumers of the library to use as they see fit.

When using the client, you just need to call `Authenticate` before sending / receiving any messages, and the handshake will be handled for you:

```
client.Authenticate()
```

Finally, to receive messages from the Speed Editor, you can call `Poll`. This will start a loop which does a blocking read, waiting for either a keypress, battery report, or jog wheel movement from the device:

```
client.Poll()
```

When any of the aforementioned events happen, the corresponding Handler function is called. 

The event handlers can be overridden by the user to implement custom functionality:

```
func customJogHandler(client speedEditor.SpeedEditorInt, report input.JogReport) {
  fmt.Printf("Jog wheel position: %d\n", report.Value)
}

client.SetJogHandler(customJogHandler)

func customKeyPressHandler(client speedEditor.SpeedEditorInt, report input.KeyPressReport) {
  for _, key := range report.Keys {
    fmt.Printf("Keys pressed: %s", key.Name)
  }
}

client.SetKeyPressHandler(customKeyPressHandler)

func customBatteryHandler(client speedEditor.SpeedEditorInt, report input.BatteryReport) {
  fmt.Printf("Battery level: %d", report.Battery)
}
```

The library also provides a complete list of keys for the device, their IDs, the IDs of their LEDs (if present for a key), their labels and their positions on the board.

This helps light LEDs based on their position such as in the lightshow and volume wheel examples.

You can light any combination of LEDs on the board:

```
keysByName := keys.ByName()

leds := []uint32{keysByName[keys.CAM7.Led], keysByName[keys.CAM5.Led], keysByName[keys.CAM3.Led]}
jogLeds := []uint8{keysByName[keys.SHTL.JogLed], keysByName[keys.JOG.JogLed], keysByName[keys.SCROLL.JogLed]}

client.SetLeds(leds)
client.SetJogLeds(jogLeds)
```

`JOG`/`SCRL`/`SHTL` are on a different system to the other LEDs, so a different function is required to light them.

Finally, there are a few different jog modes available on the device:

- `RELATIVE`
  - Reports relative position, since last report
- `ABSOLUTE`
  - Reports absolute position, where 0 is the position when the mode was set. -4096 -> 4096 = 180deg
- `RELATIVE_2`
  - Same as `RELATIVE`, but I think Davinci Resolve uses this to enable faster scrolling when jog is double pressed in later versions (according to one obscure forum post). You could replicate this in software by applying some multiplier to the relative position received from the device (or use it for any feature you like)
- `ABSOLUTE_DEADZONE`
  - Same as `ABSOLUTE` but with a deadzone around 0, so less sensitive to accidental knocks / easier to reset to 0

You can switch modes via the client:

```
client.SetJogMode(jogModes.ABSOLUTE)
```

You will have to handle lighting the buttons yourself, if you want the modes to work like they do with the editor connected to Davinci.

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

Attaching to WSL (does not persist reboot):
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

### Cross platform builds for Windows

mingw-w64 is required to compile the HID library on Linux for Windows with CGO. Installation:
```
sudo dnf install mingw64-gcc
```

To build the examples:
```
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build main.go
```

# Disclaimer

I am not affiliated with Blackmagic Design.

Blackmagic Design, Davinci Resolve, and the Speed Editor are all registered trademarks of Blackmagic Design.

I did not reverse engineer or write the original handshake algorithm.

The [EU Software Directive (2009/24/EC)](https://www.wipo.int/wipolex/en/legislation/details/8612) explicitly permits reverse engineering for interoperability. In this case the handshake algorithm is being used purely for interoperability between the Speed Editor and other software.

# Thanks

Thanks to [smunaut](https://github.com/smunaut) for reverse engineering the Speed Editor authentication algorithm and publishing it.

This is by far the hard part of getting the device working with software other than Davinci Resolve, which I had dreamed of doing but would have never realistically had time (nor probably the skills) to do.

Without their work this library would not have been possible.
