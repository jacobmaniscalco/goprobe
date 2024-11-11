# GoProbe IoT Vulnerability Scanner

<p align="center">
    <img src="./img/goprobe.png" alt="logo" width="300">
</p>


## Overview

GoProbe is a powerful, modular command-line tool designed to assess the security of Internet of Things (IoT) devices, specifically targeting vulnerabilities in two key communication protocols: Bluetooth Low Energy (BLE) and MQTT. Built in Go, this tool provides security professionals and developers with an effective means to identify and report protocol-specific vulnerabilities, especially in authentication mechanisms, across diverse IoT ecosystems.

## Project Goals

* Targeted Protocol Security: Focus on scanning and identifying vulnerabilities in BLE and MQTT, which are widely used in IoT devices.
* In-Depth Vulnerability Analysis: Emphasize in-depth analysis of BLE and MQTT communication protocols, examining areas prone to insecure implementations.
* Modular Design: Maintain a modular framework to facilitate the addition of future protocol analysis or vulnerability assessment modules.

## Core Functionalities

* Device Discovery:
    Scan for BLE-enabled IoT devices and identify active MQTT connections.
    Perform service enumeration on discovered devices to identify protocol usage and version information.
    Classify devices based on type, manufacturer, and potential vulnerabilities.

* Protocol-Specific Vulnerability Scanning:
    * Bluetooth Low Energy (BLE):
        Identify vulnerabilities in BLE 4.0 and 4.1 authentication mechanisms, focusing on weaknesses in key exchange.
        Utilize passive BLE sniffing for packet analysis, using the nRF52840 MDK USB Dongle and Wireshark with the Nordic nRF Sniffer.
    * MQTT:
        Analyze MQTT broker configurations and inspect for weak credentials or insecure authentication methods.
        Test for open MQTT channels and lack of encryption or access control.

* Protocol Analysis:
    Detect the use of insecure configurations or outdated versions of BLE and MQTT protocols.
    Identify insecure key exchange processes and other common vulnerabilities in BLE and MQTT implementations.

* Reporting:
    Generate detailed reports in formats like JSON and PDF, outlining identified protocol vulnerabilities and recommended mitigations.
    Provide summary and detailed reporting options to cater to different levels of analysis requirements.

## Getting Started
Instructions on how to install the tool will be added soon. In addition, a Dockerfile will be added with the necessary software to test 
the tool and run the tool.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.