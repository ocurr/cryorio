FROM ubuntu:latest

RUN apt update && apt install  openssh-server sudo -y

RUN useradd -rm -d /home/ubuntu -s /bin/bash -g root -G sudo -u 1000 test 

RUN  echo 'test:test' | chpasswd

RUN echo 'Port 8080' >> /etc/ssh/sshd_config

RUN service ssh start

EXPOSE 8080

CMD ["/usr/sbin/sshd","-D"]
