# ip2org

# Installation
```
go install github.com/rix4uni/ip2org@latest
```

##### via clone command
```
wget https://raw.githubusercontent.com/rix4uni/ip2org/main/ip2org.go && go build ip2org.go && mv ip2org ~/go/bin/ip2org && rm -rf ip2org.go
```

# Usage
```
  -ip string
        IP address to lookup
  -list string
        File containing IP addresses
  -o string
        File to save output
  -reverse
        Perform reverse DNS lookup
  -timeout int
        Timeout for whois lookup in seconds (default 10)
  -verbose
        Enable verbose mode
  -version
        print version information and exit
```

# Usage Example
```
echo "209.237.150.21" | ip2org
cat ips.txt | ip2org

shodan parse --fields ip_str --separator ":" *.json.gz | unew | ip2org
```

# Output
```
209.237.150.21 [Web.com Group, Inc.]
```
