# To run
To start the program simply type "go run tcphanding.go" During the programs runtime you should see a series of line in the in the format where the different data are sendt and recieved as well as when connection to the 5 clients are established. When all the Clients have been connected to the server and outprinted the 5 messages recieved the program will terminate

# Mandatory2_tpc
**a) What are packages in your implementation? What data structure did
you use to transmit data and meta-data?**
- Our implementation of TCP sends the message through our structs of packets, which we have made to contain most of the different parameters a tcp has. 
Our meta-data is also handled in our packet structure and created in our *MakePacket()* function, which initializes the fields. Most of these are random integers set in *FragmentMessage()*, they are aminly there for the understanding of TCP.
In *FragmentMessage()* we split up our message string into a slice of packets. This is our data, that we send from the *Client()* to the *Server()*. Each of the individual packets of our string are assigned a sequence number, which keeps track of the placement of the different packets.


**b) Does your implementation use threads or processes? Why is it not
realistic to use threads?**
- Our implemention uses threads, which is not realistic because in the real world, as TCP runs through the kernel and therefore can't be described through threads. Furthermore a process is any program in execution whereas a thead only runs a segment of a process. Threads also share memory, whereas a process has its own isolated memory, which makes global variables unusable(as we have done right now).


**c) How do you handle message re-ordering?**
- We use the sequence number for each of the packets that are being appended to a slice together with their corresponding packet. After all the packets have been sent and recieved they are then reordered in their correct order according to their sequence number using GO's built-in sorting.


**d) How do you handle message loss?**
- In our current code we do not simulate message loss, but if do so we would make a check in our *Server()* function that makes sure the packets recieved, and send a confirmation to *Client()*. If it doesn't get a confirmation inside a certain timeframe then it resends the packet.


**e) Why is the 3-way handshake important?**
- The 3-way handshake is essential to create a secure connection between the client and server. Furthermore it makes sure both parts are ready to send and recieve data. 
Using a 3-way handshake the client and server agree on the different parameters in form of TCP segments. These segments will check and verify incoming and outgoing packets of data.

