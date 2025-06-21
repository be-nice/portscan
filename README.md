# PortScanner by Sindustries

## Warning

Only use it on networks you have permission to scan!  
It port scanning can be a crime

## Introduction

PortScanner is a tool to scan networks for potential  
vulnerabilities in their port configuration.

## Tech specs

This tool can scan  
-- TCP ports  
-- UDP ports  
-- Stealth (not returning 3 way handshake) scan on TCP

Identifies services that are know to run on common ports.  
Designed to run multithreaded for increased performance.

## Usage

```bash
go run main.go <ip> <optional flags and parameters>
```

### Optional flags and parameters

**PORT Flag:  
Sets the ports to be scanned  
Accepts single port, port list, or port range.  
If no flag is provided, defaults to <-p 0-65535>**

```bash
<-p> [port value]
Examples-
<-p 443> single port
<-p 80,443,8080> port list
<-p 80-4000> port range
```

**PROTOCOL Flag:  
Sets the protocol used  
If no flag is provided, defaults to <-tcp->**

```bash
<-tcp> Scanning TCP protocol
<-udp> Scanning UDP protocol
<-s-> Stealt scanning TCP protocol
```

**OUTPUT Flag:  
Sets the way results are displayed  
If no flag is provided, defaults to <-a>**

```bash
<-a> Prints all results to terminal
<-o> Prints only open responses to terminal
<-f> [filename] Writes all results to text file
```

**PERFORMANCE Flag:  
Sets number of threads used for scanning  
If no flag is provided, defaults to <-w 10>  
If thread count is less than 1, defaults to <-w 10>  
If thread count is more than 100, defaults to <-w 100>**

```bash
<-w> [thread count]
```
