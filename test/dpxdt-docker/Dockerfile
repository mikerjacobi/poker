FROM ubuntu:latest

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y git python python-mysqldb phantomjs imagemagick
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y wget 
RUN apt-get autoremove -y
RUN apt-get clean all
RUN wget https://bootstrap.pypa.io/get-pip.py
RUN python get-pip.py

RUN git clone https://github.com/bslatkin/dpxdt.git
WORKDIR /dpxdt
RUN git submodule update --init --recursive
RUN pip install -r requirements.txt
RUN python setup.py install
RUN pip install -e .
RUN python -c "from dpxdt import server; server.db.drop_all(); server.db.create_all()"

#RUN sed -i 's/.setLevel(logging.DEBUG)/.setLevel(logging.INFO)/' dpxdt/runserver.py

RUN echo "SECRET_KEY = \"$(openssl rand -base64 32)\"" > secrets.py

ADD ./launch.sh /dpxdt/launch.sh

EXPOSE 5000
CMD ["./launch.sh"]
