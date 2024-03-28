```mermaid
flowchart LR
    RaspPi;PiOLED;BME280
    3v;Gnd;SCL;SDA
    Oled3v;OledSCL;OledSDA;OledGnd
    BMEVin;BMEGnd;BMESCK;BMESDI

    RaspPi --> 3v & Gnd & SCL & SDA

    3v --Red--> Oled3v & BMEVin
    Gnd --Black--> OledGnd & BMEGnd
    SCL --Blue--> OledSCL & BMESCK
    SDA --Green--> OledSDA & BMESDI

    Oled3v & OledSCL & OledSDA & OledGnd --> PiOLED
    BMEVin & BMEGnd & BMESCK & BMESDI --> BME280
```
