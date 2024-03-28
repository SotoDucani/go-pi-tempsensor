---
id: hJYjn92qjuneiTmb7CIBC
title: Wiring Diagram
desc: ''
updated: 1642869923951
created: 1642868061856
---

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
