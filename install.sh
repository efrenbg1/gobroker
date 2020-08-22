user=`id -u`
if [ $user != "0" ]; then
    echo "\e[31mThis script must be run with privileges!\e[0m"
    exit
fi
systemctl stop gobroker
rm gobroker 2> /dev/null
echo "\e[34mCompiling gobroker...\e[0m"
go build
if [ -f "gobroker" ]; then
    echo "\e[34mInstalling dependencies...\e[0m"
    cp gobroker /usr/local/bin/gobroker
    cp /etc/gobroker/settings.json /etc/gobroker/settings.json.old 2>/dev/null || :
    mkdir -p /etc/gobroker
    cp settings-example.json /etc/gobroker/settings.json
    cp -r cert /etc/gobroker/
    echo "\e[34mInstalling systemd service...\e[0m"
    cp gobroker.service /lib/systemd/system/
    systemctl daemon-reload
    systemctl enable gobroker
    echo " \e[31m→\e[0m Remember to first edit your settings in '\e[35m/etc/gobroker/settings.json\e[0m'"
    echo " \e[31m→\e[0m To start the gobroker use: '\e[35msystemctl start gobroker\e[0m'"
    echo "\e[92mAll done. Have a great day!\e[0m"
    exit
fi

echo "\e[31mCan't install gobroker (compilation was unsuccessful)\e[0m"
exit