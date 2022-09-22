# testITGO
This project was written as a test assignment for an internship

Before starting make sure that Apache couchdb has been installed on your OS
to install it wisit https://github.com/apache/couchdb

After installation open constants.go in main repository to fill login password and etc for using couchdb

Visit Makefile to set ip address and port for server

To start server type "make server" in command line
To start client type "make client" in command line

To refresh pb.go files you may use "make clean" and "make gen" commands

Client understand that commands:
"ewallet create" to create new ewallet with random balance from 1000 to 3000
"ewallet send senderUUID recipientUUID float.amount" to make transaction from one ewallet to another
"ewallet getlast" to get info about ALL getting transaction that newer seen before. After that command all "get" fields remowed

CouchDB also save all "send" transaction just in case
