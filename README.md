# Mandatory2_tpc
**a) What are packages in your implementation? What data structure did
you use to transmit data and meta-data?**
- Our implementation of TCP sends the message through our structs of packets, which we have made to contain most of the different parameters a tcp has. 
Our meta-data is also handled in our packet structure and created in our *MakePacket()* function, which initializes the our fields. Most of these are random integers set in *FragmentMessage()*, but are there for the understanding of TCP.
In *FragmentMessage()* we split up our message string into a slice of packets. This is our data, that we sent from the *Client()* to the *Server()*. Each of the individual packets of our string are assigned a sequence number, which keeps track of the placement of the different packets.


**b) Does your implementation use threads or processes? Why is it not
realistic to use threads?**
- Our implemention use threads, which is not realistic because in the real world, as the TCP runs through the kernel and therefore can't be described through threads. Furthermore a process is any program in execution whereas a thead only runs a segment of a process. Threads also share memory, where as a process has their own isolated memory, which makes global variables unusable(as we have done right now).


**c) How do you handle message re-ordering?**
- We use the a sequence number for each of the packets that are being appened to a slice together with their corresponding packet. After all the packets have been sent and recieved they are then reordered in their correct order according to their sequence number.


**d) How do you handle message loss?**
- In the current code we dont take message loss into consideration, but if we were to implement it, we would make a check in our *Server()* function that makes sure the packets recieved are the same as the length of the packets. If it isn't then the *Server()* will ask the *Client()* to send the missing packets again, and then reorder the packets and outprint the message recieved


**e) Why is the 3-way handshake important?**
- The 3-way handshake is essential to create a secure connection between the client and server. Furthermore it makes sure both parts are ready to send and recieved data. 
Using a 3-way handshake the client and server agrees on the different parameters in form of TCP segments. These segments will check and verify incoming and outgoing packets of data.

