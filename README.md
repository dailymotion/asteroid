# Asteroid

**Asteroid** is a tool designed to help you manage your Wireguard server.

It supports :
- Viewing the configuration to check peers/IP relation 
- Adding a new peer
- Removing a peer
- Generating Client config file

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

### Generating a stdout of the client Wireguard config 
If you need, for the peer added to generate a Wireguard stdout, you can do it with the flag `-generateStdout`

```
$  asteroid add -key "eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw=" -address "172.16.0.x/xx" -generateStdout
 
   Peer: eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= has been added !
   
   Client config output
   ---------------
   [Interface]
   PrivateKey = eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= 
   Address = xx.xx.xx.xx/32
   DNS = 9.9.9.9
  
   [Peer]
   PubblicKey = eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= 
   AllowedIPs = 0.0.0.0/0
   EndPoint = public_ip_of_your_wireguard
``` 

### Generating client config file
If you need, for the peer added to generate its Wireguard config, you can do it with the flag `-generateFile`  
It's possible to change it's name via the asteroid.yaml config file

```
$  asteroid add -key "eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw=" -address "172.16.0.x/xx" -generateFile
 
  Peer: eXaMPL3Ave8q+kmNVmiw4KdKiXc//M0EGOY6K9C11nw= has been added !
  The wireguard config for the client has been created with the name: wg0.conf
``` 

### Help
If you need any help with the command and their respective arguments, just use help:
```
$ asteroid help / -h
``` 

### Configuration
You just have to copy the file in the config folder into your $HOME with your configuration
```
$ cp config/asteroid_example.yaml $HOME/.asteroid.yaml
```