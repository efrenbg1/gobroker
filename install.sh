user=`id -u`
if [ $user != "0" ]; then
    echo "\e[31mThis script must be run as root!\e[0m"
    exit
fi

rm gobroker 2> /dev/null
echo "\e[32mCompiling gobroker...\e[0m"
go build
if [ -f "gobroker" ]; then
    echo "\e[32mInstalling dependencies...\e[0m"
    cp gobroker /bin/gobroker
    mkdir 
    exit
fi

echo "\e[31mCan't install gobroker (compilation was unsuccessful)\e[0m"
exit