# marvin
IRC and Slack bot from Milliways

## Instructions for Ubuntu 16.04

### Install go
```
curl https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz > go1.9.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.9.linux-amd64.tar.gz
```
### Add go to your environment/PATH
```
sudo echo -ne "export GOPATH=$HOME/go\nexport PATH=$PATH:/usr/local/go/bin" >> /etc/profile
source /etc/profile
```
### Verify go is in your path
`echo $PATH $GOPATH`

### Create go build directories
`mkdir -p $HOME/go/src/`

### Download marvin
```
go get github.com/awgh/madns
cd $HOME/go/src
```

### Install marvin

Replace /opt/marvin with the target directory of your choice.
```
cd $HOME/go/src
./install.sh /opt/marvin
```

### Create Your JSON Config Files
In the /opt/marvin/config directory (or your custom target equivalent), create a JSON file for each chat backend you wish to connect to.

At this time, only ONE CONNECTION TO SLACK is supported, but multiple different IRC networks can be connected to at once, and you can connect to Slack at the same time as multiple IRC networks... but some features aren't completely working Slack<->IRC yet.


#### Example IRC config 
Create a file in the config directory called "ANYTHING.json" and customize the following template with your own IRC information, proxy information (optional), and your free md5decrypt API credentials from md5decrypt https://md5decrypt.net/en/Api/ (optional).

```
{
        "Host":"irc.freenode.net",
        "Port":"6697",
        "Nick":"YOUR_NICK_HERE",
        "Password":"YOUR_PASSWORD_HERE",
        "Channel":"#ratnet",
        "Name":"YOUR_NAME_HERE",
        "Version":"IRCsome v1.0",
        "Quit":"bye!",
        
        "ProxyEnabled":false,
        "Proxy":"socks5://PROXY_USER:PROXY_PASSWORD@PROXY_HOST:1080",
        
        "MD5ApiUser":"YOUR_MD5_API_USER",
        "MD5ApiCode":"YOUR_MD5_API_KEY"
}
```

#### Example Slack config 
Create a file in the config directory called "slack.json" (really anything with a .json extension) and customize the following template with your own Slack information and your free md5decrypt API credentials from md5decrypt https://md5decrypt.net/en/Api/ (optional).

```
{
        "SlackAPIToken":"xorx-YOUR-SLACK-APP-API-TOKEN-HERE",
        "SlackChannel":"SLACK_CHANNEL_ID",
        
        "MD5ApiUser":"YOUR_MD5_API_USER",
        "MD5ApiCode":"YOUR_MD5_API_KEY"
}
```

You will have to create a custom Slack App and add it to your Slack in order to get an app API token.

The channel ID has to be the funky Slack channel ID string, not the human-readable channel name.

Also, proxy support has not been implemented for Slack yet... not sure if anyone actually wants that anyway.


### Starting Marvin
```
cd /opt/marvin
./marvin
```


# Marvin Help

Marvin responds to private messages privately and responds to channel commands as notices,
with the exception of the .5questions command, where the response is always broadcast to the channel.

The following commands are available:

## Milliways Commands

### .5questions [username]

(alias: .5)

will broadcast the Five Questions, with an optional greeting for username to the channel.

### .x4questions [username]

(alias .x4)

will ask additional four Questions.

## Mixed Drinks Database Commands

### .booze [booze_name_or_prefix]

(alias: .b)

will list Boozes used in the mixed drink database.  This works as a string prefix search.
If there is more than one match, all matches will be listed.  If no argument is given, all Boozes will be listed.
If only one Booze matches, the list of Drinks using that Booze will be shown.

### .drink [drink_name_or_prefix]

(alias: .d)

will display Drink recipes from the mixed drink database.  This works as a string prefix search.
If there is more than one match, all matches will be listed.  If no argument is given, all Drinks will be listed.
If only one Drink matches, the recipe for that drink will be shown.

## Answering Machine Commands

### .tell <nick> <message>

(alias: .t)

will send a message to nick the next time they join or talk in channel.  Private tells will be sent privately.

## Hash Cracking / md5decrypt Commands

If md5decrypt API credentials are provided, the following hash types can be cracked via Marvin:

### .md5 [hash]
### .md4 [hash]
### .sha1 [hash]
### .sha256 [hash]
### .sha384 [hash]
### .sha512 [hash]
### .ntlm [hash]



## Markov Chain Quote Commands

### .m

Hear from Marvin, the paranoid android.

### .mcfly

For a change of tone, hear from Marty McFly.

The quotes can be customized used the markov-gen command from this package:  https://github.com/awgh/markov
