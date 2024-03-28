# Hardware Documentation

## Hardware List

- Raspberry Pi Zero 2 W
- Adafruit BME280
    - I2C used for this project
- Adafruit 128x32 PiOLED
    - I2C used for this project
    - Uses the SSD1306 chipset

## Wiring and Wiring Diagrams

The main wiring diagram (it's very simple) is shown in [WiringDiagram.md](./WiringDiagram.md) using a Mermaid diagram.

I'll go into more details on the pins and what connects to what below for each device.

### Raspberry Pi Zero 2 W

The Raspberry Pi's pinout diagrams are shown in [rpio.diagram.md](./rpio.diagram.md), which I pulled from somewhere on the internet. As long as you're not using a super old Pi, you use the Rev 2/3 diagram. Orient yourself to the diagram by holding your Pi so the GPIO rows are running along the right side. A good reference website is [pinout.xyz](https://pinout.xyz) if you want a nicer visual.

We'll be pulling 4 wires from the Pi's GPIO, all from the very top of the GPIO rows:

- Pin 1: 3v3 Power
- Pin 3: I2C1 Data (SDA)
- Pin 5: I2C1 Clock (SCL)
- Pin 6: Ground

Pin 1 and Pin 6 are our 3v3 Power/Ground pair that feeds both the BME280 and PiOLED.

Pin 3 and Pin 5 are our I2C pair that split off to both the BME280 and PiOLED.

The PiOLED hat/diagrams from Adafruit all show the device covering the Pin 2 and Pin 4 5v pins, however they are not needed for it to function.

### Adafruit 128x32 PiOLED

Looking at the back of the display with the text right side up, the pins should be on the bottom left of the device.

- Top Left     <==> (Pi Pin 1) 3v3 Power
- Top Center   <==> (Pi Pin 3) I2C1 Data (SDA)
- Top Right    <==> (Pi Pin 5) I2C1 Clock (SCL)
- Bottom Right <==> (Pi Pin 6) Ground

The bottom left and bottom center pins can be left empty.

### Adafruit BME280

Note: I have a header-only version however the newer STEMMA QT version header pinout is the same.

The BME280 can be run in either SPI or IC2 modes. Since I'm using IC2 in this project, we'll only be populating 4 pins. Looking carefully at the board, you will see each pin marked with it's purpose.

- VIN (Voltage in) <==> (Pi Pin 1) 3v3 Power
- GND (Ground)     <==> (Pi Pin 6) Ground
- SCK (IC2 Clock)  <==> (Pi Pin 5) IC2 Clock (SCL)
- SDI (IC2 Data)   <==> (Pi Pin 3) IC2 Data (SDA)

All other pins can be left empty.