Smart home node
---------------

#####Instalation

access to serial port

sudo gpasswd --add ${USER} dialout
    
or
    
sudo usermod -a -G dialout ${USER}
    
You then need to log out and log back in again for it to be effective. 

#####Error codes
    
    1 serial port errors 
    2 modbus line errors
    3 tcp read bytes errors
    4 unmarshal bytes to json from tcp errors

#####TODO

* работа в качестве демона https://github.com/takama/daemon
* доступ по сертификату
* shell console ?

#####Протокол основанный но modbus

* ASCII
* проверка целостности пакета по контрольной сумме LRC
* ограничение по времени ожидания ответа 2сек

