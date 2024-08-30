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
echo "209.237.150.21" | ip2org
cat ips.txt | ip2org

shodan parse --fields ip_str --separator ":" *.json.gz | unew | ip2org
```

# Output
```
209.237.150.21 [Web.com Group, Inc.]
```
