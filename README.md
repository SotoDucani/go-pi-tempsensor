# go-pi-tempsensor

This is a Go project that utilizes:

- Raspberry Pi Zero 2 W
- Adafruit BME280 Environmental sensor
- Adafruit 128x32 PiOLED display

These can come together to create a small form-factor device that can monitor environmental data (temperature, humidity, and pressure), display it on a screen, and connect to a wireless network for management as well as presenting the data via prometheus metrics.

This is my first project actually building something from scratch with a Pi, and is more a learning experience than a guide or instructions.

Details on the hardware and wiring everything together is in the [docs](./docs) folder.