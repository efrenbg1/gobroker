# gobroker
MqTLS broker made with golang


Use cases for development of lemonSW:

///// topic = BC:DD:C2:08:8C:BE ///// (only MAC) (status topic)
0 -> PC off
1 -> PC on
2 -> PC suspended
3-6 -> not used
7 -> Board is updating
8 -> Board is in recovery mode
9 -> Board is off



///// topic = _BC:DD:C2:08:8C:BE ///// ("_" + MAC)  (actions topic)
0 -> Simple power button click (to turn off or on)
1 -> Force off (holds button for 5s)
2 -> Reset PC (actions 1 and 0 combined)
3-5 -> not used
4 -> Turning off
5 -> Action failed
6 -> Performed action (default status, waiting for command or cancel action)
7 -> Update
8 -> Enter recovery mode
9 -> Update settings



///// topic = !BC:DD:C2:08:8C:BE ///// ("!" + MAC)  (settings topic)
off1640on0915 -> Turn off at 16:40, turn on at 09:15
