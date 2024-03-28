```mermaid
flowchart LR
    subgraph RaspPi
        3v;Gnd;SCL;SDA
    end
    subgraph PiOLED
        Oled3v;OledSCL;OledSDA;OledGnd
    end
    subgraph BME280
        BMEVin;BMEGnd;BMESCK;BMESDI
    end

    3v --Red--> Oled3v & BMEVin
    Gnd --Black--> OledGnd & BMEGnd
    SCL --Blue--> OledSCL & BMESCK
    SDA --Green--> OledSDA & BMESDI
```
