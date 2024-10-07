## ip2org

## Installation
```
go install github.com/rix4uni/ip2org@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/ip2org/releases/download/v0.0.1/ip2org-linux-amd64-0.0.1.tgz
tar -xvzf ip2org-linux-amd64-0.0.1.tgz
rm -rf ip2org-linux-amd64-0.0.1.tgz
mv ip2org ~/go/bin/ip2org
```
Or download [binary release](https://github.com/rix4uni/ip2org/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/ip2org.git
cd ip2org; go install
```

## Usage
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

## Usage Example
```
echo "209.237.150.21" | ip2org
cat ips.txt | ip2org

shodan parse --fields ip_str --separator ":" *.json.gz | unew | ip2org
```

## Output
```
209.237.150.21 [Web.com Group, Inc.]
```
