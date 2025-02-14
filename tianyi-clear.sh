#!/bin/bash
user_password=$(cat /root/tianyi.txt)
/root/tianyi_cli -c clear -d "$user_password" -f 99359679057739 >> /root/tianyi.log
