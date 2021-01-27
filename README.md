# Overview
This is a basic script I wrote to automate the getuserrealm.srf enumeration trick.

Office 365 provides information to mail clients to be used during setup, including the authentication configuration and SSO url if federated auth is used. 

We can take advantage of this by requesting this configuration manually, and extracting useful information from it.

As of now, the tool is completely passive from an opsec perspective, as it only hits login.microsoftonline.com one time. 

In the future, I plan to add basic scraping and DNS discovery functionality so that enumeration beyond the O365 configuration is possible.

# Requirements
=======
# Example
![image](https://user-images.githubusercontent.com/32488787/80498636-f84a9500-8939-11ea-8193-71887ee4f83d.png)


## requirements
go get github.com/gookit/color

# Example Usage
![image](https://user-images.githubusercontent.com/32488787/106001583-218d8280-607e-11eb-86f3-60f9b42d0f53.png)
