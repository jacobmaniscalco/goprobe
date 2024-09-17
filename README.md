# GoProbe IoT Vulnerability Scanner

<p align="center">
    <img src="./img/goprobe.png" alt="logo" width="300">
</p>


## Overview
GoProbe is a powerful and modular command-line tool designed to assess the security of a wide range of IoT (Internet of Things) devices. Built in Go, this tool provides security professionals and developers with an effective means to identify and report vulnerabilities across diverse IoT ecosystems, regardless of device type or application.

## Project Goals

* Versatile Scanning: Develop a CLI tool in Go for scanning and identifying vulnerabilities across various IoT devices.
* Modular Design: Ensure the tool is modular, facilitating easy integration of new scanning modules and plugins.

## Core Functionalities

* Device Discovery:
    * Scan networks to identify connected IoT devices using protocols like ARP, mDNS, SSDP, and UPnP.
    * Perform detailed service enumeration to identify running services and their versions.
    * Classify devices based on type, manufacturer, and known vulnerabilities.
    * Powered by the Ullaakut/nmap Go package.

* Vulnerability Scanning:
    * Utilize known vulnerabilities databases (e.g., CVEs) to detect existing security issues.
    * Check configurations for weak or default settings and insecure firmware versions.
    * Implement brute-force modules for common IoT services (e.g., SSH, Telnet) to test for weak credentials.

* Protocol Analysis:
    * Scan for outdated or insecure firmware versions.
    * Detect the presence of hard-coded credentials and other common firmware vulnerabilities.

* Reporting:
    * Generate detailed reports in formats like JSON, XML, and PDF, outlining identified vulnerabilities and recommended actions.
    * Offer both summary and detailed modes for varied analysis needs.

## Advanced Features

* Modular Plugin System:
    * Enable users to develop and integrate custom scanning modules.
    * Provide an API for adding new vulnerability checks or device types.

* Automation & Scripting:
    * Support integration with automation frameworks (e.g., Ansible, Jenkins) for routine scans.
    * Allow creation of custom scripts for specialized scanning scenarios.


## Getting Started
Instructions on how to install the tool will be added soon. In addition, a Dockerfile will be added with the necessary software to test 
the tool with [IoTGoat](https://github.com/OWASP/IoTGoat) and [Damn Vulnerable IoT Device (DVID)](https://github.com/Vulcainreo/DVID)

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.