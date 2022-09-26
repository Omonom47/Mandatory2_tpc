# Mandatory2_tpc
a) What are packages in your implementation? What data structure d
you use to transmit data and meta-data?
- 


b) Does your implementation use threads or processes? Why is it not
realistic to use threads?
- Our implemention use threads, which is not realistic because in the real world 


c) How do you handle message re-ordering?
- We use the a sequence number for each of the packets that are being appened to a slice together with their corresponding packet. After all the packets have been sent and recieved they are then reordered in their correct order according to their sequence number.


d) How do you handle message loss?
- In the current code we dont take message loss into consideration, but if we were to implement it, we would make a check in our Server() function that makes sure the packets recieved are the same as the length of the packets. If it isn't then the Server() will ask the Client() to send the missing packets again, and then reorder the packets and outprint the message recieved


e) Why is the 3-way handshake important?
- 

