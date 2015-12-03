# defector

A utility to handle captive portal authentication within Subgraph OS.

It is implemented as a service to detect problems connecting to the Tor
network. When connection problems are detected, it prompts to user to
detect if they are within a captive portal. If so, it will open a minimal
browser (Captive Browser) to perform captive portal authentication. After
successfully authenticating, it will prompt the user to save their current
settings so the authentication may persist on subsequent connections.

defector is designed with the following requirements in mind:

1. It must prompt the user at every stage since captive portal authentication
must be performed non-anonymously (it is not possible over Tor because it is
not possible to establish a connection to the Tor network prior to 
authentication)

2. It must limit the exposure of the user during non-anonymous network access
(this is accomplished by limiting browser functionality and sandboxing
the browser)

3. It should allow the user to keep a persistent session after authenticating
(this will be accomplished by associating the current spoofed MAC address with
the captive portal network and re-using that spoofed MAC on subsequent 
connections to the network -- if the user chooses to do so)


