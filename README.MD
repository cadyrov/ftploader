    
#Install go    
    #yum update
    #wget https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
    #tar -xzf go1.12.6.linux-amd64.tar.gz
    #mv go /usr/local
    #export GOROOT=/usr/local/go
    #mcedit /etc/profile
    add to file end export PATH=$PATH:/usr/local/go/bin
    #reboot
#check 
    #go env

