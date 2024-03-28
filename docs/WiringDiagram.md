```mermaid
flowchart LR
    5v_1;5v_2;3v;Gnd;SCL;SDA
    Oled5v_1;Oled5v_2;Oled3v;OledSCL;OledSDA;OledGnd
    BMEVin;BMEGnd;BMESCK;BMESDI

    5v_1 --Red--> Oled5v_1
    5v_2 --Red--> Oled5v_2
    3v --Orange--> Oled3v & BMEVin
    Gnd --Black--> OledGnd & BMEGnd
    SCL --Blue--> OledSCL & BMESCK
    SDA --Green--> OledSDA & BMESDI
```
