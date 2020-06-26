# gobroker
This repository contains a broker of Message Queuing Telemetry over SSL/TLS made in go
# MqTLS
While the MQTT protocol performs very well under networks with limited bandwidth, the support of MQTT under TLS on microcontrollers is very limited. The very popular ESP8266 does not have hardware support for TLS which means that it must be emulated via software. Having then an even simpler subscribe-publish protocol allows to implement telemetry securely on this board.
In MqTLS publish packets have payload length from 12 to 208 bytes. This can be shortened even more, but to make implementations as easy as possible the protocolo has some overhead by encoding integers in utf-8.


# Is it the same as MQTT?
The protocol is not the same, and is not intended to be compatible. Its nature derives from a combination of both MQTT and the readable nature of HTTP. This means that newer servers and clients will have to be developed.<br><br>
Some properties:<br>
1. All data is encoded using utf-8 (including lengths of fields and the actions)<br>
2. Data has a maximum of 99 characters<br>
3. Both ip and port are stored in the server for client check (asigned when login)<br>
4. All messages are retained<br>
5. All messages should end with '\n'
6. Instead of subscribe with callbacks, the action is called watch, which any change made to a topic and its slots will send a messge back to all watchers (max. 5 per topic)<br>
7. Lengths are defined by the digits, not the bytes<br>
8. Because of performance reasons, all communications are QoS 0. Only the connect, retrieve and last-will packets should expect response.

# The structure
**<p align="center">Connect packet: → MQS004user08password15\n ←</p>**
MQS → Protocol name<br>
0 → Action type, authetificate in this case<br>
04 → length of user field<br>
user → utf-8 encoded string<br>
08 → length of password field<br>
password → utf-8 endoded string<br>
15 → time to keep the connection live in seconds (if not specified default is 10s, maximum is 99s)<br><br>

**<p align="center">Connect packet response: → MQS0\n ←</p>**
MQS → Protocol name<br>
0 → Action type (expected the same as the sent packet) if equals 9 then then the action failed<br>
<br><br>

**<p align="center">Last-will packet: → MQS317/home/temperature303Off\n ←</p>**
MQS → Protocol name<br>
3 → Action type, last-will in this case<br>
17 → length of topic field<br>
/home/temperature → utf-8 encoded string for topic<br>
3 → utf-8 encoded single digit integer to represent which slot to use in the topic (each topic has 10 slots)<br>
03 → length of data field<br>
Off → utf-8 endoded string for data<br><br>

**<p align="center">Publish packet: → MQS117/home/temperature10416ºC\n ←</p>**
MQS → Protocol name<br>
1 → Action type, publish in this case<br>
17 → length of topic field<br>
/home/temperature → utf-8 encoded string for topic<br>
1 → utf-8 encoded single digit integer to represent which slot to use in the topic (each topic has 10 slots)<br>
04 → length of data field<br>
16ºC → utf-8 endoded string for data<br><br>

**<p align="center">Retrieve packet: → MQS217/home/temperature2\n ←</p>**
MQS → Protocol name<br>
2 → Action type, retrieve in this case<br>
17 → length of topic field<br>
/home/temperature → utf-8 encoded string for topic<br>
2 → utf-8 encoded single digit integer to represent which slot to use in the topic (each topic has 10 slots)<br><br>

**<p align="center">Retrieve packet response: → MQS20416ºC\n ←</p>**
MQS → Protocol name<br>
2 → Action type, retrieve in this case<br>
04 → length of payload field<br>
16ºC → utf-8 encoded string for the payload<br><br>

**<p align="center">Watch packet: → MQS417/home/temperature\n ←</p>**
MQS → Protocol name<br>
4 → Action type: watch<br>
17 → length of topic field<br>
/home/temperature → utf-8 encoded string for the topic<br><br>

**<p align="center">Watch on update response packet: → MQS517/home/temperature00416ºC\n ←</p>**
MQS → Protocol name<br>
5 → Action type: watch update<br>
17 → length of topic field<br>
/home/temperature → utf-8 encoded string for the topic<br>
0 → utf-8 encoded single digit integer to represent which slot has been updated<br>
04 → length of payload field<br>
16ºC → utf-8 encoded string for the payload<br><br>
<br><br><br><br>

There are many more things, but as I'm the only currently using it I will leave the README as is.

# Use cases for development of lemonSW:

///// topic = BC:DD:C2:08:8C:BE ///// (only MAC) (status topic)<br>
0 -> PC off<br>
1 -> PC on<br>
2 -> PC suspended<br>
3-6 -> not used<br>
7 -> Board is updating<br>
8 -> Board is in recovery mode<br>
9 -> Board is off<br><br>



///// topic = _BC:DD:C2:08:8C:BE ///// ("_" + MAC)  (actions topic)<br>
0 -> Simple power button click (to turn off or on)<br>
1 -> Force off (holds button for 5s)<br>
2 -> Reset PC (actions 1 and 0 combined)<br>
3-5 -> not used<br>
4 -> Turning off<br>
5 -> Action failed<br>
6 -> Performed action (default status, waiting for command or cancel action)<br>
7 -> Update<br>
8 -> Enter recovery mode<br>
9 -> Update settings<br><br>



///// topic = !BC:DD:C2:08:8C:BE ///// ("!" + MAC)  (settings topic)<br>
off1640on0915 -> Turn off at 16:40, turn on at 09:15<br><br>
