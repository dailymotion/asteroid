# Asteroid

**Asteroid** is a tool designed to help you manage your Wireguard server.

It supports :
- Viewing the configuration to check peers/IP relation 
- Adding a new peer
- Removing a peer

## Use cases

### Viewing the configuration
If you want to check how much peers you have on your server and be able to pinpoint which ip they relate to, use the view flag

```
$ asteroid view
+----------------------------------------------+------------+
|                     PEER                     |  LOCAL IP  |
+----------------------------------------------+------------+
| eXaMPL3Ave8q+Da!L!m0ti0niXc//M0EGOY6K9C11nw= | 172.16.0.2 |
| eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C12nw= | 172.16.0.4 |
| eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C13nw= | 172.16.0.5 |
| eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C14nw= | 172.16.0.6 |
+----------------------------------------------+------------+
``` 

### Adding a peer
As one of the perk of asteroid is to let you add new peer, you just have to prepare the peer public key and it's Wireguard IP

```
$ asteroid add -key "eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw=" -address "172.16.0.x/xx"
  
  Peer: eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= has been added !
``` 

### Deleting a peer
If you need for any reasons to remove a peer from the server, a simple way to do it will be with the `delete` command:

```
$ asteroid delete -key "eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw="
 
  Peer: eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= has been deleted !
``` 

### Help
If you need any help with the command and their respective arguments, just use help:
```
$ asteroid help / -h
``` 