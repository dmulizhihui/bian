FROM my-ubuntu:v1
WORKDIR /root/getsomething
COPY . .
CMD ["bash","./z.sh"]