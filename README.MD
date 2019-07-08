    
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
#usage
    # git clone https://github.com/cadyrov/ftploader
    
    # go build /home/SERVERNAME/ftploader/
    
    #crontab /etc/crontab 
    #crontab -e
    * * * * * /home/SERVERNAME/ftploader/main
    
    curl or http SERVERNAME:3000/download?path=/home/SERVERNAME/ftploader/example.yml&log=/home/SERVERNAME/www/logs/log.log

    set config ftp 
    #mcedit /home/SERVERNAME/ftploader/config.yml

    set files to download 
    path=/home/SERVERNAME/ftploader/example.yml
    
