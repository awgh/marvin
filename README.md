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


## Marvin Help

Marvin responds to private messages privately and responds to channel commands as notices,
with the exception of the .5questions command, where the response is always broadcast to the channel.
The following commands are available:

**.5questions [username]**
(alias: .5)
will broadcast the Five Questions, with an optional greeting for username to the channel.

**.x4questions [username]**
(alias .x4)
will ask additional four Questions.

**.booze [booze_name_or_prefix]**
(alias: .b)
will list Boozes used in the mixed drink database.  This works as a string prefix search.
If there is more than one match, all matches will be listed.  If no argument is given, all Boozes will be listed.
If only one Booze matches, the list of Drinks using that Booze will be shown.

**.drink [drink_name_or_prefix]**
(alias: .d)
will display Drink recipes from the mixed drink database.  This works as a string prefix search.
If there is more than one match, all matches will be listed.  If no argument is given, all Drinks will be listed.
If only one Drink matches, the recipe for that drink will be shown.

**.tell <nick> <message>**
(alias: .t)
will send a message to nick the next time they join or talk in channel.  Private tells will be sent privately.
