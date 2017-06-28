# marvin

**Marvin Help**

Marvin responds to private messages privately and responds to channel commands as notices,
with the exception of the .5questions command, where the response is always broadcast to the channel.
The following commands are available:

**.5 [username]**
(alias: .5questions)
will broadcast the Five Questions, with an optional greeting for username to the channel.

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
