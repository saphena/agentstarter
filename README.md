# agentstarter

This executable is used on bootup or some other time to start a browser interface
to the Agent DVR service running on the host PC.

Agent is responsible for capturing camera feeds and is configured using a popup menu on 
the browser interface. Camera IP addresses, etc are all configured using that interface,
this application has nothing to say about those things.

The cameras themselves are TP-Link Tapo C100 units and each must be attached to the
local WiFi network and must also have a "Camera account" configured using the Tapo phone
app. The credentials for the account should be 'saphena' and '201053'.
