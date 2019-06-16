
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
